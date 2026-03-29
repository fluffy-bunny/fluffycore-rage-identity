param(
    [string]$Version = "",
    [string]$Commit = "",
    [string]$Date = ""
)

# Auto-generate version if not provided
if (-not $Version) {
    $Version = "local-" + (Get-Date -Format "yyyyMMddHHmmss")
}
if (-not $Commit) {
    $Commit = (git rev-parse --short HEAD 2>$null) ?? "unknown"
}
if (-not $Date) {
    $Date = Get-Date -Format "yyyyMMdd"
}

$tag = "fluffycore.rage.oidc:latest"

# Build the Docker image
docker build --file .\build\Dockerfile `
    --build-arg version=$Version `
    --build-arg commit=$Commit `
    --build-arg date=$Date `
    . --tag $tag

Write-Host "Built $tag (version=$Version, commit=$Commit, date=$Date)"