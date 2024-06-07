import { getCSRF } from "../utils/cookies";
import { Api } from "./Api";

export const { api, instance } = new Api({
  baseURL: "http://localhost1.com:9044",
});
