import axios from "axios";
import { ServerResponse } from "../models/server";

const BASE_URL = "http://localhost:3001/api/";

const axiosConfig = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

async function refreshToken() {
  const res = (await axiosConfig.get<ServerResponse<any>>("auth/refresh")).data;
  return res;
}

axiosConfig.interceptors.request.use(
  (config) => {
    console.info(`config: ${config.url}`, config);
    return config;
  },
  (error) => {
    console.error(`error: ${error.config.url}`, error);
    return Promise.reject(error);
  }
);

axiosConfig.interceptors.response.use(
  (config) => config,
  async (error) => {
    console.error(`error: ${error.config.url}`, error);

    const originalReq = error.config;
    const errMsg = error.response.data.msg as string;

    if (errMsg.includes("Invalid access token") && !originalReq._retry) {
      originalReq._retry = true;
      await refreshToken();

      return axiosConfig(originalReq);
    }

    return Promise.reject(error);
  }
);

export { axiosConfig };
