export function getCookieValue(name: string) {
  const nameString = name + '=';

  const value = document.cookie.split(';').filter((item) => {
    return item.includes(nameString);
  });

  if (value.length) {
    return value[0].substring(nameString.length, value[0].length);
  } else {
    return '';
  }
}

export function getCSRF() {
  return getCookieValue('_csrf');
}
