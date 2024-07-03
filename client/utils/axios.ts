import { API_URL } from "@/constants";
import axios from "axios";

// Create an instance of axios
const axiosInstance = axios.create({
  baseURL: API_URL,
  timeout: 10000,
});

// Request interceptor to add the authorization header
axiosInstance.interceptors.request.use(
  async (config) => {
    const accessToken = localStorage.getItem("accessToken");
    if (!accessToken) {
      return config;
    }

    config.headers["Authorization"] = `Bearer ${accessToken}`;
    return config;
  },
  (error) => {
    // We can add custom error handling here
    return Promise.reject(error);
  }
);

export default axiosInstance;
