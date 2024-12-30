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

## Quick Setup

1. **One-Line Installation**
   ```bash
   curl -fsSL https://raw.githubusercontent.com/kaanmertkoc/simple-monitor/main/setup.sh -o setup.sh && chmod +x setup.sh && ./setup.sh
   ```

   Or manually:

   ```bash
   # Download setup script
   wget https://raw.githubusercontent.com/kaanmertkoc/simple-monitor/main/setup.sh
   
   # Make it executable
   chmod +x setup.sh
   
   # Run setup with your domain
   ./setup.sh
   ```

That's it! Your monitoring service will be available at `https://your-domain.com`

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

## Manual Setup

If you prefer to set up manually:

1. Create directory structure:
   ```bash
   mkdir -p monitor && cd monitor
   mkdir -p nginx/conf.d data certbot/conf certbot/www
   ```

2. Download setup files:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/kaanmertkoc/simple-monitor/main/docker-compose.yml -o docker-compose.yml
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

### SSL Certificate Issues
- Ensure DNS is properly configured
- Check if ports 80 and 443 are open
- Wait for DNS propagation (can take up to 48 hours)

To manually trigger certificate renewal:
```bash
docker compose run --rm certbot certonly --webroot --webroot-path /var/www/certbot -d
```

### Connection Refused
1. Check if containers are running:
   ```bash
   docker compose ps
   ```
2. Check logs:
   ```bash
   docker compose logs -f
   ```
   
   For specific services:
   ```bash
   docker compose logs -f monitor
   docker compose logs -f nginx
   ```

### Invalid Certificate
- Make sure your domain is pointing to the correct IP
- Check if SSL certificate was generated:
  ```bash
  ls -la certbot/conf/live/your-domain.com/
  ```

### Common Issues
- Ensure DNS is properly configured
- Check if ports 80 and 443 are open
- Verify domain points to correct IP
- Wait for DNS propagation (up to 48 hours)

## Security

This setup includes:
- Automatic HTTP to HTTPS redirection
- Modern SSL configuration
- Regular certificate renewal
- Secure headers

## Support

If you encounter any issues:
1. Check the [common issues](#troubleshooting) section
2. Open an issue on GitHub
3. Join our community Discord

## License
This project is licensed under the Creative Commons Attribution-NonCommercial 4.0 International License - see the [LICENSE](LICENSE) file for details.

[![License: CC BY-NC 4.0](https://img.shields.io/badge/License-CC%20BY--NC%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc/4.0/)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.