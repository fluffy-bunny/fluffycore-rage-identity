fetch("/api/manifest", {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      "X-Csrf-Token": csrf,
    },
    
  }).then((response) => {
    return response.json();
  }).then((data) => {
    console.log("finishResponse:", data);
  });