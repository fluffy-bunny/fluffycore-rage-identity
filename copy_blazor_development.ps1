# Copy the contents of ./production/static to ./cmd/server/static
Copy-Item -Path ..\fluffycore-rage-identity-blazor\cmd\httpserver\static\management\* -Destination .\cmd\server\static\blazor\management -Recurse -Force
Copy-Item -Path ..\fluffycore-rage-identity-blazor\cmd\httpserver\static\oidc-login-ui\* -Destination .\cmd\server\static\blazor\oidc-login-ui -Recurse -Force
Copy-Item -Path .\production\static\* -Destination .\cmd\server\static -Recurse -Force

 