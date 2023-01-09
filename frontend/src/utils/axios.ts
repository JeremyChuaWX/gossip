import axios from "axios";

const BASE_URL = "http://localhost:3001/api/";

const axiosConfig = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

export { axiosConfig };
