# Parse command line arguments
param(
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\status.ps1 [-Help]"
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

Write-Host "${Yellow}=== Bibi Bot Status ===${NC}"

# Validate environment
if ([string]::IsNullOrEmpty($PROJECT_ID)) {
    Write-Host "${Red}Error: GCP_PROJECT_ID is not set${NC}" -ForegroundColor Red
    exit 1
}

# Get container status
Write-Host "`n${Yellow}Checking container status...${NC}`n"

$STATUS_COMMANDS = @'
#!/bin/bash

echo "=== Container Status ==="
if sudo docker ps | grep -q bibi-bot-container; then
    echo "Status: RUNNING"
    echo ""
    echo "Container Details:"
    sudo docker ps | grep bibi-bot
    echo ""
    echo "Container Stats:"
    sudo docker stats --no-stream bibi-bot-container
elif sudo docker ps -a | grep -q bibi-bot-container; then
    echo "Status: STOPPED"
    echo ""
    echo "Container Details:"
    sudo docker ps -a | grep bibi-bot
else
    echo "Status: NOT FOUND"
    echo "No bibi-bot container exists on this VM"
fi

echo ""
echo "=== Recent Logs (last 10 lines) ==="
sudo docker logs --tail 10 bibi-bot-container 2>/dev/null || echo "No logs available"
'@

gcloud compute ssh "$INSTANCE_NAME" `
    --zone="$ZONE" `
    --project="$PROJECT_ID" `
    --ssh-key-file="$HOME/.ssh/google_compute_engine" `
    --command="$STATUS_COMMANDS"

if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ Failed to retrieve status${NC}" -ForegroundColor Red
    exit 1
}

Write-Host "`n${Green}Status check completed${NC}" -ForegroundColor Green
