#!/bin/sh

set -e

# Run unit tests.
echo "Running unit tests..."
go test ./... -v

# Build the Go executable.
echo "Building Go executable..."
EXE_PATH="temp/main"
export CGO_ENABLED=0
go build -o "$EXE_PATH"

# Setup a SSH identity's private key from environment variable.
echo "Setting up local SSH identify key file..."
SSH_KEY_PATH="temp/ssh_key"
echo "$SSH_KEY" > "$SSH_KEY_PATH"
chmod 0600 "$SSH_KEY_PATH"

# Rollout to remote server (replace binary and restart service).
echo "Rolling out new service version..."
EXE_PATH="/usr/local/bin/website"
SSH_USERNAME="github"
SSH_HOST="localhost"
SSH_PORT=8022
SSH_FLAGS="-i $SSH_KEY_PATH -p $SSH_PORT -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
ssh $SSH_FLAGS "$SSH_USERNAME"@"$SSH_HOST" "cp $EXE_PATH $EXE_PATH.backup"
ssh $SSH_FLAGS "$SSH_USERNAME"@"$SSH_HOST" "sudo systemctl stop website.service"
scp $SSH_FLAGS "$EXE_PATH" "$SSH_USERNAME"@"$SSH_HOST":$EXE_PATH
ssh $SSH_FLAGS "$SSH_USERNAME"@"$SSH_HOST" "sudo systemctl start website.service"

# Done.
echo "✅ Done!"
