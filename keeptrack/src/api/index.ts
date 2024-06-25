import { Api } from "./Api";

export const { api, externalIdp, forgotPassword, instance } = new Api({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
  withXSRFToken: true,
});
