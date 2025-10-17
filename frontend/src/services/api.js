import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (email, password) => 
    api.post('/auth/login', { email, password }),
  register: (email, password) => 
    api.post('/auth/register', { email, password }),
};

export const vocabAPI = {
  getVocabulary: (date, search = '') => 
    api.get('/vocab', { params: { date, search } }),
  createVocabulary: (data) => 
    api.post('/vocab', data),
  updateVocabulary: (id, data) => 
    api.put(`/vocab/${id}`, data),
  deleteVocabulary: (id) => 
    api.delete(`/vocab/${id}`),
};

export default api;