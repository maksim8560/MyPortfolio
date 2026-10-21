# Сборка портфолио: сначала фронтенд (он встраивается в бинарник), потом Go.
# Запуск: .\build.ps1

$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$version = (Get-Date -Format 'yyyyMMdd-HHmm')

Write-Host '→ npm install' -ForegroundColor Cyan
Push-Location (Join-Path $root 'web')
npm install --no-audit --no-fund
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'npm install провалился' }

Write-Host '→ vite build' -ForegroundColor Cyan
npm run build
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'сборка фронтенда провалилась' }
Pop-Location

Write-Host '→ go build' -ForegroundColor Cyan
Push-Location (Join-Path $root 'server')
go build -ldflags "-X main.version=$version" -o portfolio.exe .
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'сборка Go провалилась' }
Pop-Location

Write-Host "Готово: server\portfolio.exe (версия $version)" -ForegroundColor Green
Write-Host 'Запуск: cd server; .\portfolio.exe' -ForegroundColor DarkGray
