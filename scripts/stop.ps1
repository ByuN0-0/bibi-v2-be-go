# Parse command line arguments
param(
    [switch]$Remove,
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\stop.ps1 [-Remove] [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Remove       Stop and remove the container (default: only stop)"
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

Write-Host "${Yellow}=== Bibi Bot Stop Script ===${NC}"

# Validate environment
if ([string]::IsNullOrEmpty($PROJECT_ID)) {
    Write-Host "${Red}Error: GCP_PROJECT_ID is not set${NC}" -ForegroundColor Red
    exit 1
}

# Step 1: Stop container
Write-Host "`n${Yellow}[1/1] Stopping bot container on VM...${NC}"

if ($Remove) {
    $STOP_COMMANDS = @'
#!/bin/bash

echo "Stopping and removing container..."
sudo docker stop bibi-bot-container 2>/dev/null || echo "Container not running"
sudo docker rm bibi-bot-container 2>/dev/null || echo "Container not found"

echo "Container stopped and removed"
sudo docker ps -a | grep bibi-bot || echo "No bibi-bot containers found"
'@
    Write-Host "${Yellow}Mode: Stop and remove container${NC}" -ForegroundColor Yellow
} else {
    $STOP_COMMANDS = @'
#!/bin/bash

echo "Stopping container..."
sudo docker stop bibi-bot-container 2>/dev/null || echo "Container not running"

echo "Container stopped"
sudo docker ps -a | grep bibi-bot || echo "No bibi-bot containers found"
'@
    Write-Host "${Yellow}Mode: Stop only (use -Remove to also remove the container)${NC}" -ForegroundColor Yellow
}

gcloud compute ssh "$INSTANCE_NAME" `
    --zone="$ZONE" `
    --project="$PROJECT_ID" `
    --ssh-key-file="$HOME/.ssh/google_compute_engine" `
    --command="$STOP_COMMANDS"

if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ SSH connection or stop commands failed${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "`n${Green}=== Bot stopped successfully ===${NC}" -ForegroundColor Green

if (!$Remove) {
    Write-Host "${Yellow}Note: Container is stopped but not removed. Use -Remove to also remove the container.${NC}" -ForegroundColor Yellow
}

Write-Host "`n${Yellow}Connect to VM to check status:${NC}" -ForegroundColor Yellow
Write-Host "gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID"
