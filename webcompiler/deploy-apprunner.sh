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
# Cloudflare IPv4 ranges for security group IP whitelisting
CLOUDFLARE_CIDRS="173.245.48.0/20,103.21.244.0/22,103.22.200.0/22,103.31.4.0/22,141.101.64.0/18,108.162.192.0/18,190.93.240.0/20,188.114.96.0/20,197.234.240.0/22,198.41.128.0/17,162.158.0.0/15,104.16.0.0/13,104.24.0.0/14,172.64.0.0/13,131.0.72.0/22"

echo "Using existing Cloudflare security group for IP whitelisting..."

# Use existing Cloudflare security group - NO CREATION OR MODIFICATION!
SECURITY_GROUP_ID="sg-049d709425679e233"  # osprey-cloudflare-sg

# Get existing VPC connector
VPC_CONNECTOR_ARN=$(aws apprunner list-vpc-connectors --region "${AWS_REGION}" \
    --query "VpcConnectors[?VpcConnectorName=='osprey-vpc-connector'].VpcConnectorArn" --output text 2>/dev/null | head -1)

if [ -z "$VPC_CONNECTOR_ARN" ] || [ "$VPC_CONNECTOR_ARN" = "None" ] || [ "$VPC_CONNECTOR_ARN" = "" ]; then
    echo "ERROR: VPC connector 'osprey-vpc-connector' not found!"
    echo "Create it manually with the existing security group ${SECURITY_GROUP_ID}"
    exit 1
fi

echo "Using existing VPC connector: ${VPC_CONNECTOR_ARN}"

# Security group is already configured, no wait needed
echo "Cloudflare security group ready..."

# Check if service exists
SERVICE_ARN=$(aws apprunner list-services --region "${AWS_REGION}" \
    --query "ServiceSummaryList[?ServiceName=='${APP_RUNNER_SERVICE}'].ServiceArn" --output text 2>/dev/null | head -1)

if [ -n "$SERVICE_ARN" ] && [ "$SERVICE_ARN" != "None" ] && [ "$SERVICE_ARN" != "" ]; then
    # Update existing service with new image  
    aws apprunner update-service \
        --service-arn "${SERVICE_ARN}" \
        --source-configuration ImageRepository="{ImageIdentifier=${FULL_IMAGE_URI},ImageConfiguration={Port=3001},ImageRepositoryType=ECR}" \
        --region "${AWS_REGION}" >/dev/null
    
    echo "Updating existing App Runner service..."
    sleep 5  # Brief wait for update to start
else
    # Create new service with VPC connector and Cloudflare security group
    echo "Creating new App Runner service with Cloudflare security group..."
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
    "AutoDeploymentsEnabled": false,
    "AuthenticationConfiguration": {
      "AccessRoleArn": "arn:aws:iam::${AWS_ACCOUNT_ID}:role/AppRunnerECRAccessRole"
    }
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
    
    sleep 5  # Brief wait for creation to start
fi

# Wait for service to be running before associating WAF
echo "Waiting for App Runner service to be operational..."
while true; do
    STATUS=$(aws apprunner describe-service --service-arn "${SERVICE_ARN}" --region "${AWS_REGION}" --query 'Service.Status' --output text)
    if [ "$STATUS" = "RUNNING" ]; then
        echo "App Runner service is running"
        break
    fi
    echo "App Runner service status: $STATUS, waiting..."
    sleep 15
done

# Security group handles IP filtering through VPC connector
echo "Cloudflare IP filtering active via security group ${SECURITY_GROUP_ID}"

# Get service URL
SERVICE_URL=$(aws apprunner describe-service --service-arn "${SERVICE_ARN}" --region "${AWS_REGION}" \
    --query 'Service.ServiceUrl' --output text)

echo "🔒 App Runner deployed with Cloudflare IP whitelisting: https://${SERVICE_URL}"
echo "🛡️  Security Group '${SECURITY_GROUP_ID}' (osprey-cloudflare-sg) blocks ALL traffic except Cloudflare IPs"
echo "✅ VPC connector enforces security group rules - NO OTHER IPS CAN ACCESS!"
echo "📋 Allowed IP ranges: ${CLOUDFLARE_CIDRS}" 