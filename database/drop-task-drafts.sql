-- 删除task_drafts表及相关对象
-- 由于改为直接发布任务到任务中心，不再需要草稿功能

-- 删除函数
DROP FUNCTION IF EXISTS cleanup_expired_drafts();

-- 删除索引
DROP INDEX IF EXISTS idx_task_drafts_expires;
DROP INDEX IF EXISTS idx_task_drafts_user;

-- 删除表
DROP TABLE IF EXISTS task_drafts;

-- 验证
SELECT 'task_drafts table dropped successfully' AS status;
