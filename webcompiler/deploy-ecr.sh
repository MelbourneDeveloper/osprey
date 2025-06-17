#!/bin/bash

set -e
source "$(dirname "$0")/config.sh"
cd "$(dirname "$0")"

# Verify AWS credentials
aws sts get-caller-identity >/dev/null || { echo "AWS credentials not configured"; exit 1; }

# Verify Docker is available
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed"; exit 1; }

# Get AWS account details now that credentials are verified
get_aws_account_id

# Create ECR repository if needed
if ! aws ecr describe-repositories --repository-names "${ECR_REPOSITORY}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    aws ecr create-repository \
        --repository-name "${ECR_REPOSITORY}" \
        --region "${AWS_REGION}" \
        --image-scanning-configuration scanOnPush=true >/dev/null
fi

# Login and build
aws ecr get-login-password --region "${AWS_REGION}" | docker login --username AWS --password-stdin "${ECR_URI}"
docker build --platform linux/amd64 --no-cache -t "${FULL_IMAGE_URI}" -f Dockerfile ..
docker push "${FULL_IMAGE_URI}"

echo "ECR: ${FULL_IMAGE_URI}" 