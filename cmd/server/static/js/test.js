// /api/manifest
//-----------------------------------------------------
fetch("/api/manifest", {
  method: "GET",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
})
  .then((response) => {
    return response.json();
  })
  .then((data) => {
    console.log("finishResponse:", data);
  });

// /api/start-external-login
//-----------------------------------------------------
fetch("/api/start-external-login", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    slug: "google-social",
    directive: "login",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(JSON.stringify(data)))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/verify-username
//-----------------------------------------------------
fetch("/api/verify-username", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    username: "ghstahl@gmail.com",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/verify-password-strength
//-----------------------------------------------------
fetch("/api/verify-password-strength", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@gmail.com",
    password: "ghstahl@gmail.com",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/verify-password-strength
//-----------------------------------------------------
fetch("/api/verify-password-strength", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    password: "ghstahl@gmail.com_1234567890abcdefghij",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/login-phase-one
//-----------------------------------------------------
fetch("/api/login-phase-one", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@gmail.com",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/login-phase-one for claimed domain
//-----------------------------------------------------
fetch("/api/login-phase-one", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@mapped.com",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });
// /api/start-external-login for mapped-enterprise
//-----------------------------------------------------
fetch("/api/start-external-login", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    slug: "mapped-enterprise",
    directive: "login",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(JSON.stringify(data)))
  .catch((error) => {
    console.error("Error:", error);
  });
// /api/login-password
//-----------------------------------------------------
fetch("/api/login-password", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@gmail.com",
    password: "1234",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/login-password
//-----------------------------------------------------
fetch("/api/login-password", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@gmail.com",
    password: "http://localhost:9044/signup",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/verify-code
//-----------------------------------------------------
fetch("/api/verify-code", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    code: "zwDnOR",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/signup
//-----------------------------------------------------
fetch("/api/signup", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@gmail.com",
    password: "http://localhost:9044/signup",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/password-reset-start
//-----------------------------------------------------
fetch("/api/password-reset-start", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    email: "ghstahl@gmail.com",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/password-reset-finish
//-----------------------------------------------------
fetch("/api/password-reset-finish", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    password: "http://localhost:9044/signup",
    passwordConfirm: "http://localhost:9044/signup",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/user-identity-info
//-----------------------------------------------------
fetch("/api/user-identity-info", {
  method: "GET",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
})
  .then((response) => {
    return response.json();
  })
  .then((data) => {
    console.log("finishResponse:", data);
  });

//   /api/user-profile
//-----------------------------------------------------
fetch("/api/user-profile", {
  method: "GET",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
})
  .then((response) => {
    return response.json();
  })
  .then((data) => {
    console.log("finishResponse:", data);
  });

fetch("/api/user-profile", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    givenName: 'bugs', familyName: 'bunny', phoneNumber: '555-1212'
  }),
});

// /api/password-reset-finish
//-----------------------------------------------------
fetch("/api/logout", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });

// /api/user-linked-accounts
//-----------------------------------------------------
fetch("/api/user-linked-accounts", {
  method: "GET",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
})
  .then((response) => {
    return response.json();
  })
  .then((data) => {
    console.log("finishResponse:", data);
  });

  // /api/user-linked-accounts
//-----------------------------------------------------
fetch("/api/user-linked-accounts?identity=google-social", {
  method: "DELETE",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
})
  .then((response) => {
    return response.json();
  })
  .then((data) => {
    console.log("finishResponse:", data);
  });

  // /api/user-remove-passkey
//-----------------------------------------------------
fetch("/api/user-remove-passkey", {
  method: "POST",
  credentials: "include",
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Csrf-Token": getCSRF(),
  },
  body: JSON.stringify({
    aaguid: 'AAAAAAAAAAAAAAAAAAAAAA==' 
  }),
}) .then((response) => {
  return response.json();
})
.then((data) => {
  console.log("finishResponse:", data);
});