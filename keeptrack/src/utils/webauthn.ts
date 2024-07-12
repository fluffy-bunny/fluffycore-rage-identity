// import { getCSRF } from './cookies';

// function bufferDecode(value: string): Uint8Array {
//   value = value.replace(/-/g, '+').replace(/_/g, '/');
//   return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
// }

// function bufferEncode(value: ArrayBuffer): string {
//   return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
//     .replace(/\+/g, '-')
//     .replace(/\//g, '_')
//     .replace(/=/g, '');
// }

async function registerUser(): Promise<void> {
  // try {
  //   const response = await fetch('/webauthn/register/begin');
  //   if (!response.ok) {
  //     throw new Error(`Error fetching begin data: ${response.status}`);
  //   }
  //   const beginResponse: Response = await response.json();
  //   beginResponse.publicKey.challenge = bufferDecode(
  //     beginResponse.publicKey.challenge as string,
  //   );
  //   beginResponse.publicKey.user.id = bufferDecode(
  //     beginResponse.publicKey.user.id as string,
  //   );
  //   const credentialsResponse = await navigator.credentials.create({
  //     publicKey: beginResponse.publicKey,
  //   });
  //   const csrf = getCSRF();
  //   const finishResponse = await fetch('/webauthn/register/finish', {
  //     method: 'POST',
  //     headers: {
  //       Accept: 'application/json',
  //       'Content-Type': 'application/json',
  //       'X-Csrf-Token': csrf,
  //     },
  //     body: JSON.stringify({
  //       id: credentialsResponse.id,
  //       rawId: bufferEncode(credentialsResponse.rawId),
  //       type: credentialsResponse.type,
  //       response: {
  //         attestationObject: bufferEncode(
  //           credentialsResponse.response.attestationObject,
  //         ),
  //         clientDataJSON: bufferEncode(
  //           credentialsResponse.response.clientDataJSON,
  //         ),
  //       },
  //     }),
  //   });
  //   if (!finishResponse.ok) {
  //     throw new Error(`Error registering user: ${finishResponse.status}`);
  //   }
  //   const data = await finishResponse.json();
  //   window.location.href = returnUrl;
  // } catch (error) {
  //   console.error('Error:', error);
  //   alert('Registration failed!');
  // }
}

export const webauthn = {
  registerUser,
};
