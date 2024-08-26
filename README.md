# FluffyCore Identity

## A simple and opinionated OIDC Authentication Service

![alt text](image-2.png)
This is a Proof-Of-Life authentication server.

The use case for this proof is the github.com user experience.

1. A user is a stand-alone entity.
2. A user can be linked in N number of external IDPs.
3. A user can be challenged at any time against any known IDP , and the id_token must contain what idp (external or the root) wence the identity was produced.
4. External IDPs are secret. We don't want anyone to know what external enterprises a user can be linked to.

No calls to the userinfo endpoint are supported. id_token is the only thing returned that is useful. It is meant to use that id_token as an argument to an internal token_exchange that knows more about the user in the context of that system.

## TL;DR

Just configure, build and run this thing;

### Social and Enterprise IDP

These will not work out of the box. You need to have your own credentials.

The [idps.docker.json config file](./cmd/server/config/idps.docker.json) is where you can put your own social and enterprises idps. They will pull their client_secrets from the `.env.secrets` file.

#### .env.secrets

copy `.env.secrets.example` to `.env.secrets` and fill in the blanks.

```txt
# Secrets
#--------------------------------------------------
GITHUB_68863c06bc5c9bd0c2f9_CLIENT_SECRET=**REDACTED**
GOOGLE_1096301616546-edbl612881t7rkpljp3qa3juminskulo.apps.googleusercontent.com_CLIENT_SECRET=**REDACTED**
AZUREAD_3b918868-9bff-431f-bd9c-f9896d628e6b_CLIENT_SECRET=**REDACTED**
AZUREAD_0f81aa6c-b280-4503-b130-adc0567bfbe4_CLIENT_SECRET=**REDACTED**
```

If you do nothing then the only thing that will work will be username/password logins, and passkeys.

### Windows

#### Host file

```txt
127.0.0.1 localhost.dev traefik.localhost.dev whoami.localhost.dev smtp.localhost.dev rage.localhost.dev
```

```powershell
.\mkcert.exe -install
.\mkcert.exe -cert-file certs/local-cert.pem -key-file certs/local-key.pem "localhost.dev" "*.localhost.dev"

docker build --file .\build\Dockerfile . --tag fluffycore.rage.oidc:latest
docker-compose up -d
```

Now that we have the server running in docker, lets run our client locally.

```powershell
cd cmd/go-client
go build .

$env:PORT = "5556";$env:OAUTH2_CLIENT_ID = "go-client";$env:OAUTH2_CLIENT_SECRET = "secret";$env:AUTHORITY = "https://rage.localhost.dev"; .\go-client.exe
```

