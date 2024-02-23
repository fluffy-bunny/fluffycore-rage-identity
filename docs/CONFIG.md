# Config

## Mounted Files

### Signing Keys

[generated by](https://github.com/fluffy-bunny/crypto-gen)  

There is no key rotation engine.  Instead we have mounted 10 years worth of keys that overlap.  If compromised, generate a new set and restart the service.  

[example](../cmd/server/config/signing-keys.json)  

### IDP Configuration

[example while developing](../cmd/server/config/idps.json)  
[example for docker](../cmd/server/config/idps.docker.json)  

You can also do ENV replacements in the idps config.  

```json
{
    "a_secret": "${MY_SECRET}"
}
```

In the case of the IDPs.  

```json
"client_secret": "${GOOGLE_1096301616546-edbl612881t7rkpljp3qa3juminskulo.apps.googleusercontent.com_CLIENT_SECRET}",
"client_secret": "${GITHUB_68863c06bc5c9bd0c2f9_CLIENT_SECRET}"
"client_secret": "${AZUREAD_0f81aa6c-b280-4503-b130-adc0567bfbe4_CLIENT_SECRET}",
"client_secret": "${AZUREAD_3b918868-9bff-431f-bd9c-f9896d628e6b_CLIENT_SECRET}",
```

| Key           | Description   | .etc  |
| ------------- |:-------------:| -----:|
| slug          | a stable idp identifier | required, lowercase |
| name      | If not hidden, the name of the button      |    |
| enabled      | flag to make the idp available      | true,false   |
| hidden | hidden from the login page.      | true,false |
| emailVerificationRequired | email code verification   | true,false |
| autoCreate | if an account is created on first login | true,false |
| claimedDomains | list of domains that are claim. No email/password account can be created for a claimed domain.  Typing in a claimed domain will auto redirect to the claiming idp. | array of strings, "claimedDomains": ["mapped.com"] |
| metadata.image_ref | abstract reference to an image for presentation  | extensibility |  

### OIDC Clients

These are your webapps that will use tis IDP to authenticate.  

[example for docker](../cmd/server/config/oidcClients.json.json)  

| Key           | Description   | .etc  |
| ------------- |:-------------:| -----:|
| client_id          | the client id | required |
| allowed_grant_types      | "authorization_code"  | string   |
| client_secrets      | array or $argon2id hashs. [generator](https://argon2.online/)    | required   |
| allowed_redirect_uris | list of redirect urls     | required |

#### client_secret example

```json
{
    "client_secrets": [
        {
          "id": "cim6nke5cpqs73cqthgg",
          "name": "secret",
          "expiration_unix": "64843141941",
          "hash": "$argon2id$v=19$m=16,t=2,p=1$MTIzNDU2Nzg$ZlgCW0V2GfYHJSpwaUWU1w"
        }
      ]
}
```