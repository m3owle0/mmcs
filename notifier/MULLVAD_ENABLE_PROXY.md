# How to Enable Mullvad Local SOCKS5 Proxy

The notifier uses Mullvad's local SOCKS5 proxy (`127.0.0.1:1080`) when the VPN is connected. Here's how to enable it:

## Step-by-Step Instructions

### 1. Connect to Mullvad VPN
- Open the Mullvad VPN app
- Click "Connect" to establish a VPN connection
- Wait until you see "Connected" status

### 2. Enable Local SOCKS5 Proxy
- In the Mullvad app, go to **Settings** (gear icon)
- Scroll down to find **"Local SOCKS5 proxy"** section
- Toggle **"Enable local SOCKS5 proxy"** to **ON**
- The proxy will be available at `127.0.0.1:1080`

### 3. Verify Proxy is Working
- The notifier will automatically detect if the proxy is available
- If proxy is not available, it will fallback to direct connection
- Check the logs for: `🌐 Using SOCKS5 proxy: socks5://127.0.0.1:1080 (verified)`

## Troubleshooting

### Error: "No connection could be made because the target machine actively refused it"

This means the SOCKS5 proxy is not running. Check:

1. **Is Mullvad VPN connected?**
   - The local proxy only works when VPN is connected
   - Check Mullvad app status - should show "Connected"

2. **Is Local SOCKS5 proxy enabled?**
   - Mullvad app → Settings → Local SOCKS5 proxy → Enable
   - Make sure it's toggled ON

3. **Is port 1080 available?**
   - Another application might be using port 1080
   - Check with: `netstat -an | findstr 1080`

4. **Firewall blocking?**
   - Windows Firewall might be blocking localhost connections
   - Try temporarily disabling firewall to test

### The notifier works without proxy

That's fine! The notifier will automatically fallback to direct connection if the proxy is unavailable. You'll see:
```
⚠️  SOCKS5 proxy at 127.0.0.1:1080 is not available
   Falling back to direct connection...
🌐 No proxy configured (direct connection)
```

## Benefits of Using Mullvad Proxy

- **Privacy**: All traffic routes through Mullvad VPN
- **IP Rotation**: Can use different Mullvad servers
- **Reduced Rate Limiting**: Different IP addresses for requests
- **Geo-location**: Access region-specific content

## Disabling Proxy

To disable the proxy, comment out or remove the line in `start.bat`:
```batch
REM set MULLVAD_PROXY=socks5://127.0.0.1:1080
```

Or set it to empty:
```batch
set MULLVAD_PROXY=
```
