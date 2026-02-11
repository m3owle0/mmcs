# Mullvad VPN Proxy Setup Guide

The MMCS Notifier supports Mullvad VPN through SOCKS5 proxy configuration.

## Important Note

The WireGuard connection details shown in the Mullvad app (like `23.162.8.67:45700`) are **VPN endpoints**, not proxy endpoints. For HTTP requests, you need to use **Mullvad's SOCKS5 proxy service** instead.

## Option 1: Use Mullvad's SOCKS5 Proxy Service (Recommended)

Mullvad provides SOCKS5 proxy endpoints that you can use directly:

### Available Endpoints

- **United States**: `socks5://us-wireguard.mullvad.net:1080`
- **Japan**: `socks5://jp-wireguard.mullvad.net:1080`
- **Sweden**: `socks5://se-wireguard.mullvad.net:1080`
- **Germany**: `socks5://de-wireguard.mullvad.net:1080`
- **United Kingdom**: `socks5://uk-wireguard.mullvad.net:1080`

See all available endpoints: https://mullvad.net/en/help/socks5-proxy/

### Setup

**Windows (Command Prompt):**
```cmd
set MULLVAD_PROXY=socks5://us-wireguard.mullvad.net:1080
notifier.exe
```

**Windows (PowerShell):**
```powershell
$env:MULLVAD_PROXY="socks5://us-wireguard.mullvad.net:1080"
.\notifier.exe
```

**Linux/Mac:**
```bash
export MULLVAD_PROXY=socks5://us-wireguard.mullvad.net:1080
./notifier
```

## Option 2: Use Local SOCKS5 Proxy

If you have a local SOCKS5 proxy running (e.g., through Mullvad's local proxy feature):

```cmd
set MULLVAD_PROXY=socks5://127.0.0.1:1080
```

## Option 3: Use .env File

Create a `.env` file in the `notifier` directory:

```env
MULLVAD_PROXY=socks5://us-wireguard.mullvad.net:1080
```

Note: The notifier currently reads from environment variables, so you'll need to load the `.env` file or use a tool like `dotenv`.

## Option 4: Use Windows Setup Script

Run the provided batch file:

```cmd
setup_mullvad_proxy.bat
```

This will guide you through setting up the proxy interactively.

## Verification

When you start the notifier, you should see:

```
🌐 Using SOCKS5 proxy: socks5://us-wireguard...
```

If you see:

```
🌐 No proxy configured...
```

Then the proxy is not set correctly.

## Troubleshooting

1. **Connection fails**: Make sure the SOCKS5 endpoint is accessible and your Mullvad account is active
2. **Wrong endpoint**: Try a different country endpoint
3. **Port issues**: Make sure port 1080 is not blocked by your firewall
4. **Authentication**: Mullvad SOCKS5 proxies don't require authentication, so don't include username/password in the URL

## Making Proxy Permanent

To make the proxy setting permanent on Windows:

1. Open System Properties → Environment Variables
2. Add new System Variable:
   - Name: `MULLVAD_PROXY`
   - Value: `socks5://us-wireguard.mullvad.net:1080`
3. Restart your terminal/application

## Using Your WireGuard Connection

If you want to use your specific WireGuard connection (`23.162.8.67:45700`), you'll need to:

1. Set up a local SOCKS5 proxy that routes through your WireGuard connection
2. Use tools like `ssh -D` or a SOCKS5 proxy server
3. Then point `MULLVAD_PROXY` to your local proxy (e.g., `socks5://127.0.0.1:1080`)
