#!/bin/zsh

make docker-registry

# to pull from registry and update image/container 
ssh -t "$PI_USER"@"$PI_IP" "./update.sh"
