$(document).ready(function () {
  if (!window.PublicKeyCredential) {
    alert("Client not capable. Handle error.");
    return;
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

function registerUser(returnUrl) {
  $.get(
    "/webauthn/register/begin",
    null,
    function (data) {
      return data;
    },
    "json"
  )
    .then((beginResponse) => {
      console.log("beginResponse:", beginResponse);

      beginResponse.publicKey.challenge = bufferDecode(
        beginResponse.publicKey.challenge
      );
      beginResponse.publicKey.user.id = bufferDecode(
        beginResponse.publicKey.user.id
      );

      //if (beginResponse.publicKey.excludeCredentials) {
      //  for (var i = 0; i < beginResponse.publicKey.excludeCredentials.length; i++) {
      //    beginResponse.publicKey.excludeCredentials[i].id = bufferDecode(beginResponse.publicKey.excludeCredentials[i].id);
      //  }
      //}

      return navigator.credentials.create({
        publicKey: beginResponse.publicKey,
      });
    })
    .then((credentialsResponse) => {
      console.log("credentialsResponse:", credentialsResponse);

      let attestationObject = credentialsResponse.response.attestationObject;
      let clientDataJSON = credentialsResponse.response.clientDataJSON;
      let rawId = credentialsResponse.rawId;
      let csrf = getCSRF();

      fetch("/webauthn/register/finish", {
        method: "POST",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'X-Csrf-Token': csrf
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
      })
        .then(response => response.json())
        .then(finishResponse => {
          console.log("finishResponse:", finishResponse);
          window.location.href = returnUrl;
        })
        .catch(error => {
          console.error("Error:", error);
          // Handle errors
        });

      
    })
    .then((success) => {
      console.log("success:", success);
      alert("successfully registered !");
    })
    .catch((error) => {
      console.log("failed:", error);
      alert("occur exception");
    });
}

function LoginUser() {
  $.get(
    "/webauthn/login/begin",
    null,
    function (data) {
      return data;
    },
    "json"
  )
    .then((beginResponse) => {
      console.log("beginResponse: ", beginResponse);

      beginResponse.publicKey.challenge = bufferDecode(
        beginResponse.publicKey.challenge
      );
      beginResponse.publicKey.allowCredentials.forEach(function (
        allowCredential
      ) {
        allowCredential.id = bufferDecode(allowCredential.id);
      });

      return navigator.credentials.get({
        publicKey: beginResponse.publicKey,
      });
    })
    .then((getCredential) => {
      console.log("getCredential: ", getCredential);

      let authData = getCredential.response.authenticatorData;
      let clientDataJSON = getCredential.response.clientDataJSON;
      let rawId = getCredential.rawId;
      let sig = getCredential.response.signature;
      let userHandle = getCredential.response.userHandle;
      let csrf = getCSRF();

      fetch("/webauthn/login/finish", {
        method: "POST",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'X-Csrf-Token': csrf
        },
        body: JSON.stringify({
          id: getCredential.id,
          rawId: bufferEncode(rawId),
          type: getCredential.type,
          response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
          },
        }),
      })
        .then(response => response.json())
        .then(finishResponse => {
          console.log("finishResponse:", finishResponse);
          window.location.href = finishResponse.redirectUri;
        })
        .catch(error => {
          console.error("Error:", error);
          // Handle errors
        });
    })
    .then((success) => {
      alert("successfully logged in " + username + "!");
      return;
    })
    .catch((error) => {
      console.log(error);
      alert("failed to login " + username);
    });
}
