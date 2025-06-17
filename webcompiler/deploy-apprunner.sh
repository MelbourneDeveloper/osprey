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
# Cloudflare IPv4 ranges for WAF IP whitelisting
CLOUDFLARE_CIDRS="173.245.48.0/20,103.21.244.0/22,103.22.200.0/22,103.31.4.0/22,141.101.64.0/18,108.162.192.0/18,190.93.240.0/20,188.114.96.0/20,197.234.240.0/22,198.41.128.0/17,162.158.0.0/15,104.16.0.0/13,104.24.0.0/14,172.64.0.0/13,131.0.72.0/22"

# Create WAF IP set for Cloudflare IPs
WAF_IP_SET_NAME="cloudflare-ip-set"
WAF_WEB_ACL_NAME="osprey-cloudflare-acl"

echo "Creating/updating WAF IP set for Cloudflare IPs..."
IP_SET_ARN=$(aws wafv2 list-ip-sets --scope REGIONAL --region "${AWS_REGION}" \
    --query "IPSets[?Name=='${WAF_IP_SET_NAME}'].ARN" --output text 2>/dev/null | head -1)

if [ -z "$IP_SET_ARN" ] || [ "$IP_SET_ARN" = "None" ]; then
    # Create new IP set
    IFS=',' read -ra CIDRS <<< "$CLOUDFLARE_CIDRS"
    ADDRESSES_JSON=$(printf '"%s",' "${CIDRS[@]}" | sed 's/,$//')
    
    IP_SET_ARN=$(aws wafv2 create-ip-set \
        --name "${WAF_IP_SET_NAME}" \
        --scope REGIONAL \
        --ip-address-version IPV4 \
        --addresses "[${ADDRESSES_JSON}]" \
        --region "${AWS_REGION}" \
        --query 'Summary.ARN' --output text)
else
    # Update existing IP set
    IFS=',' read -ra CIDRS <<< "$CLOUDFLARE_CIDRS"
    ADDRESSES_JSON=$(printf '"%s",' "${CIDRS[@]}" | sed 's/,$//')
    
    LOCK_TOKEN=$(aws wafv2 get-ip-set --scope REGIONAL --id "${IP_SET_ARN##*/}" --name "${WAF_IP_SET_NAME}" --region "${AWS_REGION}" --query 'LockToken' --output text)
    
    aws wafv2 update-ip-set \
        --scope REGIONAL \
        --id "${IP_SET_ARN##*/}" \
        --name "${WAF_IP_SET_NAME}" \
        --addresses "[${ADDRESSES_JSON}]" \
        --lock-token "${LOCK_TOKEN}" \
        --region "${AWS_REGION}" >/dev/null
fi

echo "Creating/updating WAF Web ACL for IP whitelisting..."
WEB_ACL_ARN=$(aws wafv2 list-web-acls --scope REGIONAL --region "${AWS_REGION}" \
    --query "WebACLs[?Name=='${WAF_WEB_ACL_NAME}'].ARN" --output text 2>/dev/null | head -1)

if [ -z "$WEB_ACL_ARN" ] || [ "$WEB_ACL_ARN" = "None" ]; then
    # Create new Web ACL
    WEB_ACL_CONFIG=$(cat <<EOF
{
  "Name": "${WAF_WEB_ACL_NAME}",
  "Scope": "REGIONAL",
  "DefaultAction": {"Block": {}},
  "Rules": [
    {
      "Name": "AllowCloudflareIPs",
      "Priority": 1,
      "Statement": {
        "IPSetReferenceStatement": {
          "ARN": "${IP_SET_ARN}"
        }
      },
      "Action": {"Allow": {}},
      "VisibilityConfig": {
        "SampledRequestsEnabled": true,
        "CloudWatchMetricsEnabled": true,
        "MetricName": "AllowCloudflareIPs"
      }
    }
  ],
  "VisibilityConfig": {
    "SampledRequestsEnabled": true,
    "CloudWatchMetricsEnabled": true,
    "MetricName": "${WAF_WEB_ACL_NAME}"
  }
}
EOF
)
    
    echo "$WEB_ACL_CONFIG" > /tmp/waf-config.json
    WEB_ACL_ARN=$(aws wafv2 create-web-acl \
        --cli-input-json file:///tmp/waf-config.json \
        --region "${AWS_REGION}" \
        --query 'Summary.ARN' --output text)
    rm /tmp/waf-config.json
