# Copy the contents of ./production/static to ./cmd/server/static
Copy-Item -Path .\production\static\* -Destination .\cmd\server\static -Recurse -Force

# Build the Docker image
docker build --file .\build\Dockerfile . --tag fluffycore.rage.oidc:latest