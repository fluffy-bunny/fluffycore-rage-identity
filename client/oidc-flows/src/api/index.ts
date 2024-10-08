import { AxiosError } from 'axios';

import { getCSRF } from '../utils/cookies';
import { Api } from './Api';

export const { api, externalIdp, forgotPassword, instance } = new Api({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
  withXSRFToken: true,
});

instance.interceptors.request.use((config) => {
  const xcrfToken = getCSRF();

  if (xcrfToken) {
    config.headers['X-Csrf-Token'] = xcrfToken;
  }

  return config;
});

export const isApiError = (error: any): error is AxiosError => {
  return error.isAxiosError;
};

export type ApiError<T> = AxiosError<T>;
