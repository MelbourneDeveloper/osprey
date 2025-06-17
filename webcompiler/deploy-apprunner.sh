#!/bin/bash

set -e
source "$(dirname "$0")/config.sh"
cd "$(dirname "$0")"

# Verify AWS credentials
aws sts get-caller-identity >/dev/null || { echo "AWS credentials not configured"; exit 1; }

# Get AWS account details now that credentials are verified
get_aws_account_id

# Verify ECR image exists
if ! aws ecr describe-images --repository-name "${ECR_REPOSITORY}" --image-ids imageTag="${IMAGE_TAG}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    echo "ECR image ${FULL_IMAGE_URI} not found!"
    echo "Run ./deploy-ecr.sh first to build and push the image."
    exit 1
fi

# https://www.cloudflare.com/en-au/ips/
# Cloudflare IPv4 ranges for AWS-level security group filtering
CLOUDFLARE_CIDRS="173.245.48.0/20,103.21.244.0/22,103.22.200.0/22,103.31.4.0/22,141.101.64.0/18,108.162.192.0/18,190.93.240.0/20,188.114.96.0/20,197.234.240.0/22,198.41.128.0/17,162.158.0.0/15,104.16.0.0/13,104.24.0.0/14,172.64.0.0/13,131.0.72.0/22"

# Get default VPC
DEFAULT_VPC_ID=$(aws ec2 describe-vpcs --filters "Name=is-default,Values=true" --query 'Vpcs[0].VpcId' --output text --region "${AWS_REGION}")

# Create/get security group for Cloudflare IPs only
SECURITY_GROUP_ID=$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=osprey-cloudflare-sg" "Name=vpc-id,Values=${DEFAULT_VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' --output text --region "${AWS_REGION}" 2>/dev/null)

if [ "$SECURITY_GROUP_ID" = "None" ] || [ -z "$SECURITY_GROUP_ID" ]; then
    SECURITY_GROUP_ID=$(aws ec2 create-security-group \
        --group-name "osprey-cloudflare-sg" \
        --description "Cloudflare IPs only for App Runner" \
        --vpc-id "${DEFAULT_VPC_ID}" \
        --region "${AWS_REGION}" \
        --query 'GroupId' --output text)
    
    # Add Cloudflare IP rules for HTTPS (443) and HTTP (80)
    IFS=',' read -ra CIDRS <<< "$CLOUDFLARE_CIDRS"
    for cidr in "${CIDRS[@]}"; do
        aws ec2 authorize-security-group-ingress \
            --group-id "${SECURITY_GROUP_ID}" \
            --protocol tcp --port 443 --cidr "${cidr}" \
            --region "${AWS_REGION}" >/dev/null 2>&1 || true
        aws ec2 authorize-security-group-ingress \
            --group-id "${SECURITY_GROUP_ID}" \
            --protocol tcp --port 80 --cidr "${cidr}" \
            --region "${AWS_REGION}" >/dev/null 2>&1 || true
    done
fi

# Get subnets for VPC connector
SUBNET_IDS=$(aws ec2 describe-subnets --filters "Name=vpc-id,Values=${DEFAULT_VPC_ID}" \
    --query 'Subnets[0:2].SubnetId' --output text --region "${AWS_REGION}")

# Create VPC connector if needed for AWS-level network filtering
VPC_CONNECTOR_ARN=$(aws apprunner list-vpc-connectors --region "${AWS_REGION}" \
    --query "VpcConnectors[?VpcConnectorName=='osprey-vpc-connector'].VpcConnectorArn" --output text 2>/dev/null | head -1)

if [ -z "$VPC_CONNECTOR_ARN" ] || [ "$VPC_CONNECTOR_ARN" = "None" ] || [ "$VPC_CONNECTOR_ARN" = "" ]; then
    echo "Creating VPC connector for Cloudflare IP filtering..."
    VPC_CONNECTOR_ARN=$(aws apprunner create-vpc-connector \
        --vpc-connector-name "osprey-vpc-connector" \
        --subnets $SUBNET_IDS \
        --security-groups "${SECURITY_GROUP_ID}" \
        --region "${AWS_REGION}" \
        --query 'VpcConnector.VpcConnectorArn' --output text)
    
    echo "Waiting for VPC connector to be ready..."
    aws apprunner wait vpc-connector-created --vpc-connector-arn "${VPC_CONNECTOR_ARN}" --region "${AWS_REGION}"
fi

# Check if service exists
SERVICE_ARN=$(aws apprunner list-services --region "${AWS_REGION}" \
    --query "ServiceSummaryList[?ServiceName=='${APP_RUNNER_SERVICE}'].ServiceArn" --output text 2>/dev/null | head -1)

if [ -n "$SERVICE_ARN" ] && [ "$SERVICE_ARN" != "None" ] && [ "$SERVICE_ARN" != "" ]; then
    # Update existing service with new image
    aws apprunner update-service \
        --service-arn "${SERVICE_ARN}" \
        --source-configuration ImageRepository="{ImageIdentifier=${FULL_IMAGE_URI},ImageConfiguration={Port=3001}}" \
        --region "${AWS_REGION}" >/dev/null
    
    echo "Updating existing App Runner service..."
    aws apprunner wait service-updated --service-arn "${SERVICE_ARN}" --region "${AWS_REGION}"
else
    # Create new service with VPC connector for AWS-level IP filtering
    echo "Creating new App Runner service with Cloudflare IP filtering..."
    SERVICE_CONFIG=$(cat <<EOF
{
  "ServiceName": "${APP_RUNNER_SERVICE}",
  "SourceConfiguration": {
    "ImageRepository": {
      "ImageIdentifier": "${FULL_IMAGE_URI}",
      "ImageConfiguration": {
        "Port": "3001"
      },
      "ImageRepositoryType": "ECR"
    },
    "AutoDeploymentsEnabled": false
  },
  "InstanceConfiguration": {
    "Cpu": "0.25 vCPU",
    "Memory": "0.5 GB"
  },
  "NetworkConfiguration": {
    "IngressConfiguration": {
      "IsPubliclyAccessible": true
    },
    "EgressConfiguration": {
      "EgressType": "VPC",
      "VpcConnectorArn": "${VPC_CONNECTOR_ARN}"
    }
  }
}
EOF
)
    
    echo "$SERVICE_CONFIG" > /tmp/apprunner-config.json
    SERVICE_ARN=$(aws apprunner create-service \
        --cli-input-json file:///tmp/apprunner-config.json \
        --region "${AWS_REGION}" \
        --query 'Service.ServiceArn' --output text)
    rm /tmp/apprunner-config.json
    
    aws apprunner wait service-created --service-arn "${SERVICE_ARN}" --region "${AWS_REGION}"
fi

# Get service URL
SERVICE_URL=$(aws apprunner describe-service --service-arn "${SERVICE_ARN}" --region "${AWS_REGION}" \
    --query 'Service.ServiceUrl' --output text)

echo "App Runner: https://${SERVICE_URL}" 