import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
const CHAT_API_URL = import.meta.env.VITE_CHAT_API_URL
const DID_LOGIN_API_URL = import.meta.env.VITE_DID_LOGIN_API_URL

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Create separate axios instance for chat (uses Function URL with longer timeout)
const chatApi = axios.create({
  baseURL: CHAT_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 120000, // 120 second timeout
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Add auth token to chat requests
chatApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle errors
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Unauthorized - redirect to login
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      window.location.href = '/'
    }
    return Promise.reject(error.response?.data || error)
  }
)

// Handle errors for chat API
chatApi.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Unauthorized - redirect to login
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      window.location.href = '/'
    }
    return Promise.reject(error.response?.data || error)
  }
)

// Chat API (uses Function URL for longer timeout support)
export const sendMessage = (messages, projectId, stream = false) => {
  return chatApi.post('/', { messages, project_id: projectId, stream })
}

// Reports API
export const saveReport = (data) => {
  return api.post('/save-report', data)
}

export const getReports = (projectId) => {
  return api.get('/reports', { params: { project_id: projectId } })
}

export const getReport = (reportId) => {
  return api.get(`/report/${reportId}`)
}

export const deleteReport = (reportId) => {
  return api.delete(`/report/${reportId}`)
}

export const updateReportItem = (reportId, itemId, data) => {
  return api.patch(`/report/${reportId}/item/${itemId}`, data)
}

// DID Login API
export const getUserProfile = () => {
  const token = localStorage.getItem('token')
  return axios.get(`${DID_LOGIN_API_URL}/api/user/profile`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
}

// Projects API (from DID Login)
export const listProjects = () => {
  const token = localStorage.getItem('token')
  return axios.get(`${DID_LOGIN_API_URL}/api/projects`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
}

export default api
