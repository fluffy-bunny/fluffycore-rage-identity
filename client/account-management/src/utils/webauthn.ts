import { apiInstance } from '../api';
import { getCSRF } from './cookies';

function bufferDecode(value: string): Uint8Array {
  value = value.replace(/-/g, '+').replace(/_/g, '/');

  return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
}

function bufferEncode(value: ArrayBuffer): string {
  return btoa(
    String.fromCharCode.apply(null, Array.from(new Uint8Array(value))),
  )
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
}

interface Response {
  publicKey: PublicKeyCredentialCreationOptions;
}

export async function registerUser(callback: () => void): Promise<void> {
  try {
    const response = await apiInstance.get('/webauthn/register/begin', {
      headers: { 'X-Csrf-Token': getCSRF() },
    });
    const beginResponse: Response = response.data;

    beginResponse.publicKey.challenge = bufferDecode(
      beginResponse.publicKey.challenge as unknown as string,
    );

    beginResponse.publicKey.user.id = bufferDecode(
      beginResponse.publicKey.user.id as unknown as string,
    );

    const credentialsResponse = (await navigator.credentials.create({
      publicKey: beginResponse.publicKey,
    })) as PublicKeyCredential;
    console.log(credentialsResponse);
    let request = {
      id: credentialsResponse.id,
      rawId: bufferEncode(credentialsResponse.rawId!),
      type: credentialsResponse.type,
      response: {
        attestationObject: bufferEncode(
          (credentialsResponse.response as AuthenticatorAttestationResponse)
            .attestationObject,
        ),
        clientDataJSON: bufferEncode(
          credentialsResponse.response.clientDataJSON,
        ),
      },
    };
    //console.log(request);
    if (credentialsResponse) {
      await apiInstance.post(
        '/webauthn/register/finish',
        request,
        {
          headers: {
            'X-Csrf-Token': getCSRF(),
          },
        },
      );

      callback();
    }
  } catch (error) {
    console.error('Error:', error);
    alert('Registration failed!');
  }
}
