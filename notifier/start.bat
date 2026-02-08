@echo off
echo ========================================
echo   Discord Notifier Service
echo ========================================
echo.

REM Set your Supabase Anon Key here
set SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzAxNzUyMjMsImV4cCI6MjA4NTc1MTIyM30.t48b38QU8QpWfDyGu__hTKdCYbjVh1rhHcrt1D7mFWU

REM Check if key is set to placeholder
if "%SUPABASE_ANON_KEY%"=="YOUR_KEY_HERE" (
    echo ERROR: Please edit start.bat and set your SUPABASE_ANON_KEY
    echo.
    echo Get your key from:
    echo 1. Go to: https://supabase.com/dashboard
    echo 2. Select your project
    echo 3. Click Project Settings gear icon or Connect button
    echo 4. Go to API Keys section
    echo 5. Copy the Publishable key OR Legacy anon key
    echo.
    pause
    exit /b 1
)

REM Check if dependencies are installed
if not exist "go.sum" (
    echo Installing dependencies...
    echo.
    set GOPROXY=https://proxy.golang.org,direct
    go mod tidy
    if errorlevel 1 (
        echo.
        echo ERROR: Failed to download dependencies.
        echo Try running: set GOPROXY=https://proxy.golang.org,direct && go mod tidy
        pause
        exit /b 1
    )
    echo.
)

REM Run the notifier
echo Starting notifier...
echo.
go run .

pause
