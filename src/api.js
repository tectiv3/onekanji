import Vue from 'vue'
import axios from 'axios';

axios.defaults.headers.post['Content-Type'] = 'application/json';

export const baseURL = 'http://localhost:7780/api';

const api = axios.create({
    baseURL,
});

api.interceptors.response.use(function ({data}) {
    return data;
}, function (error) {
    return Promise.reject(error);
});

export default api;
