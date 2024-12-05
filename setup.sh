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

# Check if domain argument is provided
if [ -z "$1" ]; then
    print_error "Please provide a domain name as an argument"
    echo "Usage: ./setup.sh yourdomain.com"
    exit 1
fi

DOMAIN=$1

# Set fixed unusual ports to avoid conflicts
HTTP_PORT=7891
HTTPS_PORT=7892

# Create project directory
PROJECT_DIR="monitor-${DOMAIN}"
print_header "Creating Project Directory"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"
print_success "Created and moved to directory: $PROJECT_DIR"

# Create docker-compose.yml
print_header "Creating Docker Compose Configuration"
cat > docker-compose.yml << EOF
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
      - "${HTTP_PORT}:80"
      - "${HTTPS_PORT}:443"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d
      - ./certbot/conf:/etc/letsencrypt
      - ./certbot/www:/var/www/certbot
    command: "/bin/sh -c 'while :; do sleep 6h & wait \$\${!}; nginx -s reload; done & nginx -g \"daemon off;\"'"
    networks:
      - monitor-network
    depends_on:
      - monitor

  certbot:
    image: certbot/certbot
    container_name: monitor-certbot
    volumes:
      - ./certbot/conf:/etc/letsencrypt
      - ./certbot/www:/var/www/certbot
    entrypoint: "/bin/sh -c 'trap exit TERM; while :; do certbot renew; sleep 12h & wait \$\${!}; done;'"
    networks:
      - monitor-network

networks:
  monitor-network:
    driver: bridge
EOF
print_success "Created docker-compose.yml"

# Create required directories
print_header "Creating Directory Structure"
mkdir -p nginx/conf.d data certbot/conf certbot/www
print_success "Created required directories"

# Get SSL certificate using standalone mode
print_header "Obtaining SSL Certificate"
docker compose down nginx
docker run -it --rm --name certbot \
    -v "$PWD/certbot/conf:/etc/letsencrypt" \
    -v "$PWD/certbot/www:/var/www/certbot" \
    -p ${HTTP_PORT}:80 \
    certbot/certbot certonly --standalone \
    --preferred-challenges http \
    --email admin@${DOMAIN} \
    --agree-tos \
    --no-eff-email \
    -d ${DOMAIN}

if [ $? -ne 0 ]; then
    print_error "Failed to obtain SSL certificate"
    exit 1
fi
print_success "Obtained SSL certificate"

# Generate nginx config
print_header "Configuring Nginx"
cat > nginx/conf.d/app.conf << EOF
server {
    listen 80;
    listen [::]:80;
    server_name ${DOMAIN};
    
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://\$host:\$server_port\$request_uri;
    }
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name ${DOMAIN};

    ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_stapling on;
    ssl_stapling_verify on;
    add_header Strict-Transport-Security "max-age=31536000" always;

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
print_header "Starting All Services"
docker compose up -d
print_success "All services started"

print_header "Setup Complete!"
echo -e "Your monitoring service has been set up in: ${GREEN}$(pwd)${NC}"
echo -e "The service should now be available at: ${GREEN}https://${DOMAIN}:${HTTPS_PORT}${NC}"
echo -e "\n${YELLOW}Note: Using non-standard ports:${NC}"
echo "HTTP port: ${HTTP_PORT}"
echo "HTTPS port: ${HTTPS_PORT}"
echo "Make sure to include the port in your URL: https://${DOMAIN}:${HTTPS_PORT}"
echo -e "\nTo check the logs:"
echo "cd ${PROJECT_DIR} && docker compose logs -f"
echo -e "\nTo stop the service:"
echo "cd ${PROJECT_DIR} && docker compose down"