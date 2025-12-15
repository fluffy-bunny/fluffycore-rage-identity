document.addEventListener("DOMContentLoaded", function () {
  if (!window.PublicKeyCredential) {
    console.error("WebAuthn not supported in this browser");
    alert("Your browser does not support passkeys. Please use a modern browser like Chrome, Safari, or Edge.");
  }
});

// Base64URL to ArrayBuffer (RFC 4648 base64url)
function bufferDecode(value) {
  // Add padding if needed
  const padding = '='.repeat((4 - (value.length % 4)) % 4);
  const base64 = value.replace(/-/g, "+").replace(/_/g, "/") + padding;
  
  try {
    const rawData = atob(base64);
    const outputArray = new Uint8Array(rawData.length);
    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
  } catch (e) {
    console.error("Failed to decode base64url:", e);
    throw new Error("Invalid base64url string");
  }
}

// ArrayBuffer to Base64URL
function bufferEncode(value) {
  const bytes = new Uint8Array(value);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");
}

// Get user-friendly error messages
function getWebAuthnErrorMessage(error) {
  if (!error) return "An unknown error occurred";
  
  const errorName = error.name || "";
  
  switch (errorName) {
    case "NotAllowedError":
      return "The operation was not allowed. You may have canceled the request or it timed out.";
    case "InvalidStateError":
      return "This passkey is already registered.";
    case "NotSupportedError":
      return "This authenticator is not supported.";
    case "SecurityError":
      return "Security error. Please ensure you're on a secure connection (HTTPS).";
    case "AbortError":
      return "The operation was aborted.";
    case "ConstraintError":
      return "The authenticator doesn't meet the requirements.";
    case "NetworkError":
      return "Network error. Please check your connection.";
    default:
      return error.message || "An error occurred during the operation.";
  }
}

