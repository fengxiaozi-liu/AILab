$ErrorActionPreference = "Stop"

$repo = "fengxiaozi-liu/AILab"
$asset = "ferryPilot-windows-amd64.exe"
$installDir = Join-Path $env:USERPROFILE "bin"
$target = Join-Path $installDir "ferryPilot.exe"
$url = "https://github.com/$repo/releases/latest/download/$asset"

New-Item -ItemType Directory -Force $installDir | Out-Null
Invoke-WebRequest -Uri $url -OutFile $target

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$paths = @()
if ($userPath) {
  $paths = $userPath -split ";"
}

if ($paths -notcontains $installDir) {
  $newPath = if ($userPath) { "$userPath;$installDir" } else { $installDir }
  [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
}

Write-Host "ferryPilot installed to $target"
Write-Host "Restart PowerShell, then run: ferryPilot -p speckit"
