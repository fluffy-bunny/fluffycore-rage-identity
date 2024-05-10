document.addEventListener("DOMContentLoaded", function () {
  if (!window.PublicKeyCredential) {
    alert("Client not capable. Handle error.");
  }
});

// Base64 to ArrayBuffer
function bufferDecode(value) {
  value = value.replace(/-/g, "+").replace(/_/g, "/");
  return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
}

// ArrayBuffer to URLBase64
function bufferEncode(value) {
  return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");
}

function registerUser() {
  fetch(`/webauthn/register/begin`)
    .then(function (response) {
      return response.json();
    })
    .then((credentialCreationOptions) => {
      credentialCreationOptions.publicKey.challenge = bufferDecode(
        credentialCreationOptions.publicKey.challenge
      );
      credentialCreationOptions.publicKey.user.id = bufferDecode(
        credentialCreationOptions.publicKey.user.id
      );

      if (credentialCreationOptions.publicKey.excludeCredentials) {
        credentialCreationOptions.publicKey.excludeCredentials.forEach(
          (item) => {
            item.id = bufferDecode(item.id);
          }
        );
      }

      return navigator.credentials.create({
        publicKey: credentialCreationOptions.publicKey,
      });
    })
    .then((credential) => {
      let attestationObject = credential.response.attestationObject;
      let clientDataJSON = credential.response.clientDataJSON;
      let rawId = credential.rawId;

      const csrf = getCSRF();

      fetch(`/webauthn/register/finish`, {
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
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
          },
        }),
      });
    })
    .then(() => {
      document.getElementById("notification").innerHTML =
        "Successfully registered.";
    })
    .catch((err) => {
      console.warn(err);
      document.getElementById("notification").innerHTML =
        "Registration failed.";
    });
}

async function registerUser2(returnUrl) {
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
    beginResponse.publicKey.user.id = bufferDecode(
      beginResponse.publicKey.user.id
    );

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

function LoginUser() {
  fetch(`/webauthn/login/begin`)
    .then(function (response) {
      return response.json();
    })
    .then((credentialRequestOptions) => {
      credentialRequestOptions.publicKey.challenge = bufferDecode(
        credentialRequestOptions.publicKey.challenge
      );
      credentialRequestOptions.publicKey.allowCredentials.forEach(function (
        listItem
      ) {
        listItem.id = bufferDecode(listItem.id);
      });

      return navigator.credentials.get({
        publicKey: credentialRequestOptions.publicKey,
      });
    })
    .then((assertion) => {
      let authData = assertion.response.authenticatorData;
      let clientDataJSON = assertion.response.clientDataJSON;
      let rawId = assertion.rawId;
      let sig = assertion.response.signature;
      let userHandle = assertion.response.userHandle;

      const csrf = getCSRF();

      fetch(`/webauthn/login/finish`, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
          "X-Csrf-Token": csrf,
        },
        body: JSON.stringify({
          id: assertion.id,
          rawId: bufferEncode(rawId),
          type: assertion.type,
          response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
          },
        }),
      })
        .then((response) => {
          return response.json();
        })
        .then((data) => {
          /* 
       document.getElementById("notification").innerHTML =
          "Successfully logged in.";
          */
          window.location.href = data.redirectUri;
        });
    })

    .then((data) => {
      /* 
     document.getElementById("notification").innerHTML =
        "Successfully logged in.";
        */
      window.location.href = data.redirectUri;
    })
    .catch((error) => {
      console.log(error);
      document.getElementById("notification").innerHTML = "Login failed.";
    });
}

async function LoginUser22() {
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

    const credential = await navigator.credentials.get(beginResponse);
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
