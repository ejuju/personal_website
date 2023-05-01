#!/bin/sh

# Run Go checks and build executable
go mod tidy
go mod verify
go vet ./...
CGO_ENABLED=0 go build -o personal_website

# Setup SSH key (for Github Workflow env)
if test -z "$KEY"; then
	echo "$KEY" > .ssh_key
	chmod 0600 .ssh_key
fi

# Deploy to production server (replace binary and restart service)
scp \
	-i .ssh_key \
	-o StrictHostKeyChecking=no \
	-o UserKnownHostsFile=/dev/null \
	personal_website \
	"$USERNAME"@"$HOST":/usr/local/bin/
ssh \
	-i .ssh_key \
	-o StrictHostKeyChecking=no \
	-o UserKnownHostsFile=/dev/null \
	"$USERNAME"@"$HOST" \
	"sudo systemctl restart personal_website"