import { AxiosError } from 'axios';

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

export const isApiError = (error: any): error is AxiosError => {
  return error.isAxiosError;
};
