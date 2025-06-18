#!/bin/bash

# Configuration - App Runner cheapest in us-east-1
AWS_REGION="us-east-1"
ECR_REPOSITORY="osprey-web-compiler"
IMAGE_TAG="latest"
APP_RUNNER_SERVICE="osprey-web-compiler"

# Load sensitive config from .env file if it exists
if [ -f "$(dirname "$0")/.env" ]; then
    set -a
    source "$(dirname "$0")/.env"
    set +a
    echo "Loaded config from .env file"
fi

# Get AWS account ID (from .env or dynamically)
get_aws_account_id() {
    if [ -z "$AWS_ACCOUNT_ID" ]; then
        AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
        echo "Got account ID dynamically: ${AWS_ACCOUNT_ID}"
    else
        echo "Using account ID from .env: ${AWS_ACCOUNT_ID}"
    fi
    ECR_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
    FULL_IMAGE_URI="${ECR_URI}/${ECR_REPOSITORY}:${IMAGE_TAG}"
} 