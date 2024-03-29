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

Just build and run this thing;

```bash
 docker build --file .\build\Dockerfile . --tag fluffycore.rage.oidc:latest
 docker-compose up -d
```

Now that we have the server running in docker, lets run our client locally.

```powershell
cd cmd/go-client
go build .

$env:PORT = "5556";$env:OAUTH2_CLIENT_ID = "go-client";$env:OAUTH2_CLIENT_SECRET = "secret";$env:AUTHORITY = "http://localhost:9044/"; .\go-client.exe
```

Open your browser, [Edge](https://www.microsoft.com/en-us/edge) is best and we all know it!

Navigate to [http://localhost:5556/login](http://localhost:5556/login)

Any username and password will work.

You should see a json response like this.

```json
{
  "OAuth2Token": {
    "access_token": "eyJhbGciOiJFUzI1NiIsImtpZCI6IjBiMmNkMmU1NGM5MjRjZTg5ZjAxMGYyNDI4NjIzNjdkIiwidHlwIjoiSldUIn0.eyJhdWQiOiJnby1jbGllbnQiLCJjbGllbnRfaWQiOiJnby1jbGllbnQiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJleHAiOjE3MDc1MjQ0ODUsImlhdCI6MTcwNzUyMDg4NSwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo5MDQ0IiwianRpIjoiY24zYjZ0YWkycW5iMzc4MjgwbjAiLCJuYmYiOjE3MDc1MjA1ODUsInBlcm1pc3Npb25zIjpbInJlYWQiLCJ3cml0ZSJdLCJzdWIiOiIxMjMifQ.R9zQX2njveB-iUhQTO698logMjPtFdDbe7Ne2scSoT8kcMEMk3wEIz2D8tyzcjSlsqSSoXoAP6YKo1dIfnFOOQ",
    "token_type": "bearer",
    "refresh_token": "refresh_token",
    "expiry": "2024-02-09T16:21:26.0012199-08:00"
  },
  "IDTokenClaims": {
    "aud": "go-client",
    "client_id": "go-client",
    "email": "test@test.com",
    "exp": 1707524485,
    "iat": 1707520885,
    "iss": "http://localhost:9044",
    "jti": "cn3b6tai2qnb378280mg",
    "nbf": 1707520585,
    "nonce": "AeJpC-NrPt0ED3-Qh2M34g",
    "sub": "123"
  },
  "IDToken": "eyJhbGciOiJFUzI1NiIsImtpZCI6IjBiMmNkMmU1NGM5MjRjZTg5ZjAxMGYyNDI4NjIzNjdkIiwidHlwIjoiSldUIn0.eyJhdWQiOiJnby1jbGllbnQiLCJjbGllbnRfaWQiOiJnby1jbGllbnQiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJleHAiOjE3MDc1MjQ0ODUsImlhdCI6MTcwNzUyMDg4NSwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo5MDQ0IiwianRpIjoiY24zYjZ0YWkycW5iMzc4MjgwbWciLCJuYmYiOjE3MDc1MjA1ODUsIm5vbmNlIjoiQWVKcEMtTnJQdDBFRDMtUWgyTTM0ZyIsInN1YiI6IjEyMyJ9.0KuxDAlXX4DIh5Lh0KSXTahY8gQicRYVWMd-4Ic8J5ZwbFwnrFPk_sgE2cGetcaAFiReHg1SYAszsY8Sahds6A"
}
```

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
swag init  --dir ./,../../internal
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
