import { instance } from '../api';
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
  publicKey: PublicKeyCredentialRequestOptions;
}

export async function loginUser() {
  try {
    const response = await instance.get<Response>('/webauthn/login/begin', {
      headers: { 'X-Csrf-Token': getCSRF() },
    });

    const beginResponse: Response = response.data;

    beginResponse.publicKey.challenge = bufferDecode(
      beginResponse.publicKey.challenge as unknown as string,
    );

    if (beginResponse.publicKey.allowCredentials) {
      beginResponse.publicKey.allowCredentials.forEach((allowCredential) => {
        allowCredential.id = bufferDecode(
          allowCredential.id as unknown as string,
        );
      });
    }

    const credential = (await navigator.credentials.get(
      beginResponse,
    )) as PublicKeyCredential;

    const { authenticatorData, clientDataJSON, signature, userHandle, rawId } =
      credential.response as AuthenticatorAssertionResponse & {
        rawId: ArrayBuffer;
      };

    const finishResponse = await instance.post<{ redirectUrl: string }>(
      '/webauthn/login/finish',
      {
        id: credential.id,
        rawId: bufferEncode(rawId),
        type: credential.type,
        response: {
          authenticatorData: bufferEncode(authenticatorData),
          clientDataJSON: bufferEncode(clientDataJSON),
          signature: bufferEncode(signature),
          userHandle: userHandle ? bufferEncode(userHandle) : undefined, // Handle optional userHandle
        },
      },
      {
        headers: {
          'X-Csrf-Token': getCSRF(),
        },
      },
    );

    window.location.href = finishResponse.data.redirectUrl;
  } catch (error) {
    throw error;
  }
}
