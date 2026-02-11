@echo off
REM Setup Mullvad Proxy for MMCS Notifier
REM This script sets the MULLVAD_PROXY environment variable

echo ========================================
echo MMCS Notifier - Mullvad Proxy Setup
echo ========================================
echo.

REM Option 1: Use Mullvad's SOCKS5 proxy service (recommended)
REM Get endpoints from: https://mullvad.net/en/help/socks5-proxy/
REM Format: socks5://[country]-wireguard.mullvad.net:1080
REM Example: socks5://us-wireguard.mullvad.net:1080

REM Option 2: Use local SOCKS5 proxy (if you have one running)
REM Format: socks5://127.0.0.1:1080

REM Option 3: Use HTTP proxy (if configured)
REM Format: http://127.0.0.1:8080

echo Choose proxy type:
echo 1. Mullvad SOCKS5 endpoint (recommended)
echo 2. Local SOCKS5 proxy (127.0.0.1:1080)
echo 3. Custom proxy URL
echo 4. Skip (no proxy)
echo.
set /p choice="Enter choice (1-4): "

if "%choice%"=="1" (
    echo.
    echo Available Mullvad SOCKS5 endpoints:
    echo - us-wireguard.mullvad.net:1080 (United States)
    echo - jp-wireguard.mullvad.net:1080 (Japan)
    echo - se-wireguard.mullvad.net:1080 (Sweden)
    echo - de-wireguard.mullvad.net:1080 (Germany)
    echo - uk-wireguard.mullvad.net:1080 (United Kingdom)
    echo - See more at: https://mullvad.net/en/help/socks5-proxy/
    echo.
    set /p country="Enter country code (e.g., us, jp, se): "
    set MULLVAD_PROXY=socks5://%country%-wireguard.mullvad.net:1080
    echo.
    echo Set MULLVAD_PROXY=%MULLVAD_PROXY%
) else if "%choice%"=="2" (
    set MULLVAD_PROXY=socks5://127.0.0.1:1080
    echo.
    echo Set MULLVAD_PROXY=%MULLVAD_PROXY%
) else if "%choice%"=="3" (
    set /p custom="Enter proxy URL (e.g., socks5://127.0.0.1:1080): "
    set MULLVAD_PROXY=%custom%
    echo.
    echo Set MULLVAD_PROXY=%MULLVAD_PROXY%
) else (
    echo.
    echo No proxy configured. Skipping...
    goto :end
)

echo.
echo ========================================
echo Proxy Configuration:
echo MULLVAD_PROXY=%MULLVAD_PROXY%
echo ========================================
echo.
echo This environment variable is set for this session only.
echo To make it permanent, add it to your system environment variables.
echo.
echo To run the notifier with this proxy, use:
echo   set MULLVAD_PROXY=%MULLVAD_PROXY%
echo   notifier.exe
echo.
echo Or run this script before starting the notifier.
echo.

:end
pause
