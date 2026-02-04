import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { getReports, deleteReport } from '../api'
import './ReportsPage.css'

function ReportsPage({ selectedProject }) {
  const [reports, setReports] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (selectedProject) {
      loadReports()
    }
  }, [selectedProject])

  const loadReports = async () => {
    setLoading(true)
    setError(null)

    try {
      const response = await getReports(selectedProject.project_id)
      if (response.success) {
        setReports(response.data || [])
      }
    } catch (err) {
      setError(err.error || '加载报告失败')
      console.error('Load reports error:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (reportId) => {
    if (!confirm('确定要删除这份报告吗？')) return

    try {
      const response = await deleteReport(reportId)
      if (response.success) {
        setReports(reports.filter(r => r.report_id !== reportId))
      }
    } catch (err) {
      alert(err.error || '删除失败')
      console.error('Delete report error:', err)
    }
  }

  const countPublishedTasks = (recommendations) => {
    if (!recommendations) return 0
    
    let count = 0
    if (recommendations.ai_workflows) {
      count += recommendations.ai_workflows.filter(w => w.status === 'published').length
    }
    if (recommendations.human_roles) {
      count += recommendations.human_roles.filter(r => r.status === 'published').length
    }
    return count
  }

  if (loading) {
    return (
      <div className="reports-page">
        <div className="loading-container">
          <div className="spinner"></div>
          <p>加载中...</p>
        </div>
      </div>
    )
  }

  if (!selectedProject) {
    return (
      <div className="reports-page">
        <div className="empty-state">
          <p>请先选择项目</p>
        </div>
      </div>
    )
  }

  return (
    <div className="reports-page">
      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {reports.length === 0 ? (
        <div className="empty-state">
          <div className="empty-icon">📋</div>
          <p>还没有保存的报告</p>
          <Link to="/" className="btn btn-primary">
            开始咨询
          </Link>
        </div>
      ) : (
        <div className="reports-grid">
          {reports.map((report) => (
            <div key={report.report_id} className="report-card">
              <div className="report-header">
                <h3 className="report-title">{report.business_goal}</h3>
                <div className="report-date">
                  {new Date(report.created_at).toLocaleDateString('zh-CN')}
                </div>
              </div>

              <div className="report-stats">
                <div className="stat">
                  <span className="stat-label">AI工作流</span>
                  <span className="stat-value">
                    {report.recommendations?.ai_workflows?.length || 0}
                  </span>
                </div>
                <div className="stat">
                  <span className="stat-label">真人岗位</span>
                  <span className="stat-value">
                    {report.recommendations?.human_roles?.length || 0}
                  </span>
                </div>
                <div className="stat">
                  <span className="stat-label">已发布任务</span>
                  <span className="stat-value">
                    {countPublishedTasks(report.recommendations)}
                  </span>
                </div>
              </div>

              {report.recommendations?.summary && (
                <div className="report-summary">
                  {report.recommendations.summary}
                </div>
              )}

              <div className="report-actions">
                <Link 
                  to={`/reports/${report.report_id}`}
                  className="btn btn-primary btn-sm"
                >
                  查看详情
                </Link>
                <button 
                  className="btn btn-danger btn-sm"
                  onClick={() => handleDelete(report.report_id)}
                >
                  删除
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default ReportsPage
