#!/bin/bash

# Colors for pretty output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print styled messages
print_header() { echo -e "\n${BLUE}=== $1 ===${NC}\n"; }
print_success() { echo -e "${GREEN}✔ $1${NC}"; }
print_error() { echo -e "${RED}✘ $1${NC}"; }
print_warning() { echo -e "${YELLOW}! $1${NC}"; }

# Get server's public IP
get_public_ip() {
    curl -s ifconfig.me || curl -s icanhazip.com || curl -s ipinfo.io/ip
}

# Create project directory
PROJECT_DIR="simple-monitor"
print_header "Creating Project Directory"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"
print_success "Created and moved to directory: $PROJECT_DIR"

# Create docker-compose.yml
print_header "Creating Docker Compose Configuration"
cat > docker-compose.yml << 'EOF'
services:
  monitor:
    image: kaanmertkoc1/simple-monitor:latest
    container_name: simple-monitor
    expose:
      - "8080"
    volumes:
      - ./data:/data
    restart: unless-stopped
    privileged: true
    networks:
      - monitor-network

  nginx:
    image: nginx:alpine
    container_name: monitor-nginx
    ports:
      - "9090:80"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d
    networks:
      - monitor-network
    depends_on:
      - monitor

networks:
  monitor-network:
    driver: bridge
EOF
print_success "Created docker-compose.yml"

# Create required directories
print_header "Creating Directory Structure"
mkdir -p nginx/conf.d data
print_success "Created required directories"

# Generate nginx config
print_header "Configuring Nginx"
cat > nginx/conf.d/app.conf << EOF
server {
    listen 80;
    listen [::]:80;
    
    location / {
        proxy_pass http://monitor:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
print_success "Generated Nginx configuration"

# Start all services
print_header "Starting Services"
docker compose up -d
print_success "All services started"

# Get public IP
PUBLIC_IP=$(get_public_ip)
SERVICE_URL="http://${PUBLIC_IP}:9090"

# Generate QR code using ASCII art (requires qrencode)
print_header "QR Code for Mobile Access"
if ! command -v qrencode &> /dev/null; then
    if command -v apt-get &> /dev/null; then
        print_warning "Installing qrencode..."
        apt-get update && apt-get install -y qrencode
    elif command -v yum &> /dev/null; then
        print_warning "Installing qrencode..."
        yum install -y qrencode
    else
        print_warning "Please install qrencode manually to generate QR code"
    fi
fi

if command -v qrencode &> /dev/null; then
    echo -e "\nScan this QR code with your mobile device:"
    qrencode -t ANSI "${SERVICE_URL}"
fi

print_header "Setup Complete!"
echo -e "Your monitoring service has been set up in: ${GREEN}$(pwd)${NC}"
echo -e "Service URL: ${GREEN}${SERVICE_URL}${NC}"
echo -e "\nTo check the logs:"
echo "cd ${PROJECT_DIR} && docker compose logs -f"
echo -e "\nTo stop the service:"
echo "cd ${PROJECT_DIR} && docker compose down"

# Test the service
print_header "Testing Service"
sleep 5  # Wait for services to start
if curl -s "${SERVICE_URL}/health" | grep -q "ok"; then
    print_success "Service is running correctly"
else
    print_warning "Service might not be running properly. Check logs with: docker compose logs -f"
fi