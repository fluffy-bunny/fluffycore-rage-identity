import { useEffect, useState } from 'react';

export const usePublicKeyCredentialSupport = () => {
  const [isSupported, setIsSupported] = useState<boolean>();

  useEffect(() => {
    const checkSupport = () => {
      if (window.PublicKeyCredential) {
        console.log('PublicKeyCredential is supported in this browser.');
        setIsSupported(true);
      } else {
        console.log('PublicKeyCredential is not supported in this browser.');
        setIsSupported(false);
      }
    };

    checkSupport();
  }, []);

  return isSupported;
};
