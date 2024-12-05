# Server Resource Monitor

A lightweight server monitoring solution that tracks CPU and RAM usage through a secure REST API.

## Prerequisites

Before setting up the monitoring service, ensure you have:

1. **Domain Name**
   - A domain or subdomain that you can configure
   - Access to your domain's DNS settings

2. **Server Requirements**
   - A VPS or server with:
     - Docker installed
     - Docker Compose installed
     - Ports 80 and 443 available
     - At least 1GB RAM recommended

3. **DNS Configuration**
   If using a subdomain (e.g., `monitor.yourdomain.com`):
   ```
   Type: A
   Host: monitor
   Value: Your-Server-IP
   TTL: 3600 (or default)
   ```
   
   If using a root domain (e.g., `yourdomain.com`):
   ```
   Type: A
   Host: @
   Value: Your-Server-IP
   TTL: 3600 (or default)
   ```

## Quick Start

1. **Clone the Repository**
   ```bash
   git clone https://github.com/kaanmertkoc/simple-monitor.git
   cd simple-monitor
   ```

2. **Make the Setup Script Executable**
   ```bash
   chmod +x setup.sh
   ```

3. **Run the Setup Script**
   ```bash
   ./setup.sh your-domain.com
   ```

The setup script will:
- Check all prerequisites
- Set up Nginx as a reverse proxy
- Obtain an SSL certificate automatically
- Configure HTTPS
- Start the monitoring service

## Common DNS Providers Instructions

### Cloudflare
1. Log in to Cloudflare Dashboard
2. Select your domain
3. Go to DNS settings
4. Click "Add Record"
5. Choose A record
6. Enter subdomain (if using)
7. Enter your server's IP
8. Save

### GoDaddy
1. Log in to GoDaddy
2. Go to DNS Management
3. Find the DNS records section
4. Add new record
5. Choose A record
6. Enter subdomain (if using)
7. Enter your server's IP
8. Save

### NameCheap
1. Log in to Namecheap
2. Go to Domain List
3. Click "Manage" next to your domain
4. Go to Advanced DNS
5. Add new record
6. Choose A record
7. Enter subdomain (if using)
8. Enter your server's IP
9. Save

## Troubleshooting

### SSL Certificate Issues
- Ensure DNS is properly configured
- Check if ports 80 and 443 are open
- Wait for DNS propagation (can take up to 48 hours)

### Connection Refused
1. Check if containers are running:
   ```bash
   docker compose ps
   ```
2. Check logs:
   ```bash
   docker compose logs
   ```

### Invalid Certificate
- Make sure your domain is pointing to the correct IP
- Check if SSL certificate was generated:
  ```bash
  ls -la certbot/conf/live/your-domain.com/
  ```

## Security Notes

This setup includes:
- Automatic HTTP to HTTPS redirection
- Modern SSL configuration
- Regular certificate renewal
- Secure headers

## Updating

To update the monitoring service:

```bash
git pull
docker compose down
docker compose up -d --build
```

## License

[Your License]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


# Simple Monitor

A lightweight server monitoring solution that tracks CPU and RAM usage through a secure REST API.

## Quick Setup

1. **Prerequisites**
   - A server with Docker and Docker Compose installed
   - A domain name pointing to your server
   - Ports 80 and 443 available

2. **One-Line Installation**
   ```bash
   curl -fsSL https://raw.githubusercontent.com/yourusername/simple-monitor/main/setup.sh -o setup.sh && chmod +x setup.sh && ./setup.sh your-domain.com
   ```

   Or manually:

   ```bash
   # Download setup script
   wget https://raw.githubusercontent.com/yourusername/simple-monitor/main/setup.sh
   
   # Make it executable
   chmod +x setup.sh
   
   # Run setup with your domain
   ./setup.sh your-domain.com
   ```

That's it! Your monitoring service will be available at `https://your-domain.com`

## DNS Configuration

Before running the setup, make sure your domain points to your server:

1. Get your server's IP address
2. Add an A record in your DNS settings:
   ```
   Type: A
   Host: @ (or subdomain)
   Value: Your-Server-IP
   TTL: 3600
   ```

## Manual Setup

If you prefer to set up manually:

1. Create directory structure:
   ```bash
   mkdir -p monitor && cd monitor
   mkdir -p nginx/conf.d data certbot/conf certbot/www
   ```

2. Download setup files:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/yourusername/simple-monitor/main/docker-compose.yml -o docker-compose.yml
   ```

3. Start the service:
   ```bash
   docker compose up -d
   ```

## Usage

Access your monitoring dashboard at `https://your-domain.com`

API Endpoints:
- `/metrics` - Current metrics
- `/metrics/history` - Historical data
- `/health` - Health check

## Management Commands

```bash
# View logs
docker compose logs -f

# Stop service
docker compose down

# Update to latest version
docker compose pull
docker compose up -d

# Restart service
docker compose restart
```

## Troubleshooting

1. **Certificate Issues**
   ```bash
   # Manually trigger certificate renewal
   docker compose run --rm certbot certonly --webroot --webroot-path /var/www/certbot -d your-domain.com
   ```

2. **Check Logs**
   ```bash
   # All services
   docker compose logs -f

   # Specific service
   docker compose logs -f monitor
   docker compose logs -f nginx
   ```

3. **Common Issues**
   - Ensure DNS is properly configured
   - Check if ports 80 and 443 are open
   - Verify domain points to correct IP
   - Wait for DNS propagation (up to 48 hours)

## Security

This setup includes:
- Automatic HTTPS redirection
- Modern SSL configuration
- Regular certificate renewal
- Secure headers

## Support

If you encounter any issues:
1. Check the [common issues](#troubleshooting) section
2. Open an issue on GitHub
3. Join our community Discord

## License

[Your License]