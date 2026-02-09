import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axios from 'axios'
import { getReport, updateReportItem, identifyProfessionTags } from '../api'
import './ReportDetailPage.css'

const TASK_UI_URL = import.meta.env.VITE_TASK_UI_URL
const TASK_API_URL = import.meta.env.VITE_TASK_API_URL

function ReportDetailPage({ selectedProject }) {
  const { reportId } = useParams()
  const navigate = useNavigate()
  const [report, setReport] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [activeTab, setActiveTab] = useState('workflows')
  const [publishingWorkflow, setPublishingWorkflow] = useState(null)
  const [publishingRole, setPublishingRole] = useState(null)

  useEffect(() => {
    window.scrollTo(0, 0)
    loadReport()
  }, [reportId])

  const loadReport = async () => {
    setLoading(true)
    setError(null)

    try {
      console.log('Loading report:', reportId)
      const response = await getReport(reportId)
      console.log('Report response:', response)
      if (response.success) {
        setReport(response.data)
      } else {
        setError(response.error || '加载报告失败')
      }
    } catch (err) {
      console.error('Load report error:', err)
      setError(err.error || err.message || '加载报告失败')
    } finally {
      setLoading(false)
    }
  }

  const handlePublishWorkflow = async (workflow, index) => {
    if (!confirm(`确定要发布任务"${workflow.name}"吗？`)) return

    setPublishingWorkflow(index)
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        alert('请先登录')
        setPublishingWorkflow(null)
        return
      }

      // Step 1: Identify profession tags from task description
      const fullDescription = `${workflow.name}\n${workflow.description}\n输入要求：${workflow.input_requirements}\n输出要求：${workflow.output_requirements}`
      
      let professionTags = []
      try {
        console.log('Identifying profession tags for workflow...')
        const tagsResponse = await identifyProfessionTags(fullDescription)
        if (tagsResponse.success && tagsResponse.data?.profession_tags) {
          professionTags = tagsResponse.data.profession_tags
          console.log('Identified tags:', professionTags)
        }
      } catch (err) {
        console.error('Failed to identify profession tags:', err)
        // Continue without tags if identification fails
      }

      // Step 2: Create task with profession tags
      const taskData = {
        project_id: selectedProject?.project_id,
        task_name: `开发AI工作流：${workflow.name}`,
        task_description: `${workflow.description}\n\n**输入要求**：\n${workflow.input_requirements}\n\n**输出要求**：\n${workflow.output_requirements}`,
        acceptance_criteria: workflow.acceptance_criteria || '1. API调用成功\n2. 输出格式正确\n3. 性能达标',
        reward_amount: String(workflow.estimated_cost || 0),
        visibility: 'global',
        profession_tags: professionTags,
      }

      const response = await axios.post(`${TASK_API_URL}/tasks`, taskData, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      })

      if (response.data.success) {
        // Update item status
        await updateReportItem(reportId, `wf-${index}`, {
          status: 'published',
          task_id: response.data.data.task_id,
        })

        alert(`任务发布成功！任务ID: ${response.data.data.task_id}`)
        
        // Open task in task center
        window.open(`${TASK_UI_URL}/tasks/${response.data.data.task_id}`, '_blank')

        // Reload report to show updated status
        setTimeout(() => loadReport(), 1000)
      }
    } catch (err) {
      console.error('Publish workflow error:', err)
      alert(err.response?.data?.error || err.message || '发布任务失败')
    } finally {
      setPublishingWorkflow(null)
    }
  }

  const handlePublishRole = async (role, index) => {
    if (!confirm(`确定要发布招聘任务"${role.title}"吗？`)) return

    setPublishingRole(index)
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        alert('请先登录')
        setPublishingRole(null)
        return
      }

      // Step 1: Identify profession tags from role description
      const fullDescription = `${role.title}\n职责：${role.responsibilities?.join(', ') || ''}\n要求：${role.requirements?.join(', ') || ''}`
      
      let professionTags = []
      try {
        console.log('Identifying profession tags for role...')
        const tagsResponse = await identifyProfessionTags(fullDescription)
        if (tagsResponse.success && tagsResponse.data?.profession_tags) {
          professionTags = tagsResponse.data.profession_tags
          console.log('Identified tags:', professionTags)
        }
      } catch (err) {
        console.error('Failed to identify profession tags:', err)
        // Continue without tags if identification fails
      }

      // Step 2: Create task with profession tags
      const taskData = {
        project_id: selectedProject?.project_id,
        task_name: `招聘：${role.title}`,
        task_description: `**职责**：\n${role.responsibilities?.join('\n') || ''}\n\n**要求**：\n${role.requirements?.join('\n') || ''}\n\n**工作时间**：\n${role.work_hours || '待定'}`,
        acceptance_criteria: role.trial_period_criteria || '试用期1个月，考核标准：\n1. 按时完成工作\n2. 沟通顺畅\n3. 质量达标',
        reward_amount: String(role.monthly_budget || 0),
        visibility: 'global',
        profession_tags: professionTags,
      }

      const response = await axios.post(`${TASK_API_URL}/tasks`, taskData, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      })

      if (response.data.success) {
        // Update item status
        await updateReportItem(reportId, `role-${index}`, {
          status: 'published',
          task_id: response.data.data.task_id,
        })

        alert(`任务发布成功！任务ID: ${response.data.data.task_id}`)
        
        // Open task in task center
        window.open(`${TASK_UI_URL}/tasks/${response.data.data.task_id}`, '_blank')

        // Reload report to show updated status
        setTimeout(() => loadReport(), 1000)
      }
    } catch (err) {
      console.error('Publish role error:', err)
      alert(err.response?.data?.error || err.message || '发布任务失败')
    } finally {
      setPublishingRole(null)
    }
  }

  const handleCancelPublish = async (itemId) => {
    if (!confirm('确定要取消发布吗？')) return

    try {
      await updateReportItem(reportId, itemId, {
        status: null,
        task_id: null,
      })
      loadReport()
    } catch (err) {
      alert(err.error || '取消失败')
      console.error('Cancel publish error:', err)
    }
  }

  const handleExportTxt = () => {
    if (!report) return
    
    const conversation = {
      messages: [], // No conversation history in saved report
    }
    exportAsTxt(conversation, report.recommendations)
  }

  if (loading) {
    return (
      <div className="report-detail-page">
        <div className="loading-container">
          <div className="spinner"></div>
          <p>加载中...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="report-detail-page">
        <div className="error-message">{error}</div>
      </div>
    )
  }

  if (!report) {
    return (
      <div className="report-detail-page">
        <div className="error-message">报告不存在</div>
      </div>
    )
  }

  const { recommendations } = report

  return (
    <div className="report-detail-page">
      <div className="page-header">
        <h1>{report.business_goal}</h1>
      </div>

      <div className="report-meta">
        <span>创建时间: {new Date(report.created_at).toLocaleString('zh-CN')}</span>
      </div>

      {recommendations?.summary && (
        <div className="report-summary-box">
          <h3>方案概述</h3>
          <p>{recommendations.summary}</p>
        </div>
      )}

      <div className="tabs">
        <button 
          className={`tab ${activeTab === 'workflows' ? 'active' : ''}`}
          onClick={() => setActiveTab('workflows')}
        >
          AI工作流 ({recommendations?.ai_workflows?.length || 0})
        </button>
        <button 
          className={`tab ${activeTab === 'roles' ? 'active' : ''}`}
          onClick={() => setActiveTab('roles')}
        >
          真人岗位 ({recommendations?.human_roles?.length || 0})
        </button>
        <button 
          className={`tab ${activeTab === 'budget' ? 'active' : ''}`}
          onClick={() => setActiveTab('budget')}
        >
          预算规划
        </button>
      </div>

      <div className="tab-content">
        {activeTab === 'workflows' && (
          <div className="workflows-list">
            {recommendations?.ai_workflows?.map((workflow, index) => (
              <div key={index} className="item-card">
                <div className="item-header">
                  <h3>{workflow.name}</h3>
                  <div className="item-badges">
                    <span className={`badge badge-${workflow.priority}`}>
                      {workflow.priority}
                    </span>
                    <span className="badge badge-secondary">
                      {workflow.complexity}
                    </span>
                  </div>
                </div>
                <p className="item-description">{workflow.description}</p>
                
                <div className="item-details">
                  <div className="detail-row">
                    <strong>输入要求：</strong>
                    <span>{workflow.input_requirements}</span>
                  </div>
                  <div className="detail-row">
                    <strong>输出要求：</strong>
                    <span>{workflow.output_requirements}</span>
                  </div>
                  {workflow.tools_needed && (
                    <div className="detail-row">
                      <strong>所需工具：</strong>
                      <span>{workflow.tools_needed.join(', ')}</span>
                    </div>
                  )}
                  <div className="detail-row">
                    <strong>预估成本：</strong>
                    <span className="cost">{workflow.estimated_cost} XZT/月</span>
                  </div>
                </div>

                <div className="item-actions">
                  {workflow.status === 'published' ? (
                    <span className="status-badge status-published">
                      ✓ 已发布
                    </span>
                  ) : workflow.status === 'draft_created' ? (
                    <>
                      <span className="status-badge status-draft">
                        ⏳ 等待发布
                      </span>
                      <button 
                        className="btn btn-secondary btn-sm"
                        onClick={() => handleCancelPublish(`wf-${index}`)}
                      >
                        取消发布
                      </button>
                    </>
                  ) : (
                    <button 
                      className="btn btn-primary"
                      onClick={() => handlePublishWorkflow(workflow, index)}
                      disabled={publishingWorkflow === index}
                    >
                      {publishingWorkflow === index ? (
                        <>
                          <span className="spinner-small"></span>
                          发布中...
                        </>
                      ) : (
                        '发布任务'
                      )}
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'roles' && (
          <div className="roles-list">
            {recommendations?.human_roles?.map((role, index) => (
              <div key={index} className="item-card">
                <div className="item-header">
                  <h3>{role.title}</h3>
                  <span className={`badge badge-${role.priority}`}>
                    {role.priority}
                  </span>
                </div>
                
                <div className="item-details">
                  <div className="detail-row">
                    <strong>职责：</strong>
                    <ul>
                      {role.responsibilities?.map((resp, i) => (
                        <li key={i}>{resp}</li>
                      ))}
                    </ul>
                  </div>
                  <div className="detail-row">
                    <strong>要求：</strong>
                    <ul>
                      {role.requirements?.map((req, i) => (
                        <li key={i}>{req}</li>
                      ))}
                    </ul>
                  </div>
                  <div className="detail-row">
                    <strong>工作时间：</strong>
                    <span>{role.work_hours}</span>
                  </div>
                  <div className="detail-row">
                    <strong>月度预算：</strong>
                    <span className="cost">{role.monthly_budget} XZT</span>
                  </div>
                  {role.trial_period_criteria && (
                    <div className="detail-row">
                      <strong>试用期考核：</strong>
                      <pre>{role.trial_period_criteria}</pre>
                    </div>
                  )}
                </div>

                <div className="item-actions">
                  {role.status === 'published' ? (
                    <span className="status-badge status-published">
                      ✓ 已发布
                    </span>
                  ) : role.status === 'draft_created' ? (
                    <>
                      <span className="status-badge status-draft">
                        ⏳ 等待发布
                      </span>
                      <button 
                        className="btn btn-secondary btn-sm"
                        onClick={() => handleCancelPublish(`role-${index}`)}
                      >
                        取消发布
                      </button>
                    </>
                  ) : (
                    <button 
                      className="btn btn-primary"
                      onClick={() => handlePublishRole(role, index)}
                      disabled={publishingRole === index}
                    >
                      {publishingRole === index ? (
                        <>
                          <span className="spinner-small"></span>
                          发布中...
                        </>
                      ) : (
                        '发布招聘'
                      )}
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'budget' && (
          <div className="budget-list">
            {recommendations?.phases?.map((phase, index) => (
              <div key={index} className="phase-card">
                <h3>{phase.phase_name}</h3>
                <div className="phase-meta">
                  <span>时长: {phase.duration}</span>
                  <span className="phase-budget">
                    月度预算: {phase.monthly_budget} XZT
                  </span>
                </div>

                <div className="budget-breakdown">
                  <h4>预算明细</h4>
                  {(() => {
                    console.log('Phase data:', phase);
                    console.log('budget_breakdown:', phase.budget_breakdown);
                    console.log('budget_breakdown type:', typeof phase.budget_breakdown);
                    return null;
                  })()}
                  <table className="budget-table">
                    <tbody>
                      {phase.budget_breakdown && typeof phase.budget_breakdown === 'object' && !Array.isArray(phase.budget_breakdown) ? (
                        Object.entries(phase.budget_breakdown).map(([key, value]) => (
                          <tr key={key}>
                            <td>{key}</td>
                            <td className="amount">{value} XZT</td>
                          </tr>
                        ))
                      ) : (
                        <tr>
                          <td colSpan="2">暂无明细</td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>

                {phase.milestones && phase.milestones.length > 0 && (
                  <div className="milestones">
                    <h4>里程碑目标</h4>
                    <ul>
                      {phase.milestones.map((milestone, i) => (
                        <li key={i}>{milestone}</li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default ReportDetailPage
