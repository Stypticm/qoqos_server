# Force UTF-8 for this session to avoid mojibake
$OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "--- QOQOS SYSTEM STARTUP ---" -ForegroundColor Green

Write-Host "[1/4] Starting Infrastructure (PostgreSQL, Qdrant)..." -ForegroundColor Cyan
docker-compose up -d

Write-Host "[2/4] Waiting for DB readiness..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

Write-Host "[3/4] Starting OpenClaw Gateway..." -ForegroundColor Cyan
# Start OpenClaw Gateway using npx
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd d:\github\projects\qoqos_server\openclaw; npx openclaw gateway --port 18789"

Write-Host "[4/4] Starting Telegram Bot (Valera)..." -ForegroundColor Cyan
# Start Telegram Bot (Valera)
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd d:\github\projects\qoqos_server; py valera_brain.py"

Write-Host "SUCCESS: All systems are running!" -ForegroundColor Green
Write-Host "Visit: http://localhost:5173" -ForegroundColor White
