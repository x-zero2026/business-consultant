// LocalStorage keys
const CONVERSATION_KEY = 'business_consultant_conversation'
const LAST_PROJECT_KEY = 'business_consultant_last_project'

// Save conversation to localStorage
export const saveConversation = (conversation) => {
  try {
    localStorage.setItem(CONVERSATION_KEY, JSON.stringify(conversation))
  } catch (e) {
    console.error('Failed to save conversation:', e)
  }
}

// Load conversation from localStorage
export const loadConversation = () => {
  try {
    const data = localStorage.getItem(CONVERSATION_KEY)
    return data ? JSON.parse(data) : null
  } catch (e) {
    console.error('Failed to load conversation:', e)
    return null
  }
}

// Clear conversation
export const clearConversation = () => {
  localStorage.removeItem(CONVERSATION_KEY)
}

// Save last selected project
export const saveLastProject = (projectId) => {
  localStorage.setItem(LAST_PROJECT_KEY, projectId)
}

// Load last selected project
export const loadLastProject = () => {
  return localStorage.getItem(LAST_PROJECT_KEY)
}

// Export conversation as TXT
export const exportAsTxt = (conversation, recommendations) => {
  let content = '# 一人公司商业咨询报告\n\n'
  content += `生成时间: ${new Date().toLocaleString('zh-CN')}\n\n`
  
  // Add conversation history
  content += '## 对话记录\n\n'
  conversation.messages.forEach((msg) => {
    if (msg.role === 'user') {
      content += `**用户**: ${msg.content}\n\n`
    } else if (msg.role === 'assistant') {
      content += `**顾问**: ${msg.content}\n\n`
    }
  })
  
  // Add recommendations
  if (recommendations) {
    content += '\n## 推荐方案\n\n'
    content += `**商业目标**: ${recommendations.business_goal}\n\n`
    
    if (recommendations.summary) {
      content += `**方案概述**: ${recommendations.summary}\n\n`
    }
    
    // AI Workflows
    if (recommendations.ai_workflows && recommendations.ai_workflows.length > 0) {
      content += '### AI工作流推荐\n\n'
      recommendations.ai_workflows.forEach((wf, index) => {
        content += `${index + 1}. **${wf.name}**\n`
        content += `   - 描述: ${wf.description}\n`
        content += `   - 预估成本: ${wf.estimated_cost} XZT/月\n`
        content += `   - 优先级: ${wf.priority}\n`
        content += `   - 复杂度: ${wf.complexity}\n\n`
      })
    }
    
    // Human Roles
    if (recommendations.human_roles && recommendations.human_roles.length > 0) {
      content += '### 真人岗位推荐\n\n'
      recommendations.human_roles.forEach((role, index) => {
        content += `${index + 1}. **${role.title}**\n`
        content += `   - 职责: ${role.responsibilities.join(', ')}\n`
        content += `   - 要求: ${role.requirements.join(', ')}\n`
        content += `   - 工作时间: ${role.work_hours}\n`
        content += `   - 月度预算: ${role.monthly_budget} XZT\n\n`
      })
    }
    
    // Phases
    if (recommendations.phases && recommendations.phases.length > 0) {
      content += '### 分阶段预算\n\n'
      recommendations.phases.forEach((phase, index) => {
        content += `${index + 1}. **${phase.phase_name}**\n`
        content += `   - 时长: ${phase.duration}\n`
        content += `   - 月度预算: ${phase.monthly_budget} XZT\n`
        content += `   - 预算明细:\n`
        Object.entries(phase.breakdown).forEach(([key, value]) => {
          content += `     * ${key}: ${value} XZT\n`
        })
        content += `   - 里程碑: ${phase.milestones.join(', ')}\n\n`
      })
    }
  }
  
  content += '\n---\n'
  content += '当前汇率: 1 XZT ≈ 1 CNY\n'
  content += '以上建议仅供参考，实际执行时请根据市场变化和个人情况调整。\n'
  
  // Create download
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `商业咨询报告_${new Date().getTime()}.txt`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
