# Parse command line arguments
param(
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\start.ps1 [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Help         Show this help message"
    exit 0
}

# Configuration
$PROJECT_ID = if ($env:GCP_PROJECT_ID) { $env:GCP_PROJECT_ID } else { "by-251105" }
$ZONE = if ($env:GCP_ZONE) { $env:GCP_ZONE } else { "us-west1-b" }
$INSTANCE_NAME = if ($env:INSTANCE_NAME) { $env:INSTANCE_NAME } else { "bibi-bot-vm" }
$CONTAINER_NAME = "bibi-bot-container"

# Color codes for output
$Yellow = "`e[1;33m"
$Green = "`e[0;32m"
$Red = "`e[0;31m"
$NC = "`e[0m"

Write-Host "${Yellow}=== Bibi Bot Start Script ===${NC}"

# Validate environment
if ([string]::IsNullOrEmpty($PROJECT_ID)) {
    Write-Host "${Red}Error: GCP_PROJECT_ID is not set${NC}" -ForegroundColor Red
    exit 1
}

# Step 1: Start container
Write-Host "`n${Yellow}[1/1] Starting bot container on VM...${NC}"

$START_COMMANDS = @'
#!/bin/bash

# Check if container exists but is stopped
if sudo docker ps -a | grep -q bibi-bot-container; then
    echo "Found existing container, starting..."
    sudo docker start bibi-bot-container

    if [ $? -eq 0 ]; then
        echo "Container started successfully"
        sudo docker ps | grep bibi-bot
    else
        echo "Error: Failed to start existing container"
        exit 1
    fi
else
    echo "No existing container found. Starting new container..."

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

    if [ $? -eq 0 ]; then
        echo "Container started successfully"
        sudo docker ps | grep bibi-bot
    else
        echo "Error: Failed to start container"
        exit 1
    fi
fi
'@

gcloud compute ssh "$INSTANCE_NAME" `
    --zone="$ZONE" `
    --project="$PROJECT_ID" `
    --ssh-key-file="$HOME/.ssh/google_compute_engine" `
    --command="$START_COMMANDS"

if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ SSH connection or start commands failed${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "`n${Green}=== Bot started successfully ===${NC}" -ForegroundColor Green
Write-Host "`n${Yellow}Connect to VM to check logs:${NC}" -ForegroundColor Yellow
Write-Host "gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID"
Write-Host "${Yellow}View container logs:${NC}" -ForegroundColor Yellow
Write-Host "sudo docker logs -f bibi-bot-container"
