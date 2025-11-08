# Parse command line arguments
param(
    [switch]$Follow,
    [int]$Tail = 100,
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\logs.ps1 [-Follow] [-Tail <number>] [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Follow       Follow log output (like tail -f)"
    Write-Host "  -Tail         Number of lines to show (default: 100)"
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

Write-Host "${Yellow}=== Bibi Bot Logs ===${NC}"

# Validate environment
if ([string]::IsNullOrEmpty($PROJECT_ID)) {
    Write-Host "${Red}Error: GCP_PROJECT_ID is not set${NC}" -ForegroundColor Red
    exit 1
}

# Build log command
if ($Follow) {
    $LOG_COMMAND = "sudo docker logs -f --tail $Tail $CONTAINER_NAME"
    Write-Host "${Yellow}Following logs (Ctrl+C to stop)...${NC}`n" -ForegroundColor Yellow
} else {
    $LOG_COMMAND = "sudo docker logs --tail $Tail $CONTAINER_NAME"
    Write-Host "${Yellow}Showing last $Tail lines...${NC}`n" -ForegroundColor Yellow
}

# Execute SSH command
gcloud compute ssh "$INSTANCE_NAME" `
    --zone="$ZONE" `
    --project="$PROJECT_ID" `
    --ssh-key-file="$HOME/.ssh/google_compute_engine" `
    --command="$LOG_COMMAND"

if ($LASTEXITCODE -ne 0) {
    Write-Host "${Red}✗ Failed to retrieve logs${NC}" -ForegroundColor Red
    exit 1
}
