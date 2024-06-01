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
    email: "ghstahl@gmail.com",
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
    password: "1234"
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
    password: "http://localhost:9044/signup"
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Error:", error);
  });
 