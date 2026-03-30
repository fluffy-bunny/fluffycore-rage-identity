.\mkcert.exe -install
.\mkcert.exe -cert-file certs/local-cert.pem -key-file certs/local-key.pem "localhost.dev" "*.localhost.dev" "rage.localhost.dev"  "*.rage.localhost.dev"

$Version = "local-" + (Get-Date -Format "yyyyMMddHHmmss")
$Commit = (git rev-parse --short HEAD 2>$null) ?? "unknown"
$Date = Get-Date -Format "yyyyMMdd"

docker build --file .\build\Dockerfile `
    --build-arg version=$Version `
    --build-arg commit=$Commit `
    --build-arg date=$Date `
    . --tag fluffycore.rage.oidc:latest

Write-Host "Built fluffycore.rage.oidc:latest (version=$Version, commit=$Commit, date=$Date)"