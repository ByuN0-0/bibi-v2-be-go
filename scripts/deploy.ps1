# Parse command line arguments
param(
    [switch]$SkipBuild,
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\deploy.ps1 [-SkipBuild] [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -SkipBuild    Skip Docker build and push steps, only update VM"
    Write-Host "  -Help         Show this help message"
    exit 0
}

# Configuration
$PROJECT_ID = if ($env:GCP_PROJECT_ID) { $env:GCP_PROJECT_ID } else { "by-251105" }
$REGION = if ($env:GCP_REGION) { $env:GCP_REGION } else { "us-west1" }
$ZONE = if ($env:GCP_ZONE) { $env:GCP_ZONE } else { "us-west1-b" }
$INSTANCE_NAME = if ($env:INSTANCE_NAME) { $env:INSTANCE_NAME } else { "bibi-bot-vm" }
$IMAGE_NAME = "bibi-bot"
$CONTAINER_NAME = "bibi-bot-container"
$REGISTRY = "${REGION}-docker.pkg.dev/${PROJECT_ID}/bibi-bot"

# Color codes for output
$Yellow = "`e[1;33m"
$Green = "`e[0;32m"
$Red = "`e[0;31m"
$NC = "`e[0m"

Write-Host "${Yellow}=== Bibi Bot Deployment Script ===${NC}"

# Validate environment
if ([string]::IsNullOrEmpty($PROJECT_ID)) {
    Write-Host "${Red}Error: GCP_PROJECT_ID is not set${NC}" -ForegroundColor Red
    exit 1
}

# Get the directory of this script
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$PROJECT_ROOT = Split-Path -Parent $SCRIPT_DIR

if ($SkipBuild) {
    Write-Host "`n${Yellow}Skipping Docker build and push (using existing image in GCR)${NC}" -ForegroundColor Yellow
} else {
    # Step 1: Build Docker image
    Write-Host "`n${Yellow}[1/5] Building Docker image...${NC}"
docker build -t "${IMAGE_NAME}:latest" -f "${PROJECT_ROOT}/docker/Dockerfile" "${PROJECT_ROOT}"
if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ Docker build failed${NC}" -ForegroundColor Red
    exit 1
}

docker tag "${IMAGE_NAME}:latest" "${REGISTRY}/${IMAGE_NAME}:latest"
docker tag "${IMAGE_NAME}:latest" "${REGISTRY}/${IMAGE_NAME}:$(Get-Date -Format 'yyyyMMdd_HHmmss')"

Write-Host "${Green}✓ Docker image built successfully${NC}" -ForegroundColor Green

# Step 2: Authenticate with Artifact Registry
Write-Host "`n${Yellow}[2/5] Authenticating with Artifact Registry...${NC}"
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet
if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ Authentication failed${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "${Green}✓ Authentication successful${NC}" -ForegroundColor Green

# Step 3: Push image to Artifact Registry
Write-Host "`n${Yellow}[3/5] Pushing image to Artifact Registry...${NC}"
docker push "${REGISTRY}/${IMAGE_NAME}:latest"
if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ Artifact Registry push failed${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "${Green}✓ Image pushed to Artifact Registry${NC}" -ForegroundColor Green
}

# Step 3.5: Upload .env file to VM
Write-Host "`n${Yellow}[3.5/6] Uploading environment variables to VM...${NC}"

$ENV_FILE = "${PROJECT_ROOT}/.env"
if (Test-Path $ENV_FILE) {
    # Upload .env to home directory (without ~ which doesn't work in Windows gcloud)
    gcloud compute scp "$ENV_FILE" "${INSTANCE_NAME}:.bibi-bot.env" --zone="$ZONE" --project="$PROJECT_ID"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "${Red}✗ Failed to upload .env file${NC}" -ForegroundColor Red
        exit 1
    }

    # Set permissions
    gcloud compute ssh "$INSTANCE_NAME" `
        --zone="$ZONE" `
        --project="$PROJECT_ID" `
        --command="chmod 600 .bibi-bot.env"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "${Red}✗ Failed to set permissions on .env file${NC}" -ForegroundColor Red
        exit 1
    }

    Write-Host "${Green}✓ Environment variables uploaded to VM${NC}" -ForegroundColor Green
} else {
    Write-Host "${Yellow}Warning: .env file not found at $ENV_FILE${NC}" -ForegroundColor Yellow
    Write-Host "${Yellow}Container will start without environment variables${NC}" -ForegroundColor Yellow
}

# Step 4: SSH into VM and stop old container
Write-Host "`n${Yellow}[4/6] Updating container on VM...${NC}"

$REMOTE_COMMANDS = @'
#!/bin/bash
set -e

# Stop and remove old container
echo "Stopping old container..."
sudo docker stop bibi-bot-container 2>/dev/null || true
sudo docker rm bibi-bot-container 2>/dev/null || true

# Pull new image using metadata service to authenticate
echo "Pulling new image from Artifact Registry..."
TOKEN=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')

# Use /tmp for Docker config (writable in COS)
mkdir -p /tmp/.docker
export DOCKER_CONFIG=/tmp/.docker
echo $TOKEN | sudo -E docker login -u oauth2accesstoken --password-stdin us-west1-docker.pkg.dev 2>&1 | grep -v WARNING
sudo -E docker pull us-west1-docker.pkg.dev/by-251105/bibi-bot/bibi-bot:latest

echo "Container updated successfully"
'@

gcloud compute ssh "$INSTANCE_NAME" `
    --zone="$ZONE" `
    --project="$PROJECT_ID" `
    --ssh-key-file="$HOME/.ssh/google_compute_engine" `
    --command="$REMOTE_COMMANDS"

if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ SSH connection or remote commands failed${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "${Green}✓ Container update commands sent to VM${NC}" -ForegroundColor Green

# Step 5: Start new container on VM
Write-Host "`n${Yellow}[5/6] Starting new container on VM...${NC}"

$START_CONTAINER = @'
#!/bin/bash

# Load environment variables from home directory
if [ -f .bibi-bot.env ]; then
    export $(cat .bibi-bot.env | grep -v '^#' | grep -v '^$' | xargs)
else
    echo "Warning: .bibi-bot.env not found. Container will start without env vars."
fi

# Run the new container
sudo docker run -d \
    --name bibi-bot-container \
    --restart unless-stopped \
    -e DISCORD_BOT_TOKEN="${DISCORD_BOT_TOKEN}" \
    -e DISCORD_CLIENT_ID="${DISCORD_CLIENT_ID}" \
    -e DISCORD_CLIENT_SECRET="${DISCORD_CLIENT_SECRET}" \
    -e DISCORD_PUBLIC_KEY="${DISCORD_PUBLIC_KEY}" \
    us-west1-docker.pkg.dev/by-251105/bibi-bot/bibi-bot:latest

echo "Container started successfully"
sudo docker ps | grep bibi-bot
'@

gcloud compute ssh "$INSTANCE_NAME" `
    --zone="$ZONE" `
    --project="$PROJECT_ID" `
    --ssh-key-file="$HOME/.ssh/google_compute_engine" `
    --command="$START_CONTAINER"

if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ Container start failed${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "`n${Green}=== Deployment completed successfully ===${NC}" -ForegroundColor Green
Write-Host "${Yellow}Connect to VM:${NC}" -ForegroundColor Yellow
Write-Host "gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID"
