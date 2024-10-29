import axios from "axios";

// const baseURL = "http://defo1215.cn:3030/api";
const baseURL = "http://localhost:3030/api";


const request = axios.create({
    baseURL,
    timeout: 1000,
    headers: {
        "Content-Type": "application/json",
    },
});

// 请求拦截器
request.interceptors.request.use(
    (config) => { return config },
    (error) => { return Promise.reject(error); }
);

// 响应拦截器
request.interceptors.response.use(
    (response) => { return response; },
    (error) => { return Promise.reject(error); }
);

// 暴露所写的内容
export { request, baseURL };