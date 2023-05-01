#!/bin/sh

set -e

# Run Go checks and build executable
go mod tidy
go mod verify
go vet ./...
CGO_ENABLED=0 go build -o personal_website

# Setup SSH key (for Github Workflow env)
echo "$KEY" > ssh_key
chmod 0600 ssh_key

# Deploy to production server (replace binary and restart service)
OPTS="-i ssh_key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
ssh $OPTS "$USERNAME"@"$HOST" "systemctl stop personal_website"
scp $OPTS "$USERNAME"@"$HOST" personal_website "$USERNAME"@"$HOST":/usr/local/bin
ssh $OPTS "$USERNAME"@"$HOST" "systemctl restart personal_website"

# Cleanup (remove created files)
rm -f ssh_key
rm -f personal_website