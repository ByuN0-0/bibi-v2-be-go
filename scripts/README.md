# Bibi Bot Management Scripts

PowerShell scripts for deploying and managing Bibi Bot on Google Cloud Platform.

## Prerequisites

- PowerShell (Windows, macOS, or Linux)
- Google Cloud SDK (`gcloud`) installed and configured
- Access to the GCP project `by-251105`

## Environment Variables (Optional)

You can override default values with these environment variables:

```powershell
$env:GCP_PROJECT_ID = "by-251105"         # Default: by-251105
$env:GCP_REGION = "us-west1"              # Default: us-west1
$env:GCP_ZONE = "us-west1-b"              # Default: us-west1-b
$env:INSTANCE_NAME = "bibi-bot-vm"        # Default: bibi-bot-vm
```

## Scripts

### deploy.ps1
Deploy the bot to GCP (build, push, and start on VM)

```powershell
# Full deployment (build + push + deploy)
.\deploy.ps1

# Deploy without rebuilding Docker image
.\deploy.ps1 -SkipBuild

# Show help
.\deploy.ps1 -Help
```

**What it does:**
1. Builds Docker image
2. Authenticates with Artifact Registry
3. Pushes image to registry
4. Uploads `.env` file to VM
5. Stops old container
6. Pulls new image
7. Starts new container

---

### start.ps1
Start the bot on GCP VM

```powershell
# Start the bot
.\start.ps1

# Show help
.\start.ps1 -Help
```

**What it does:**
- If container exists but is stopped: starts it
- If container doesn't exist: creates and starts new container with latest image

---

### stop.ps1
Stop the running bot on GCP VM

```powershell
# Stop the container (keeps it for restart)
.\stop.ps1

# Stop and remove the container
.\stop.ps1 -Remove

# Show help
.\stop.ps1 -Help
```

**What it does:**
- `-Remove` flag: stops and removes the container completely
- Without flag: only stops the container (can be restarted with `start.ps1`)

---

### status.ps1
Check the bot's current status

```powershell
# Check status
.\status.ps1

# Show help
.\status.ps1 -Help
```

**What it shows:**
- Container running/stopped/not found status
- Container details (ID, image, uptime)
- Container resource usage (CPU, memory)
- Recent logs (last 10 lines)

---

### logs.ps1
View bot logs

```powershell
# View last 100 lines (default)
.\logs.ps1

# Follow logs in real-time
.\logs.ps1 -Follow

# View last 50 lines
.\logs.ps1 -Tail 50

# Follow with custom tail
.\logs.ps1 -Follow -Tail 20

# Show help
.\logs.ps1 -Help
```

**Options:**
- `-Follow`: Stream logs in real-time (like `tail -f`)
- `-Tail <number>`: Number of lines to show (default: 100)

---

## Common Workflows

### First Deployment
```powershell
.\deploy.ps1
```

### Update Bot Code
```powershell
# Make code changes, then:
.\deploy.ps1
```

### Restart Bot
```powershell
.\stop.ps1
.\start.ps1
```

### Quick Check
```powershell
.\status.ps1
```

### Debug Issues
```powershell
# View recent logs
.\logs.ps1

# Follow logs in real-time
.\logs.ps1 -Follow
```

### Complete Redeployment
```powershell
# Stop and remove old container
.\stop.ps1 -Remove

# Deploy fresh
.\deploy.ps1
```

---

## Troubleshooting

### Script fails with "GCP_PROJECT_ID is not set"
Set the environment variable:
```powershell
$env:GCP_PROJECT_ID = "by-251105"
```

### SSH connection fails
1. Check if VM is running: `gcloud compute instances list`
2. Check your SSH keys: `gcloud compute ssh bibi-bot-vm --zone=us-west1-b`

### Container not starting
1. Check logs: `.\logs.ps1`
2. Check `.env` file exists with valid `DISCORD_BOT_TOKEN`
3. SSH into VM and check manually:
   ```bash
   gcloud compute ssh bibi-bot-vm --zone=us-west1-b
   sudo docker logs bibi-bot-container
   ```

### Docker image pull fails
1. Check Artifact Registry authentication
2. Verify image exists:
   ```bash
   gcloud artifacts docker images list us-west1-docker.pkg.dev/by-251105/bibi-bot
   ```
