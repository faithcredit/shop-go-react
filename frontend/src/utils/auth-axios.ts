import axios from 'axios';
// import { baseUrl } from './helper';

const authAxios = axios.create({
  baseURL: `${process.env.REACT_APP_SERVER_ENDPOINT}/api`,
});

export const authorizationProvider = (store: any) => {
  authAxios.interceptors.request.use((config: any) => {
    const token = store.getState().login.userInfo.token;
    config.headers.Authorization = `Bearer ${token}`;

    // response.setHeader("Access-Control-Allow-Origin", "*");
    // config.setHeader("Access-Control-Allow-Credentials", "true");
    // config.setHeader("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT");
    // config.setHeader("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers");
    return config;
  });
};

export default authAxios;
