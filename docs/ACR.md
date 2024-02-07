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