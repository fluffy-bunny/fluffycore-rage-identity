document.addEventListener("DOMContentLoaded", function () {
  if (!window.PublicKeyCredential) {
    alert("Client not capable. Handle error.");
  }
});

function bufferDecode(base64URL) {
  const base64 = base64URL.replace(/\-/g, "+").replace(/\_/g, "/");
  return Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
}

function bufferEncode(value) {
  return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");
}

async function registerUser(returnUrl) {
  try {
    const response = await fetch("/webauthn/register/begin");
    if (!response.ok) {
      throw new Error(`Error fetching begin data: ${response.status}`);
    }
    const beginResponse = await response.json();
    console.log("beginResponse:", beginResponse);

    beginResponse.publicKey.challenge = bufferDecode(
      beginResponse.publicKey.challenge
    );
    beginResponse.publicKey.user.id = bufferDecode(beginResponse.publicKey.user.id);

    const credentialsResponse = await navigator.credentials.create({
      publicKey: beginResponse.publicKey,
    });
    console.log("credentialsResponse:", credentialsResponse);

    const attestationObject = credentialsResponse.response.attestationObject;
    const clientDataJSON = credentialsResponse.response.clientDataJSON;
    const rawId = credentialsResponse.rawId;
    const csrf = getCSRF();

    const finishResponse = await fetch("/webauthn/register/finish", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        "X-Csrf-Token": csrf,
      },
      body: JSON.stringify({
        id: credentialsResponse.id,
        rawId: bufferEncode(rawId),
        type: credentialsResponse.type,
        response: {
          attestationObject: bufferEncode(attestationObject),
          clientDataJSON: bufferEncode(clientDataJSON),
        },
      }),
    });

    if (!finishResponse.ok) {
      throw new Error(`Error registering user: ${finishResponse.status}`);
    }

    const data = await finishResponse.json();
    console.log("finishResponse:", data);
    window.location.href = returnUrl;
  } catch (error) {
    console.error("Error:", error);
    alert("Registration failed!"); // More specific message based on error
  }
}


async function LoginUser() {
  try {
    const response = await fetch("/webauthn/login/begin");
    if (!response.ok) {
      throw new Error(`Error fetching begin data: ${response.status}`);
    }
    const beginResponse = await response.json();
    console.log("beginResponse: ", beginResponse);

    beginResponse.publicKey.challenge = bufferDecode(
      beginResponse.publicKey.challenge
    );
    beginResponse.publicKey.allowCredentials.forEach((allowCredential) => {
      allowCredential.id = bufferDecode(allowCredential.id);
    });

    const credential = await navigator.credentials.get({
      publicKey: beginResponse.publicKey,
    });
    console.log("getCredential: ", credential);

    const authData = credential.response.authenticatorData;
    const clientDataJSON = credential.response.clientDataJSON;
    const rawId = credential.rawId;
    const sig = credential.response.signature;
    const userHandle = credential.response.userHandle || null; // Handle optional userHandle
    const csrf = getCSRF();

    const finishResponse = await fetch("/webauthn/login/finish", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        "X-Csrf-Token": csrf,
      },
      body: JSON.stringify({
        id: credential.id,
        rawId: bufferEncode(rawId),
        type: credential.type,
        response: {
          authenticatorData: bufferEncode(authData),
          clientDataJSON: bufferEncode(clientDataJSON),
          signature: bufferEncode(sig),
          userHandle: userHandle ? bufferEncode(userHandle) : undefined, // Handle optional userHandle
        },
      }),
    });

    if (!finishResponse.ok) {
      throw new Error(`Error logging in user: ${finishResponse.status}`);
    }

    const data = await finishResponse.json();
    console.log("finishResponse:", data);
    window.location.href = data.redirectUri;
  } catch (error) {
    console.error("Error:", error);
    alert("Login failed!"); // More specific message based on error
  }
}

