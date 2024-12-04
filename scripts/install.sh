#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Installing Simple Monitor...${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${BLUE}Installing Docker...${NC}"
    curl -fsSL https://get.docker.com | sh
fi

# Create directory for data
sudo mkdir -p /opt/simple-monitor/data

# Create docker-compose.yml
cat > /opt/simple-monitor/docker-compose.yml << EOL
version: '3.8'

services:
  monitor:
    image: yourdockerhubusername/simple-monitor:latest
    container_name: simple-monitor
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    restart: unless-stopped
    privileged: true
EOL

# Change to installation directory
cd /opt/simple-monitor

# Pull and start the container
docker-compose up -d

echo -e "${GREEN}Installation complete!${NC}"
echo -e "Monitor is running on port 8080"
echo -e "\nTest with: curl http://localhost:8080/metrics"
echo -e "\nTo update in the future, run:"
echo -e "cd /opt/simple-monitor && docker-compose pull && docker-compose up -d"