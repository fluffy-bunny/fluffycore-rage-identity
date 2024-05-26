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

  fetch('/api/verify-username', {
    method: 'POST', 
    credentials: 'include', 
    headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'X-Csrf-Token': csrf
                },
    body: JSON.stringify({
        username: 'ghstahl@gmail.com'
    })
})
.then(response => response.json())
.then(data => console.log(data))
.catch((error) => {
console.error('Error:', error);
});