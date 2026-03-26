param(
    [ValidateSet("htmx", "wasm")]
    [string]$UIMode = "htmx"
)

# Copy production WASM assets if they exist (only needed for wasm builds)
if (Test-Path .\production\static) {
    Copy-Item -Path .\production\static\* -Destination .\cmd\server\static -Recurse -Force
} elseif ($UIMode -eq "wasm") {
    Write-Warning "production/static not found - WASM build may be missing static assets"
}

$tag = "fluffycore.rage.oidc:latest"
if ($UIMode -eq "wasm") {
    $tag = "fluffycore.rage.oidc:latest-wasm"
}

# Build the Docker image
docker build --file .\build\Dockerfile --build-arg UI_MODE=$UIMode . --tag $tag

Write-Host "Built $tag with UI_MODE=$UIMode"