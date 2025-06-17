#!/bin/bash

set -e
cd "$(dirname "$0")"

echo "Deploying Osprey Web Compiler..."

./deploy-ecr.sh
./deploy-apprunner.sh

echo "Deployment complete!" 