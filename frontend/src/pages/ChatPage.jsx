import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { sendMessage, saveReport } from '../api'
import { saveConversation, loadConversation, clearConversation, exportAsTxt } from '../utils/storage'
import './ChatPage.css'

function ChatPage({ selectedProject }) {
  const navigate = useNavigate()
  const [conversation, setConversation] = useState({
    messages: [],
    stage: 'initial',
    businessGoal: '',
    recommendations: null,
  })
  const [inputValue, setInputValue] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [showContinuePrompt, setShowContinuePrompt] = useState(false)
  const messagesEndRef = useRef(null)

  useEffect(() => {
    // Load last conversation
    const savedConversation = loadConversation()
    if (savedConversation && savedConversation.messages.length > 0) {
      setConversation(savedConversation)
      setShowContinuePrompt(true)
    } else {
      // Start new conversation
      startNewConversation()
    }
  }, [])

  useEffect(() => {
    // Auto scroll to bottom
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [conversation.messages])

  useEffect(() => {
    // Save conversation to localStorage
    if (conversation.messages.length > 0) {
      saveConversation(conversation)
    }
  }, [conversation])

  const startNewConversation = () => {
    const initialMessage = {
      role: 'assistant',
      content: '您好！我是您的一人公司商业顾问。请告诉我您的商业目标是什么？例如：跨境电商、SaaS产品、内容创作等。',
      timestamp: new Date().toISOString(),
    }
    
    setConversation({
      messages: [initialMessage],
      stage: 'initial',
      businessGoal: '',
      recommendations: null,
    })
    setShowContinuePrompt(false)
    clearConversation()
  }

  const handleSendMessage = async () => {
    if (!inputValue.trim() || loading) return
    if (!selectedProject) {
      setError('请先选择项目')
      return
    }

    const userMessage = {
      role: 'user',
      content: inputValue.trim(),
      timestamp: new Date().toISOString(),
    }

    const newMessages = [...conversation.messages, userMessage]
    setConversation(prev => ({
      ...prev,
      messages: newMessages,
    }))
    setInputValue('')
    setError(null)
    setLoading(true)

    try {
      // Only send the last 6 messages (3 rounds) to reduce API response time
      // This keeps the context relevant while staying under API Gateway's 29s timeout
      const messagesToSend = newMessages.slice(-6)
      const response = await sendMessage(messagesToSend, selectedProject.project_id)
      
      if (!response || !response.success) {
        throw new Error(response?.error || '发送消息失败')
      }
      
      const aiResponse = response.data
        
        // Parse AI response
        let assistantMessage
        let newStage = conversation.stage
        let newRecommendations = conversation.recommendations

        if (aiResponse.stage === 'recommending' && aiResponse.recommendations) {
          // AI provided recommendations - format the content to include recommendations
          let content = aiResponse.message || '根据您的情况，我为您制定了以下方案：'
          
          const recs = aiResponse.recommendations
          console.log('AI Recommendations:', recs)
          
          // Add summary
          if (recs.summary) {
            content += '\n\n📋 方案概述：\n' + recs.summary
          }
          
          // Add AI workflows
          if (recs.ai_workflows && recs.ai_workflows.length > 0) {
            content += '\n\n🤖 AI自动化工作流：'
            recs.ai_workflows.forEach((workflow, i) => {
              content += `\n\n${i + 1}. ${workflow.name}`
              if (workflow.description) content += `\n   描述：${workflow.description}`
              if (workflow.input_requirements) content += `\n   输入要求：${workflow.input_requirements}`
              if (workflow.output_requirements) content += `\n   输出要求：${workflow.output_requirements}`
              if (workflow.estimated_cost) content += `\n   预算：${workflow.estimated_cost} XZT/月`
              if (workflow.priority) content += `\n   优先级：${workflow.priority}`
            })
          }
          
          // Add human roles
          if (recs.human_roles && recs.human_roles.length > 0) {
            content += '\n\n👥 人力资源配置：'
            recs.human_roles.forEach((role, i) => {
              content += `\n\n${i + 1}. ${role.title}`
              if (role.responsibilities && role.responsibilities.length > 0) {
                content += `\n   职责：${role.responsibilities.join('、')}`
              }
              if (role.requirements && role.requirements.length > 0) {
                content += `\n   要求：${role.requirements.join('、')}`
              }
              if (role.work_hours) content += `\n   工作时间：${role.work_hours}`
              if (role.monthly_budget) content += `\n   预算：${role.monthly_budget} XZT/月`
              if (role.priority) content += `\n   优先级：${role.priority}`
            })
          }
          
          // Add phases
          if (recs.phases && recs.phases.length > 0) {
            content += '\n\n📅 实施阶段：'
            recs.phases.forEach((phase, i) => {
              content += `\n\n${i + 1}. ${phase.phase_name}`
              if (phase.duration) content += `\n   时长：${phase.duration}`
              if (phase.monthly_budget) content += `\n   月预算：${phase.monthly_budget} XZT`
              if (phase.budget_breakdown) {
                content += `\n   预算明细：`
                // Check if budget_breakdown is an object
                if (typeof phase.budget_breakdown === 'object' && !Array.isArray(phase.budget_breakdown)) {
                  Object.entries(phase.budget_breakdown).forEach(([key, value]) => {
                    content += `\n     - ${key}: ${value} XZT`
                  })
                } else if (typeof phase.budget_breakdown === 'string') {
                  // If it's a string, just display it
                  content += `\n     ${phase.budget_breakdown}`
                }
              }
            })
          }
          
          assistantMessage = {
            role: 'assistant',
            content,
            timestamp: new Date().toISOString(),
          }
          newStage = 'recommending'
          newRecommendations = aiResponse.recommendations
        } else if (aiResponse.stage === 'questioning') {
          // AI asking questions
          let content = aiResponse.message || ''
          if (aiResponse.questions && aiResponse.questions.length > 0) {
            content += '\n\n' + aiResponse.questions.map((q, i) => `${i + 1}. ${q}`).join('\n')
          }
          assistantMessage = {
            role: 'assistant',
            content,
            timestamp: new Date().toISOString(),
          }
          newStage = 'questioning'
        } else {
          // Default response
          assistantMessage = {
            role: 'assistant',
            content: aiResponse.message || aiResponse.content || '请继续...',
            timestamp: new Date().toISOString(),
          }
        }

        setConversation(prev => ({
          ...prev,
          messages: [...newMessages, assistantMessage],
          stage: newStage,
          recommendations: newRecommendations,
        }))
    } catch (err) {
      setError(err.error || err.message || '发送消息失败，请重试')
      console.error('Send message error:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSaveReport = async () => {
    if (!conversation.recommendations) {
      setError('还没有生成推荐方案')
      return
    }
    if (!selectedProject) {
      setError('请先选择项目')
      return
    }

    setLoading(true)
    setError(null)

    try {
      const response = await saveReport({
        project_id: selectedProject.project_id,
        business_goal: conversation.recommendations.business_goal || conversation.businessGoal,
        recommendations: conversation.recommendations,
      })

      if (response.success) {
        alert('报告已保存！')
        navigate('/reports')
      }
    } catch (err) {
      setError(err.error || '保存报告失败')
      console.error('Save report error:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleExportTxt = () => {
    exportAsTxt(conversation, conversation.recommendations)
  }

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  return (
    <div className="chat-page">
      {showContinuePrompt && (
        <div className="continue-prompt">
          <p>检测到上次未完成的对话</p>
          <div className="continue-actions">
            <button 
              className="btn btn-secondary"
              onClick={() => setShowContinuePrompt(false)}
            >
              继续上次对话
            </button>
            <button 
              className="btn btn-primary"
              onClick={startNewConversation}
            >
              开启新对话
            </button>
          </div>
        </div>
      )}

      <div className="chat-container">
        <div className="messages-container">
          {conversation.messages.map((message, index) => (
            <div 
              key={index} 
              className={`message ${message.role === 'user' ? 'message-user' : 'message-assistant'}`}
            >
              <div className="message-avatar">
                {message.role === 'user' ? '👤' : '👔'}
              </div>
              <div className="message-content">
                <div className="message-text">{message.content}</div>
                <div className="message-time">
                  {new Date(message.timestamp).toLocaleTimeString('zh-CN')}
                </div>
              </div>
            </div>
          ))}
          
          {loading && (
            <div className="message message-assistant">
              <div className="message-avatar">👔</div>
              <div className="message-content">
                <div className="loading-container">
                  <div className="spinner"></div>
                  <p className="loading-text">思考中...</p>
                </div>
              </div>
            </div>
          )}
          
          <div ref={messagesEndRef} />
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {conversation.stage === 'recommending' && conversation.recommendations && (
          <div className="recommendations-actions">
            <button 
              className="btn btn-primary"
              onClick={handleSaveReport}
              disabled={loading}
            >
              💾 保存报告
            </button>
            <button 
              className="btn btn-secondary"
              onClick={startNewConversation}
            >
              🔄 开启新对话
            </button>
          </div>
        )}

        <div className="input-container">
          <textarea
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="输入您的回答..."
            className="message-input"
            rows="3"
            disabled={loading}
          />
          <button 
            className="btn btn-primary btn-send"
            onClick={handleSendMessage}
            disabled={loading || !inputValue.trim()}
          >
            {loading ? '发送中...' : '发送'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default ChatPage
