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

function getAllCookies() {
    return document.cookie;
}

async function sendRequestWithCookies(url, method, body) {
    try {
        const csrfToken = getCSRF();

        const requestOptions = {
            method: method,
            credentials: 'include', // Ensures cookies are included
            headers: {
                'Content-Type': 'application/json',
                'X-Csrf-Token': csrfToken
            }
        };

        if (body) {
            requestOptions.body = JSON.stringify(body);
        }

        const response = await fetch(url, requestOptions);

        let wrappedResponse = {
            statusCode: response.status
        };
        wrappedResponse.response = await response.json();


        return wrappedResponse;
    } catch (error) {
        let wrappedResponse = {
            status: 500,
            error: error
        };

        console.error('Error sending request: ', error);
        return wrappedResponse;
    }
}
window.setFocus = (element) => {
    element.focus();
};

