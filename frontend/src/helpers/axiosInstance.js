import axios from "axios";

export default () => {
  const baseURL = process.env.REACT_APP_BE_URL;

  let headers = {};

  if (localStorage.access_token) {
    headers.Authorization = `Bearer ${localStorage.access_token}`;
  }

  const axiosInstance = axios.create({
    baseURL: baseURL,
    headers,
  });

  axiosInstance.interceptors.response.use(
    (response) =>
      new Promise((resolve, reject) => {
        resolve(response);
      }),
    (error) => {
      return new Promise((resolve, reject) => {
        reject(error);
      });
      // if (!error.response) {
      //   return new Promise((resolve, reject) => {
      //     reject(error);
      //   });
      // }

      // if (error.response.status === 401) {
      //   localStorage.removeItem("access_token");

      //   window.location = "/login";

      // } else {
      //   return new Promise((resolve, reject) => {
      //     reject(error);
      //   });
      // }
    }
  );

  return axiosInstance;
};
