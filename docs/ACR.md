# ACR

Let's dive into **OIDC ACR (Authentication Context Class Reference)** values and their significance:

1. **Purpose of ACR Values**:

   - ACR values serve as a way to communicate the **level of authentication** that occurred during the user's login process.
   - They provide the client (your application) with confidence about the quality and strength of the authentication.
   - These values are agreed upon between the client and the identity provider (IdP).

2. **What Are ACR Values?**:

   - The ACR value is a **URI** (Uniform Resource Identifier) that identifies the specific **authentication context class reference**.
   - Different services and standards may have varying sets of supported ACR values.
   - The IdP presents the supported ACR values as an array of strings during the OIDC IDP discovery process.

3. **Flexibility and Interpretation**:

   - There are no **officially standardized** ACR values. Instead, they are arbitrary labels agreed upon by the client and IdP.
   - The ACR values are not required parameters, which means there's flexibility in how they are implemented and interpreted.
   - Clients can use the `acr_values_supported` parameter from the OIDC discovery response to understand which ACR values are supported by the IdP ¹².

4. **Examples**:
   - An example ACR value might be `"urn:mace:incommon:iap:silver"`, indicating a certain level of authentication assurance ¹.
   - However, it's essential to note that ACR values can vary based on the context and the specific system.

In summary, ACR values help convey the authentication level achieved during the login process, allowing clients to make informed decisions based on the quality of authentication. 🛡️

Source: Conversation with Bing, 2/7/2024
(1) Can someone explain ACR return values in OIDC? - Stack Overflow. https://stackoverflow.com/questions/52632690/can-someone-explain-acr-return-values-in-oidc.
(2) Understanding ACR and AMR Claims in Authentication: Practical Use Cases. https://cloudentity.com/developers/blog/practical_use_of_acr_and_amr_claims/.
(3) Step-up authentication using ACR values | Okta Developer. https://developer.okta.com/docs/guides/step-up-authentication/main/.
(4) OpenID Connect Extended Authentication Profile (EAP) ACR Values 1.0 .... https://openid.net/specs/openid-connect-eap-acr-values-1_0.html.
(5) undefined. https://www.iana.org/assignments/loa-profiles/loa-profiles.xhtml.
(6) undefined. https://openid.net/specs/openid-connect-discovery-1_0.html.

## ACR Terms

These are request to you make to the server requesting a certain level of authentication.

```bash
http://localhost:9044/oauth2/default/v1/authorize?client_id={clientId}
&response_type=code
&scope=openid email profile
&redirect_uri=https://${yourapp}/authorization-code/callback
&state=296bc9a0-a2a2-4a57-be1a-d0e2fd9bb601
&acr_values=urn:{root_provider}:loa:1fa:any
```

| Term                                   | Descrition                                                                                                                                                                                                     |
| -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| urn:{root_provider}:loa:1fa:any        | any type of auth is fine, mfa not required, any idp is ok                                                                                                                                                      |
| urn:{root_provider}:loa:mfa:any        | mfa is required. Since external IDPs cannot be counted on to provide this information, the MFA must come from the root idp.                                                                                    |
| urn:{root_provider}:loa:1fa:psk        | root passkey auth is required. Can be combined with also requesting mfa                                                                                                                                        |
| urn:{root_provider}:loa:idp:{idp_slug} | an authentication to an external idp is required. i.e. urn:{root_provider}:loa:idp:usa.ca.gov.dmv, where `urn:{root_provider}:loa:idp:usa.ca.gov.dmv` is the external idp of California DMV employee user base |

## ARM Terms

| Term       | Descrition                                                                                      |
| ---------- | ----------------------------------------------------------------------------------------------- |
| loa        | Level of authentication                                                                         |
| pwd        | User used password                                                                              |
| psk        | User used a passkey                                                                             |
| 1fa        | User used only 1 factor of authentication                                                       |
| mfa:{type} | User used more than 1 factor or authentication. `mfa:sms`, `mfa:auth_app`, `mfa:psk`, `mfa:pwd` |
| idp        | User was authenticated through IDP. Level of authentication on the idp is unknown               |

Since there are no requirements around AMR. Depending on your model you could also encode the idp(s) as ARM values. i.e. `idp:{idp_slug}`. I would suggest doing both since `idp` as a first class claim seems to be everywhere.

## ACR/AMR examples

`fluffyroot` is used as the root idp.

```bash
http://localhost:9044/oauth2/default/v1/authorize?client_id={clientId}
&response_type=code
&scope=openid email profile
&redirect_uri=https://${yourapp}/authorization-code/callback
&state=296bc9a0-a2a2-4a57-be1a-d0e2fd9bb601
&acr_values=urn:fluffyroot:loa:mfa:any urn:fluffyroot:loa:idp:usa.ca.gov.dmv
```

Response id_token

```json
{
    "sub": "00u47ijy7sRLaeSdC0g7",
    "ver": 1,
    "iss": "https://{yourOktaDomain}/oauth2/default",
    "aud": "0oa48e74ox4t7mQJX0g7",
    "iat": 1661289624,
    "exp": 1661293224,
    "jti": "ID.dz6ibX-YnBNlt14huAtBULam_Z0_wPG0ig5SWCy8XQU",
    "amr": [
        "idp",
        "mfa:auth_app"
    ],
    "acr": [
        "urn:fluffyidp:loa:mfa:any",
        "urn:fluffyidp:loa:idp:usa.ca.gov.dmv"
    ],
    "idp": [
        "fluffyroot",
        "usa.ca.gov.dmv"
    ],
    "auth_time": 1661289603,
    "at_hash": "w6BLQV3642TKWvaVwTAJuw"
}
```

Here we see the 2 idp's have vouched for the user. Also we see that the root idp honored a mfa request.

Depending on your model you could also encode the idp(s) as ARM values. i.e. `idp:{idp_slug}`.
