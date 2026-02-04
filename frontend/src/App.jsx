import { useState, useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Link, NavLink, useNavigate } from 'react-router-dom'
import { getUserInfo, saveUserInfo } from './utils/auth'
import { listProjects } from './api'
import { saveLastProject, loadLastProject } from './utils/storage'
import ChatPage from './pages/ChatPage'
import ReportsPage from './pages/ReportsPage'
import ReportDetailPage from './pages/ReportDetailPage'
import StarfieldBackground from './components/StarfieldBackground'
import './App.css'

function App() {
  const [userInfo, setUserInfo] = useState(null)
  const [projects, setProjects] = useState([])
  const [selectedProject, setSelectedProject] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    initializeApp()
  }, [])

  const initializeApp = async () => {
    // Check if token is in URL (from DID login redirect)
    const urlParams = new URLSearchParams(window.location.search)
    const tokenFromUrl = urlParams.get('token')
    
    if (tokenFromUrl) {
      console.log('✅ Token received from URL')
      localStorage.setItem('token', tokenFromUrl)
      
      try {
        const payload = JSON.parse(atob(tokenFromUrl.split('.')[1]))
        const userInfo = {
          did: payload.did,
          username: payload.username,
        }
        saveUserInfo(userInfo)
        setUserInfo(userInfo)
        
        // Remove token from URL
        window.history.replaceState({}, document.title, window.location.pathname)
      } catch (err) {
        console.error('Failed to parse token:', err)
      }
    } else {
      const user = getUserInfo()
      setUserInfo(user)
    }

    // Load projects
    if (getUserInfo()) {
      await loadProjects()
    }
    
    setLoading(false)
  }

  const loadProjects = async () => {
    try {
      const response = await listProjects()
      const projectList = response.data.data || []
      setProjects(projectList)
      
      // Load last selected project or use first one
      const lastProjectId = loadLastProject()
      const project = lastProjectId 
        ? projectList.find(p => p.project_id === lastProjectId)
        : projectList[0]
      
      if (project) {
        setSelectedProject(project)
      }
    } catch (err) {
      console.error('Load projects error:', err)
    }
  }

  const handleProjectChange = (projectId) => {
    const project = projects.find(p => p.project_id === projectId)
    setSelectedProject(project)
    saveLastProject(projectId)
  }

  if (loading) {
    return (
      <div className="app">
        <div className="loading-screen">
          <div className="loading-spinner-center">
            <div className="spinner"></div>
            <p>加载中...</p>
          </div>
        </div>
      </div>
    )
  }

  if (!userInfo) {
    return (
      <div className="app">
        <div className="login-prompt">
          <div className="login-card">
            <h1>👔 一人公司智能商业顾问</h1>
            <p>AI驱动的商业咨询服务，帮助您规划资源配置和预算</p>
            <div className="login-features">
              <div className="feature">
                <span className="feature-icon">💬</span>
                <span>智能对话咨询</span>
              </div>
              <div className="feature">
                <span className="feature-icon">📊</span>
                <span>个性化推荐</span>
              </div>
              <div className="feature">
                <span className="feature-icon">💰</span>
                <span>预算规划</span>
              </div>
            </div>
            <p className="login-hint">请先登录 X-Zero 系统以使用商业顾问服务</p>
            <a 
              href="https://main.d2fozf421c6ftf.amplifyapp.com" 
              className="btn btn-primary btn-large"
            >
              前往登录
            </a>
          </div>
        </div>
      </div>
    )
  }

  return (
    <Router>
      <div className="app">
        <StarfieldBackground />
        {/* Header */}
        <header className="header">
          <div className="header-left">
            <Link to="/" className="header-title">
              👔 商业顾问
            </Link>
            <nav className="header-nav">
              <NavLink to="/" className={({ isActive }) => isActive ? "nav-link active" : "nav-link"} end>对话</NavLink>
              <NavLink to="/reports" className={({ isActive }) => isActive ? "nav-link active" : "nav-link"}>我的报告</NavLink>
            </nav>
          </div>
          <div className="header-right">
            {projects.length > 0 && selectedProject && (
              <select
                value={selectedProject.project_id}
                onChange={(e) => handleProjectChange(e.target.value)}
                className="project-select"
              >
                {projects.map((project) => (
                  <option key={project.project_id} value={project.project_id}>
                    {project.project_name}
                  </option>
                ))}
              </select>
            )}
            <div className="user-info">
              <span className="user-name">{userInfo.username}</span>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="main-content">
          <Routes>
            <Route 
              path="/" 
              element={<ChatPage selectedProject={selectedProject} />} 
            />
            <Route 
              path="/reports" 
              element={<ReportsPage selectedProject={selectedProject} />} 
            />
            <Route 
              path="/reports/:reportId" 
              element={<ReportDetailPage selectedProject={selectedProject} />} 
            />
          </Routes>
        </main>

        {/* Footer */}
        <footer className="footer">
          <p>当前汇率: 1 XZT ≈ 1 CNY</p>
          <p>以上建议仅供参考，实际执行时请根据市场变化和个人情况调整</p>
        </footer>
      </div>
    </Router>
  )
}

export default App
