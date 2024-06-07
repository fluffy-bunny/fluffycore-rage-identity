import { getCSRF } from "../utils/cookies";
import { Api } from "./Api";

export const { api, instance } = new Api({
  baseURL: "http://localhost1.com:9044",
});

instance.interceptors.request.use((config) => {
  const csrf = getCSRF();

  if (csrf) {
    config.headers["X-Csrf-Token"] = csrf;
  }

  return config;
});
