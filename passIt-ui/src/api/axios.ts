import axios from 'axios';
import config from '../config/env';

const publicAxios = axios.create({
  baseURL: config.apiBaseUrl,
  withCredentials: false, 
});

const authAxios = axios.create({
  baseURL: config.apiBaseUrl,
  withCredentials: true, // send cookies with requests
});

// Interceptor for handling 401 Unauthorized globally
// authAxios.interceptors.response.use(
//   response => response,
//   error => {
//     if (error.response && error.response.status === 401) {
//       window.location.href = '/login';
//     }
//     return Promise.reject(error);
//   }
// );

export default authAxios;
export { publicAxios };