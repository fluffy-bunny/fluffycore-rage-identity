.\mkcert.exe -install
.\mkcert.exe -cert-file certs/local-cert.pem -key-file certs/local-key.pem "localhost.dev" "*.localhost.dev" "rage.localhost.dev"  "*.rage.localhost.dev"
docker build --file .\build\Dockerfile . --tag fluffycore.rage.oidc:latest