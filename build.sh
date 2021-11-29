#!/bin/bash
mkdir -p ./dist
env GOOS=linux go build -o ./dist/proxy
cp ./hosts.json ./dist/
cp ./start.sh ./dist/
echo "Production build complete."