Open your browser, [Edge](https://www.microsoft.com/en-us/edge) is best and we all know it!

Navigate to [http://localhost:5556/login](http://localhost:5556/login)

Any username and password will work.

You should see a json response like this.

```json
{
  "OAuth2Token": {
    "access_token": "eyJhbGciOiJFUzI1NiIsImtpZCI6ImYzZTlmMjRjYTQ3MzRjNGU4YTQ4ZDI3ZjRhMmVmMjUyIiwidHlwIjoiSldUIn0.eyJhdWQiOiJnby1jbGllbnQiLCJjbGllbnRfaWQiOiJnby1jbGllbnQiLCJleHAiOjE3MTUxODIzMTAsImlhdCI6MTcxNTE3ODcxMCwiaXNzIjoiaHR0cHM6Ly9yYWdlLmxvY2FsaG9zdC5kZXYiLCJqdGkiOiJjb3RvcGxoM2NyaHBwa2RzdHE3ZyIsIm5iZiI6MTcxNTE3ODQxMCwic3ViIjoicmFnZV9jb3RvcGpoM2NyaHBwa2RzdHByZyJ9.ivUv29f2_bwtH-h1vM0Tb9VV18-cBBKJMfGAn4oCHxxW10UVwWo2UHzDU5BCUuIuvMav8bbNNy6aWQbDFfTyoQ",
    "token_type": "bearer",
    "expiry": "2024-05-08T08:31:50.2769354-07:00"
  },
  "IDTokenClaims": {
    "acr": ["urn:rage:idp:root", "urn:rage:password"],
    "amr": ["pwd", "idp", "mfa", "emailcode"],
    "aud": "go-client",
    "client_id": "go-client",
    "email": "ghstahl@gmail.com",
    "email_verified": false,
    "exp": 1715182310,
    "iat": 1715178710,
    "idp": ["root"],
    "iss": "https://rage.localhost.dev",
    "jti": "cotoplh3crhppkdstq70",
    "nbf": 1715178410,
    "nonce": "NoUCzwQrGM9WXqDwNHrpWw",
    "sub": "rage_cotopjh3crhppkdstprg"
  },
  "IDToken": "eyJhbGciOiJFUzI1NiIsImtpZCI6ImYzZTlmMjRjYTQ3MzRjNGU4YTQ4ZDI3ZjRhMmVmMjUyIiwidHlwIjoiSldUIn0.eyJhY3IiOlsidXJuOnJhZ2U6aWRwOnJvb3QiLCJ1cm46cmFnZTpwYXNzd29yZCJdLCJhbXIiOlsicHdkIiwiaWRwIiwibWZhIiwiZW1haWxjb2RlIl0sImF1ZCI6ImdvLWNsaWVudCIsImNsaWVudF9pZCI6ImdvLWNsaWVudCIsImVtYWlsIjoiZ2hzdGFobEBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImV4cCI6MTcxNTE4MjMxMCwiaWF0IjoxNzE1MTc4NzEwLCJpZHAiOlsicm9vdCJdLCJpc3MiOiJodHRwczovL3JhZ2UubG9jYWxob3N0LmRldiIsImp0aSI6ImNvdG9wbGgzY3JocHBrZHN0cTcwIiwibmJmIjoxNzE1MTc4NDEwLCJub25jZSI6Ik5vVUN6d1FyR005V1hxRHdOSHJwV3ciLCJzdWIiOiJyYWdlX2NvdG9wamgzY3JocHBrZHN0cHJnIn0.tWHugvPE8AN-QPicdx3Jdm1OfvpE77CtMz367tKr2_QeY9YC6Obx21AJDj0FT7qZLpjl-ylzf1MTniV2q-Wl5w"
}
```

**_Note the following claims in the id_token;_**

```json
{
  "acr": ["urn:rage:idp:root", "urn:rage:password"],
  "amr": ["pwd", "idp", "mfa", "emailcode"],
  "idp": ["root"],
  "sub": "rage_cotopjh3crhppkdstprg"
}
```

Context is important. The id_token normalizes the user to the sub claim. No matter how you login, passkey, password, social, enterprise, etc. The sub claim is always the same. The id_token will contain the acr and amr claims that tell you how the user was authenticated. The idp claim tells you where the user was authenticated. This is important because the user can be linked to multiple external IDPs. The id_token will tell you which one was used.

In cases like github, a user will get challenged to login with their linked enterprise account, even though they are already logged in using their github username/password. If it all goes well, the acr, amr, and idp will reflect the enterprise IDP.

If you fail the challenge, you don't get access to the github enterprise resources.

## Protos

Note: I had to run bash on windows so I could pass `./api/proto/**/*.proto`

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/fluffy-bunny/fluffycore/protoc-gen-go-fluffycore-di/cmd/protoc-gen-go-fluffycore-di@latest

```

```bash
go get github.com/fluffy-bunny/fluffycore

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/helloworld/helloworld.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/types/primitives.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/types/filter.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/types/pagination.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/types/phone_number.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/models/client.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/models/idp.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/client/client.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/idp/idp.proto

protoc --go_out=. --go_opt paths=source_relative  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/models/user.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/user/user.proto

protoc --go_out=. --go_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/oidc/flows/oidc_flow.proto

protoc --go_out=. --go_opt paths=source_relative  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/external/models/user.proto
protoc --go_out=. --go_opt paths=source_relative  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/external/models/metadata.proto
protoc --go_out=. --go_opt paths=source_relative  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/external/user/user.proto


protoc --go_out=. --go_opt paths=source_relative  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out=allow_merge=true,merge_file_name=proto:./proto --go-grpc_out . --go-grpc_opt paths=source_relative --go-fluffycore-di_out .  --go-fluffycore-di_opt paths=source_relative,grpc_gateway=true  ./proto/types/webauthn/webauthn.proto

```

## Private OAuth2 server

The kit comes with a self contained oauth2 server.

Your apis need tokens, and [here](./cmd/server/config/client.json) we can define exactly what claims a given client will mint.

The client_credenitials flow is the only thing supported.

[discovery](http://localhost:50053/.well-known/openid-configuration)  
[jwks](http://localhost:50053/.well-known/jwks.json)

client_credentials example:

```bash
curl --location 'http://localhost:50053/oauth/token' --header 'Content-Type: application/x-www-form-urlencoded' --header 'Authorization: Basic Y2xpZW50MTpzZWNyZXQ=' --data-urlencode 'grant_type=client_credentials'
```

```json
{
  "access_token": "eyJhbGciOiJFUzI1NiIsImtpZCI6IjBiMmNkMmU1NGM5MjRjZTg5ZjAxMGYyNDI4NjIzNjdkIiwidHlwIjoiSldUIn0.eyJjbGllbnRfaWQiOiJjbGllbnQxIiwiZXhwIjoxNjk5MjI3MzY3LCJpYXQiOjE2OTkyMjM3NjcsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwNTMiLCJwZXJtaXNzaW9ucyI6WyJyZWFkIiwid3JpdGUiXSwic3ViIjoiY2xpZW50MSJ9.hAtAa5W81NATUZmNDVQdQLYSmA_0Wx4HvmSMOcqGMdQMS7ay99v1RmKf-kT2l8Xm6rDMG8klIiEU9M-FK-400w",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

## Docker Build

```bash
 docker build --file .\build\Dockerfile . --tag fluffycore.rage.oidc:latest
```

## Health check

[go-healthcheck](https://github.com/phramz/go-healthcheck)

```yaml
COPY --from=gregthebunny/go-healthcheck /bin/healthcheck /bin/healthcheck
ENV PROBE='{{ .Assert.HTTPBodyContains .HTTP.Handler "GET" "http://localhost:50052/healthz" nil "SERVING" }}'
HEALTHCHECK --start-period=10s --retries=3 --timeout=10s --interval=10s \
CMD ["/bin/healthcheck", "probe", "$PROBE"]
```

Now all that is needed for another service to check health is a `condition: service_healthy`

```yaml
whoami:
  container_name: whoami
  extends:
    file: ./docker-compose-common.yml
    service: micro
  image: containous/whoami
  security_opt:
    - no-new-privileges:true
  depends_on:
    starterkit:
      condition: service_healthy
```

## Docker Compose

```bash
docker-compose -f .\docker-compose.yml up -d
```

## Swagger

[echo-swagger](https://github.com/swaggo/echo-swagger)

We want to swag init the general dir first, which is in the cmd/server directory. Then we want to include the swaggers in the internal

```bash
cd cmd/server
swag init  --dir ./,../../pkg,../../example/services/echo/account/api/api_user_profile
```

## GO OIDC CLIENT

```powershell
cd cmd/go-client
go build .

$env:PORT = "5556";$env:OAUTH2_CLIENT_ID = "go-client";$env:OAUTH2_CLIENT_SECRET = "secret";$env:AUTHORITY = "http://localhost:9044"; .\go-client.exe

$env:PORT = "5556";$env:OAUTH2_CLIENT_ID = "go-client";$env:OAUTH2_CLIENT_SECRET = "secret";$env:AUTHORITY = "http://localhost:9044"; .\go-client.exe
```

### Dev Client

```bash
cd cmd/oidc-client
go build .

.\oidc-client.exe serve    --authority http://localhost:9044 --client_id go-client --client_secret secret --port 5556

.\oidc-client.exe serve --acr_values "urn:rage:idp:google-social"   --authority http://localhost:9044 --client_id go-client --client_secret secret --port 5556

.\oidc-client.exe serve --acr_values "urn:rage:idp:mapped-enterprise" --acr_values "urn:rage:root_candidate:cnf07331og1ecp4r680g"  --authority http://localhost:9044 --client_id go-client --client_secret secret --port 5556

.\oidc-client.exe serve --acr_values "urn:rage:idp:mapped-enterprise"  --authority http://localhost:9044 --client_id go-client --client_secret secret --port 5556
```

.\oidc-client.exe serve --authority https://3156-47-150-126-75.ngrok-free.app --client_id go-client --client_secret secret --port 5556

### Docker Clients

```bash
.\oidc-client.exe serve    --authority https://rage.localhost.dev --client_id go-client --client_secret secret --port 5556

.\oidc-client.exe serve --acr_values "urn:rage:idp:mapped-enterprise"   --authority https://rage.localhost.dev --client_id go-client --client_secret secret --port 5556

.\oidc-client.exe serve --acr_values "urn:rage:idp:mapped-enterprise" --acr_values "urn:rage:root_candidate:cnf08ok1fnuu73eq91vg"   --authority https://rage.localhost.dev --client_id go-client --client_secret secret --port 5556


```

## PassKeys

For developement we need https. This is where ngrok comes in.

**NOTE**: Because we use ngrok we don't have a stable domain. So all IDP logins will fail, because we need to register a stable https domain with google, github, microsoft, azure, etc.

Passkey development can only work for simple username/password accounts.

```powershell
ngrok http http://localhost:9044
```

This will give you your ngrok url.

```cmd
ngrok
Forwarding  https://3156-47-150-126-75.ngrok-free.app -> http://localhost:9044
```

Update [.env.ngrok](./.env.ngrok) with the ngrok domain. In this case it would be `3156-47-150-126-75.ngrok-free.app`

We launch the server using vscode [launch.json](./.vscode/launch.json)

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "server-ngrok"
    }
  ]
}
```

## Recommended Configurations

- require email verification for password and social logins
- require email multifactor
- totp auth app is not really needed. The reason is that a users email is more important than your app. They should be doing way more multi factor over there. i.e. github, google or microsoft social. Enterprise IDPS have required multifactor and it looks like social accounts are now requiring it as well. Doing it here is just redundant and leads to account recovery problems.
- offer passkeys
