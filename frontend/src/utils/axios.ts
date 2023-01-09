import axios from "axios";
import { ServerResponse } from "../models/server";

const BASE_URL = "http://localhost:3001/api/";

const axiosConfig = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

async function refreshToken() {
  const res = (await axiosConfig.get<ServerResponse<any>>("auth/refresh")).data;

  if (res.error) {
    throw Error("Error refreshing access token");
  } else {
    return res;
  }
}

axiosConfig.interceptors.response.use(
  (res) => {
    return res;
  },
  async (error) => {
    const originalReq = error.config;
    const errMsg = error.response.data.message as string;

    if (errMsg.includes("") && !originalReq._retry) {
      originalReq._retry = true;
      await refreshToken();
      return axiosConfig(originalReq);
    }

    return Promise.reject(error);
  }
);

export { axiosConfig };
