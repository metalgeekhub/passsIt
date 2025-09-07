import axios from 'axios';

const publicAxios = axios.create({
  baseURL: 'http://localhost:8080', // your backend URL
  withCredentials: false, 
});

const authAxios = axios.create({
  baseURL: 'http://localhost:8080', // your backend URL
  withCredentials: true,            // send cookies with requests
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