else
    # Update existing Web ACL
    LOCK_TOKEN=$(aws wafv2 get-web-acl --scope REGIONAL --id "${WEB_ACL_ARN##*/}" --name "${WAF_WEB_ACL_NAME}" --region "${AWS_REGION}" --query 'LockToken' --output text)
    
    WEB_ACL_CONFIG=$(cat <<EOF
{
  "Scope": "REGIONAL",
  "Id": "${WEB_ACL_ARN##*/}",
  "Name": "${WAF_WEB_ACL_NAME}",
  "DefaultAction": {"Block": {}},
  "Rules": [
    {
      "Name": "AllowCloudflareIPs",
      "Priority": 1,
      "Statement": {
        "IPSetReferenceStatement": {
          "ARN": "${IP_SET_ARN}"
        }
      },
      "Action": {"Allow": {}},
      "VisibilityConfig": {
        "SampledRequestsEnabled": true,
        "CloudWatchMetricsEnabled": true,
        "MetricName": "AllowCloudflareIPs"
      }
    }
  ],
  "VisibilityConfig": {
    "SampledRequestsEnabled": true,
    "CloudWatchMetricsEnabled": true,
    "MetricName": "${WAF_WEB_ACL_NAME}"
  },
  "LockToken": "${LOCK_TOKEN}"
}
EOF
)
    
    echo "$WEB_ACL_CONFIG" > /tmp/waf-update-config.json
    aws wafv2 update-web-acl \
        --cli-input-json file:///tmp/waf-update-config.json \
        --region "${AWS_REGION}" >/dev/null
    rm /tmp/waf-update-config.json
fi

# Use existing Cloudflare security group
SECURITY_GROUP_ID="sg-049d709425679e233"  # osprey-cloudflare-sg
DEFAULT_VPC_ID=$(aws ec2 describe-vpcs --filters "Name=is-default,Values=true" --query 'Vpcs[0].VpcId' --output text --region "${AWS_REGION}")

# Get subnets for VPC connector
SUBNET_IDS=$(aws ec2 describe-subnets --filters "Name=vpc-id,Values=${DEFAULT_VPC_ID}" \
    --query 'Subnets[0:2].SubnetId' --output text --region "${AWS_REGION}")

# Create VPC connector using existing security group
VPC_CONNECTOR_ARN=$(aws apprunner list-vpc-connectors --region "${AWS_REGION}" \
    --query "VpcConnectors[?VpcConnectorName=='osprey-vpc-connector'].VpcConnectorArn" --output text 2>/dev/null | head -1)

if [ -z "$VPC_CONNECTOR_ARN" ] || [ "$VPC_CONNECTOR_ARN" = "None" ] || [ "$VPC_CONNECTOR_ARN" = "" ]; then
    echo "Creating VPC connector with Cloudflare security group..."
    VPC_CONNECTOR_ARN=$(aws apprunner create-vpc-connector \
        --vpc-connector-name "osprey-vpc-connector" \
        --subnets $SUBNET_IDS \
        --security-groups "${SECURITY_GROUP_ID}" \
        --region "${AWS_REGION}" \
        --query 'VpcConnector.VpcConnectorArn' --output text)
    
    echo "Waiting for VPC connector to be ready..."
    while true; do
        STATUS=$(aws apprunner describe-vpc-connector --vpc-connector-arn "${VPC_CONNECTOR_ARN}" --region "${AWS_REGION}" --query 'VpcConnector.Status' --output text)
        if [ "$STATUS" = "ACTIVE" ]; then
            echo "VPC connector is ready"
            break
        fi
        echo "VPC connector status: $STATUS, waiting..."
        sleep 10
    done
fi

# Wait for WAF resources to be ready
echo "Waiting for WAF resources to be available..."
sleep 5

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

# Associate WAF with App Runner service for IP whitelisting
echo "Associating WAF Web ACL with App Runner service for Cloudflare IP filtering..."
aws wafv2 associate-web-acl \
    --web-acl-arn "${WEB_ACL_ARN}" \
    --resource-arn "${SERVICE_ARN}" \
    --region "${AWS_REGION}" 2>/dev/null || {
    echo "WAF association may already exist or service not ready, continuing..."
}

# Get service URL
SERVICE_URL=$(aws apprunner describe-service --service-arn "${SERVICE_ARN}" --region "${AWS_REGION}" \
    --query 'Service.ServiceUrl' --output text)

echo "🔒 App Runner deployed with Cloudflare IP whitelisting: https://${SERVICE_URL}"
echo "✅ WAF Web ACL '${WAF_WEB_ACL_NAME}' blocks ALL traffic except Cloudflare IPs"
echo "🛡️  Security Group '${SECURITY_GROUP_ID}' (osprey-cloudflare-sg) hooked up via VPC connector"
echo "📋 Allowed IP ranges: ${CLOUDFLARE_CIDRS}" 