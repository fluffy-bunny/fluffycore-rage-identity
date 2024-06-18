import { Api } from "./Api";

export const { api, externalIdp, instance } = new Api({
  baseURL: "http://localhost1.com:9044",
  withCredentials: true,
  withXSRFToken: true,
});
