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
REM Mullvad SOCKS5 proxy endpoints use the format: [country]-wireguard.mullvad.net
REM However, if DNS resolution fails, try using the IP address directly or use a local proxy
REM 
REM Option 1: Use Mullvad's SOCKS5 endpoint (if DNS works):
REM set MULLVAD_PROXY=socks5://jp-wireguard.mullvad.net:1080
REM 
REM Option 2: Use local SOCKS5 proxy (recommended if you have Mullvad app running):
REM The Mullvad app provides a local SOCKS5 proxy on port 1080
set MULLVAD_PROXY=socks5://127.0.0.1:1080
REM 
REM Option 3: Disable proxy:
REM set MULLVAD_PROXY=

if defined MULLVAD_PROXY (
    echo [Mullvad Proxy] Enabled: %MULLVAD_PROXY%
    echo   Make sure Mullvad VPN app is running if using local proxy (127.0.0.1:1080)
) else (
    echo [Mullvad Proxy] Disabled (no proxy)
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
