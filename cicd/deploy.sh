#!/bin/sh

set -e

# Run Go checks and build executable
go mod tidy
go mod verify
go vet ./...
CGO_ENABLED=0 go build -o main

# Setup SSH key (for Github Workflow env)
echo "$KEY" > ssh_key
chmod 0600 ssh_key

# Deploy to production server (replace binary and restart service)
OPTS="-i ssh_key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
ssh "$OPTS" "$USERNAME"@"$HOST" "sudo systemctl stop website.service"
scp "$OPTS" main "$USERNAME"@"$HOST":/usr/local/bin
ssh "$OPTS" "$USERNAME"@"$HOST" "sudo systemctl start website.service"

# Cleanup (remove created files)
rm -f ssh_key
rm -f main
