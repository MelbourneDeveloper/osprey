#!/bin/bash

# Configuration - App Runner cheapest in us-east-1
AWS_REGION="us-east-1"
ECR_REPOSITORY="osprey-web-compiler"
IMAGE_TAG="latest"
APP_RUNNER_SERVICE="osprey-web-compiler"

# Get AWS account ID dynamically (will be set after credential check)
get_aws_account_id() {
    AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
    ECR_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
    FULL_IMAGE_URI="${ECR_URI}/${ECR_REPOSITORY}:${IMAGE_TAG}"
} 