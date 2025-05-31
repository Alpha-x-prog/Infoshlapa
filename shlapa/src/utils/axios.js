import axios from 'axios';

const instance = axios.create({
    baseURL: 'http://localhost:8080'
});

// Add a request interceptor
instance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Add a response interceptor
instance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            // If we receive a 401 response, clear the token and redirect to login
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/profile';
        }
        return Promise.reject(error);
    }
);

export default instance; 