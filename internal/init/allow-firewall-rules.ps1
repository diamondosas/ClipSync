# Get the full path of an app named 'app.exe' in the current folder
$CurrentDirApp = Join-Path -Path $pwd -ChildPath "clipsync.exe"

Write-Host "Allowing mDNS for: $CurrentDirApp"

# Inbound Rule
New-NetFirewallRule -DisplayName "mDNS Inbound (Current Dir)" `
  -Program $CurrentDirApp -Direction Inbound -Protocol UDP -LocalPort 5353 -Action Allow

# Outbound Rule
New-NetFirewallRule -DisplayName "mDNS Outbound (Current Dir)" `
  -Program $CurrentDirApp -Direction Outbound -Protocol UDP -LocalPort 5353 -Action Allow

Write-Host "`nPress any key to close this window..." -ForegroundColor Yellow
$Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") | Out-Null