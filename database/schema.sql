-- Business Consultant Database Schema

-- 商业咨询报告表
CREATE TABLE IF NOT EXISTS business_reports (
  report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_did VARCHAR(255) NOT NULL,
  project_id UUID NOT NULL,
  business_goal TEXT NOT NULL,
  recommendations JSONB NOT NULL,  -- 包含 ai_workflows, human_roles, phases
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_business_reports_user ON business_reports(user_did);
CREATE INDEX IF NOT EXISTS idx_business_reports_project ON business_reports(project_id);
CREATE INDEX IF NOT EXISTS idx_business_reports_created ON business_reports(created_at DESC);

COMMENT ON TABLE business_reports IS '商业咨询报告，存储AI生成的推荐内容';
COMMENT ON COLUMN business_reports.recommendations IS 'JSON格式：{ai_workflows: [], human_roles: [], phases: []}';
