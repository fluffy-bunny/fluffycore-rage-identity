// common.js
function getCookieValue(name) {
  // Create a regular expression to match the cookie name
  const nameEQ = name + "=";
  const cookies = document.cookie.split(';');
  
  for (let i = 0; i < cookies.length; i++) {
      let cookie = cookies[i];
      // Remove leading spaces
      while (cookie.charAt(0) === ' ') {
          cookie = cookie.substring(1, cookie.length);
      }
      // Check if the cookie name matches
      if (cookie.indexOf(nameEQ) === 0) {
          return cookie.substring(nameEQ.length, cookie.length);
      }
  }
  return null;
}

function getCSRF() {
  let csrf = getCookieValue('_csrf');
  return csrf;
}
 