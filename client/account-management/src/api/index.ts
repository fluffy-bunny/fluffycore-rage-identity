import { AxiosError } from 'axios';

import { getCSRF } from '../utils/cookies';
import { Api } from './Api';

export const {
  api,
  externalIdp,
  forgotPassword,
  instance: apiInstance,
} = new Api({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
  withXSRFToken: true,
});

apiInstance.interceptors.request.use((config) => {
  const xcrfToken = getCSRF();

  if (xcrfToken) {
    config.headers['X-Csrf-Token'] = xcrfToken;
  }

  return config;
});

export const isApiError = (error: any): error is AxiosError => {
  return error.isAxiosError;
};
