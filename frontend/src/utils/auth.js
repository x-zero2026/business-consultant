// Get user info from localStorage
export const getUserInfo = () => {
  const userInfoStr = localStorage.getItem('userInfo')
  if (!userInfoStr) return null
  
  try {
    return JSON.parse(userInfoStr)
  } catch (e) {
    return null
  }
}

// Save user info to localStorage
export const saveUserInfo = (userInfo) => {
  localStorage.setItem('userInfo', JSON.stringify(userInfo))
}

// Get token from localStorage
export const getToken = () => {
  return localStorage.getItem('token')
}

// Save token to localStorage
export const saveToken = (token) => {
  localStorage.setItem('token', token)
}

// Clear auth data
export const clearAuth = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
}

// Check if user is logged in
export const isLoggedIn = () => {
  return !!getToken()
}
