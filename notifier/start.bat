@echo off
REM MMCS Discord Notifier Startup Script
REM This script starts the Discord notifier service

echo ========================================
echo MMCS Discord Notifier
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Change to the notifier directory
cd /d "%~dp0"

REM Check if .env file exists for environment variables
if exist .env (
    echo Loading environment variables from .env file...
    for /f "usebackq tokens=*" %%a in (".env") do set %%a
)

REM Set Mullvad Proxy (SOCKS5)
REM The Mullvad app provides a local SOCKS5 proxy on port 1080 when VPN is connected
REM 
REM IMPORTANT: To use Mullvad proxy:
REM 1. Connect to Mullvad VPN first
REM 2. In Mullvad app: Settings → Local SOCKS5 proxy → Enable
REM 3. The proxy will be available at 127.0.0.1:1080
REM 
REM Option 1: Use local SOCKS5 proxy (recommended - works when Mullvad VPN is connected):
REM Uncomment the line below to enable:
set MULLVAD_PROXY=socks5://127.0.0.1:1080
REM 
REM Option 2: Disable proxy (works without VPN):
REM Comment out the line above or set:
REM set MULLVAD_PROXY=

if defined MULLVAD_PROXY (
    echo [Mullvad Proxy] Enabled: %MULLVAD_PROXY%
    echo   Note: Proxy will only work if:
    echo   - Mullvad VPN is connected
    echo   - Local SOCKS5 proxy is enabled in Mullvad app settings
    echo   - If proxy unavailable, will fallback to direct connection
) else (
    echo [Mullvad Proxy] Disabled (direct connection, no VPN)
)
echo.

REM Set default Supabase Anon Key (if not set via environment)
if "%SUPABASE_SERVICE_ROLE_KEY%"=="" (
    if "%SUPABASE_ANON_KEY%"=="" (
        REM Set default anon key
        set SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzAxNzUyMjMsImV4cCI6MjA4NTc1MTIyM30.t48b38QU8QpWfDyGu__hTKdCYbjVh1rhHcrt1D7mFWU
        echo.
        echo Using default Supabase Anon Key (configured in start.bat)
        echo NOTE: For production, use SERVICE_ROLE_KEY instead (set via .env file)
        echo.
    )
)

REM Download dependencies if needed
echo Checking dependencies...
go mod download

REM Build the notifier
echo.
echo Building notifier...
go build -o notifier.exe

if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to build notifier
    pause
    exit /b 1
)

REM Run the notifier
echo.
echo Starting notifier...
echo Press Ctrl+C to stop
echo.
notifier.exe

REM If the notifier exits, pause so user can see any error messages
pause