async function registerUser(returnUrl, friendlyName) {
  let abortController = null;
  
  try {
    // Fetch registration options from server
    const response = await fetch("/webauthn/register/begin");

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Server error: ${response.status} - ${errorText}`);
    }
    
    const optionsJSON = await response.json();
    console.log("Registration options:", optionsJSON);

    // Use parseCreationOptionsFromJSON if available (recommended by Google/W3C)
    // This handles base64url decoding automatically
    let options;
    if (PublicKeyCredential.parseCreationOptionsFromJSON) {
      options = PublicKeyCredential.parseCreationOptionsFromJSON(optionsJSON);
    } else {
      // Fallback: manual decoding for older browsers
      options = optionsJSON;
      options.publicKey.challenge = bufferDecode(options.publicKey.challenge);
      options.publicKey.user.id = bufferDecode(options.publicKey.user.id);
      
      if (options.publicKey.excludeCredentials) {
        options.publicKey.excludeCredentials.forEach((cred) => {
          cred.id = bufferDecode(cred.id);
        });
      }
    }

    // Create abort controller for timeout
    abortController = new AbortController();
    const timeoutId = setTimeout(() => {
      abortController.abort();
    }, 120000); // 2 minute timeout

    // Create credential
    const credentialsResponse = await navigator.credentials.create({
      publicKey: options.publicKey || options,
      signal: abortController.signal,
    });

    clearTimeout(timeoutId);

    if (!credentialsResponse) {
      throw new Error("Failed to create credential");
    }

    console.log("Credential created:", credentialsResponse);

    // Get transports if available (improves future authentication UX)
    const transports = credentialsResponse.response.getTransports 
      ? credentialsResponse.response.getTransports() 
      : [];

    const csrf = getCSRF();

    // Use toJSON() for automatic encoding (recommended by Google/W3C)
    let credentialJSON;
    if (credentialsResponse.toJSON) {
      credentialJSON = credentialsResponse.toJSON();
      // Add transports and friendly name
      credentialJSON.response.transports = transports;
      credentialJSON.friendlyName = friendlyName || "My Passkey";
    } else {
      // Fallback: manual encoding for older browsers
      credentialJSON = {
        id: credentialsResponse.id,
        rawId: bufferEncode(credentialsResponse.rawId),
        type: credentialsResponse.type,
        response: {
          attestationObject: bufferEncode(credentialsResponse.response.attestationObject),
          clientDataJSON: bufferEncode(credentialsResponse.response.clientDataJSON),
          transports: transports,
        },
        friendlyName: friendlyName || "My Passkey",
      };
    }

    // Send to server
    const finishResponse = await fetch("/webauthn/register/finish", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        "X-Csrf-Token": csrf,
      },
      body: JSON.stringify(credentialJSON),
    });

    if (!finishResponse.ok) {
      const errorData = await finishResponse.json().catch(() => ({}));
      throw new Error(errorData.error || `Server error: ${finishResponse.status}`);
    }

    const data = await finishResponse.json();
    console.log("Registration complete:", data);
    
    // Redirect on success
    window.location.href = returnUrl;
  } catch (error) {
    console.error("Registration error:", error);
    
    // Handle InvalidStateError specially (passkey already exists)
    if (error.name === "InvalidStateError") {
      alert("You already have a passkey registered on this device.");
      // Still redirect since the user has a passkey
      if (returnUrl) {
        window.location.href = returnUrl;
      }
      return;
    }
    
    // User-friendly error message
    const message = getWebAuthnErrorMessage(error);
    alert(`Passkey registration failed: ${message}`);
    
    // Don't redirect on error, let user retry
  } finally {
    if (abortController) {
      abortController = null;
    }
  }
}

async function LoginUser(returnFailedUrl, useConditionalUI = false) {
  let abortController = null;
  
  try {
    // Fetch authentication options from server
    const response = await fetch("/webauthn/login/begin");
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Server error: ${response.status} - ${errorText}`);
    }
    
    const optionsJSON = await response.json();
    console.log("Authentication options:", optionsJSON);

    // Use parseRequestOptionsFromJSON if available (recommended by Google/W3C)
    // This handles base64url decoding automatically
    let options;
    if (PublicKeyCredential.parseRequestOptionsFromJSON) {
      options = PublicKeyCredential.parseRequestOptionsFromJSON(optionsJSON);
    } else {
      // Fallback: manual decoding for older browsers
      options = optionsJSON;
      options.publicKey.challenge = bufferDecode(options.publicKey.challenge);
      
      if (options.publicKey.allowCredentials) {
        options.publicKey.allowCredentials.forEach((allowCredential) => {
          allowCredential.id = bufferDecode(allowCredential.id);
        });
      }
    }

    // Create abort controller for timeout
    abortController = new AbortController();
    const timeoutId = setTimeout(() => {
      abortController.abort();
    }, 120000); // 2 minute timeout

    // Configure credential request options
    const credentialRequestOptions = {
      publicKey: options.publicKey || options,
      signal: abortController.signal,
    };

    // Add conditional UI (autofill) support if requested and available
    if (useConditionalUI && window.PublicKeyCredential.isConditionalMediationAvailable) {
      const conditionalAvailable = await PublicKeyCredential.isConditionalMediationAvailable();
      if (conditionalAvailable) {
        credentialRequestOptions.mediation = 'conditional';
        console.log("Using conditional UI for autofill");
      }
    }

    // Get credential
    const credential = await navigator.credentials.get(credentialRequestOptions);
    
    clearTimeout(timeoutId);

    if (!credential) {
      throw new Error("Failed to get credential");
    }

    console.log("Credential retrieved:", credential);

    const csrf = getCSRF();

    // Use toJSON() for automatic encoding (recommended by Google/W3C)
    let credentialJSON;
    if (credential.toJSON) {
      credentialJSON = credential.toJSON();
    } else {
      // Fallback: manual encoding for older browsers
      credentialJSON = {
        id: credential.id,
        rawId: bufferEncode(credential.rawId),
        type: credential.type,
        response: {
          authenticatorData: bufferEncode(credential.response.authenticatorData),
          clientDataJSON: bufferEncode(credential.response.clientDataJSON),
          signature: bufferEncode(credential.response.signature),
          userHandle: credential.response.userHandle ? bufferEncode(credential.response.userHandle) : null,
        },
      };
    }

    // Send to server
    const finishResponse = await fetch("/webauthn/login/finish", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        "X-Csrf-Token": csrf,
      },
      body: JSON.stringify(credentialJSON),
    });

    if (!finishResponse.ok) {
      const errorData = await finishResponse.json().catch(() => ({}));
      
      // Signal unknown credential to passkey provider (per Google recommendations)
      if (finishResponse.status === 404 && PublicKeyCredential.signalUnknownCredential) {
        try {
          await PublicKeyCredential.signalUnknownCredential({
            rpId: window.location.hostname,
            credentialId: credential.id,
          });
          console.log("Signaled unknown credential to passkey provider");
        } catch (signalError) {
          console.warn("Failed to signal unknown credential:", signalError);
        }
      }
      
      throw new Error(errorData.error || `Server error: ${finishResponse.status}`);
    }

    const data = await finishResponse.json();
    console.log("Authentication complete:", data);
    
    // Redirect on success
    window.location.href = data.redirectUri || returnFailedUrl;
  } catch (error) {
    console.error("Authentication error:", error);
    
    // User-friendly error message
    const message = getWebAuthnErrorMessage(error);
    alert(`Passkey authentication failed: ${message}`);
    
    // Redirect to failed URL
    if (returnFailedUrl) {
      window.location.href = returnFailedUrl;
    }
  } finally {
    if (abortController) {
      abortController = null;
    }
  }
}

// Check if conditional UI (autofill) is available
async function isConditionalUIAvailable() {
  if (window.PublicKeyCredential && PublicKeyCredential.isConditionalMediationAvailable) {
    return await PublicKeyCredential.isConditionalMediationAvailable();
  }
  return false;
}

// Check if user verifying platform authenticator is available
async function isPlatformAuthenticatorAvailable() {
  if (window.PublicKeyCredential && PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable) {
    return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
  }
  return false;
}
