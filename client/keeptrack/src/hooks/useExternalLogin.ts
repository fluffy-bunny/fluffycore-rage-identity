import { api } from '../api';

export function useExternalLogin() {
  const execute = async (slug?: string) => {
    if (!slug) {
      return;
    }

    const { data: externalLoginData } = await api.startExternalLoginCreate({
      slug,
      directive: 'login',
    });

    if (externalLoginData.redirectUri) {
      window.location.href = externalLoginData.redirectUri;
    }
  };

  return [execute];
}
