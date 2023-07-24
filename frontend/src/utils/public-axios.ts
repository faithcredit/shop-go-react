import axios from 'axios';
// import { baseUrl } from './helper';

const publicAxios = axios.create({
  baseURL: `${process.env.REACT_APP_SERVER_ENDPOINT}/api`,
});

export default publicAxios;
