// WebAuthn helper functions for passkey management
// This file must be loaded before the WASM application
console.log("🔵 webauthn.js: Script started loading...");

// Get CSRF token from cookie
function getCSRF() {
  const name = "_csrf=";
  const decodedCookie = decodeURIComponent(document.cookie);
  const ca = decodedCookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) === 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

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

async function LoginUser(returnFailedUrl, useConditionalUI = false, errorCallback = null) {
  let abortController = null;
  
  console.log("LoginUser called with errorCallback:", errorCallback, "type:", typeof errorCallback);
  
  try {
    // Fetch authentication options from server
    // The server will automatically clear any stale signin cookies
    const response = await fetch("/webauthn/login/begin");
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Server error: ${response.status} - ${errorText}`);
    }
    
    const optionsJSON = await response.json();

    // Use parseRequestOptionsFromJSON if available (recommended by Google/W3C)
    // This handles base64url decoding automatically
    let credentialRequestOptions;
    if (PublicKeyCredential.parseRequestOptionsFromJSON) {
      // parseRequestOptionsFromJSON expects the publicKey object directly, not wrapped
      const publicKeyOptions = optionsJSON.publicKey || optionsJSON;
      const parsedOptions = PublicKeyCredential.parseRequestOptionsFromJSON(publicKeyOptions);
      
      // parseRequestOptionsFromJSON returns the publicKey options, wrap it for navigator.credentials.get
      credentialRequestOptions = {
        publicKey: parsedOptions
      };
    } else {
      // Fallback: manual decoding for older browsers
      const options = optionsJSON;
      options.publicKey.challenge = bufferDecode(options.publicKey.challenge);
      
      if (options.publicKey.allowCredentials) {
        options.publicKey.allowCredentials.forEach((allowCredential) => {
          allowCredential.id = bufferDecode(allowCredential.id);
        });
      }
      
      credentialRequestOptions = {
        publicKey: options.publicKey,
      };
    }

    // Create abort controller for timeout
    const abortController = new AbortController();
    const timeoutId = setTimeout(() => {
      abortController.abort();
    }, 120000); // 2 minute timeout

    // Add signal to options
    credentialRequestOptions.signal = abortController.signal;

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
    
    // Redirect on success (check both redirectUrl and redirectUri for compatibility)
    const redirectUrl = data.redirectUrl || data.redirectUri || returnFailedUrl;
    window.location.href = redirectUrl;
  } catch (error) {
    console.error("Authentication error:", error);
    
    // User-friendly error message
    const message = getWebAuthnErrorMessage(error);
    
    // Use callback if provided (from WASM), otherwise use alert
    if (errorCallback && typeof errorCallback === 'function') {
      errorCallback(message);
    } else {
      alert(`Passkey authentication failed: ${message}`);
      
      // Redirect to failed URL only if using alert (not callback)
      if (returnFailedUrl) {
        window.location.href = returnFailedUrl;
      }
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

// Register a new passkey (for management page - no redirect)
async function registerPasskey(friendlyName) {
  console.log("🟢 registerPasskey() called with friendlyName:", friendlyName);
  console.log("🟢 Function type:", typeof registerPasskey);
  console.log("🟢 Window object has registerPasskey:", "registerPasskey" in window);
  
  let abortController = null;
  
  try {
    console.log("🔵 Step 1: Fetching registration options from /webauthn/register/begin");
    // Fetch registration options from server
    const response = await fetch("/webauthn/register/begin");

    if (!response.ok) {
      const errorText = await response.text();
      console.error("❌ Server error:", response.status, errorText);
      throw new Error(`Server error: ${response.status} - ${errorText}`);
    }
    
    const optionsJSON = await response.json();
    console.log("✅ Step 1 complete: Registration options received:", optionsJSON);
    console.log("🔍 optionsJSON structure:", JSON.stringify(optionsJSON, null, 2));
    console.log("🔍 optionsJSON.publicKey:", optionsJSON.publicKey);
    console.log("🔍 optionsJSON.publicKey.challenge:", optionsJSON.publicKey?.challenge);

    console.log("🔵 Step 2: Parsing credential creation options");
    // Use parseCreationOptionsFromJSON if available (recommended by Google/W3C)
    let options;
    if (PublicKeyCredential.parseCreationOptionsFromJSON) {
      // Pass optionsJSON.publicKey directly, not the wrapper object
      options = PublicKeyCredential.parseCreationOptionsFromJSON(optionsJSON.publicKey);
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
    console.log("✅ Step 2 complete: Options parsed successfully");

    console.log("🔵 Step 3: Creating credential (browser will show dialog)");
    // Create abort controller for timeout
    abortController = new AbortController();
    const timeoutId = setTimeout(() => {
      console.log("⏱️ Timeout reached, aborting...");
      abortController.abort();
    }, 120000); // 2 minute timeout

    // Create credential
    const credentialsResponse = await navigator.credentials.create({
      publicKey: options.publicKey || options,
      signal: abortController.signal,
    });

    clearTimeout(timeoutId);

    if (!credentialsResponse) {
      console.error("❌ Failed to create credential - no response");
      throw new Error("Failed to create credential");
    }

    console.log("✅ Step 3 complete: Credential created:", credentialsResponse);

    console.log("🔵 Step 4: Preparing credential data for server");
    // Get transports if available
    const transports = credentialsResponse.response.getTransports 
      ? credentialsResponse.response.getTransports() 
      : [];

    const csrf = getCSRF();

    // Use toJSON() for automatic encoding
    let credentialJSON;
    if (credentialsResponse.toJSON) {
      credentialJSON = credentialsResponse.toJSON();
      credentialJSON.response.transports = transports;
      credentialJSON.friendlyName = friendlyName || "My Passkey";
    } else {
      // Fallback: manual encoding
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
      console.error("❌ Server rejected credential:", finishResponse.status, errorData);
      throw new Error(errorData.error || `Server error: ${finishResponse.status}`);
    }

    const data = await finishResponse.json();
    console.log("✅ Step 5 complete: Server accepted credential:", data);
    console.log("✅✅✅ PASSKEY REGISTRATION SUCCESSFUL! ✅✅✅");
    
    // Return success
    return true;
  } catch (error) {
    console.error("❌❌❌ PASSKEY REGISTRATION FAILED:", error);
    console.error("Error name:", error.name);
    console.error("Error message:", error.message);
    console.error("Error stack:", error.stack);
    
    // Handle InvalidStateError specially (passkey already exists)
    if (error.name === "InvalidStateError") {
      console.warn("⚠️ Passkey already exists on this device");

      alert("You already have a passkey registered on this device.");
      return false;
    }
    
    // User-friendly error message
    const message = getWebAuthnErrorMessage(error);
    alert(`Passkey registration failed: ${message}`);
    return false;
  } finally {
    if (abortController) {
      abortController = null;
    }
  }
}

console.log("🔵 webauthn.js: registerPasskey function defined");
console.log("🔵 webauthn.js: typeof registerPasskey =", typeof registerPasskey);

// Ensure all functions are immediately available on window object
// These assignments must happen synchronously during script load
console.log("🔵 webauthn.js: Assigning functions to window object...");
window.registerPasskey = registerPasskey;
window.registerUser = registerUser;
window.LoginUser = LoginUser;

// Log when script is fully loaded
console.log("✅✅✅ webauthn.js: SCRIPT FULLY LOADED ✅✅✅");
console.log("✅ All functions registered on window object:");
console.log("  - window.registerPasskey:", typeof window.registerPasskey);
console.log("  - window.registerUser:", typeof window.registerUser);  
console.log("  - window.LoginUser:", typeof window.LoginUser);
console.log("✅ Script ready - WASM can now call these functions");

// Check browser support after DOM loads
document.addEventListener("DOMContentLoaded", function () {
  if (!window.PublicKeyCredential) {
    console.error("WebAuthn not supported in this browser");
  } else {
    console.log("✅ WebAuthn supported in this browser");
  }
});