-- 添加友链推荐标记字段
ALTER TABLE friends ADD COLUMN IF NOT EXISTS is_featured boolean DEFAULT false;
COMMENT ON COLUMN friends.is_featured IS '是否推荐显示在页脚';

-- 将现有忽略检查的友链自动标记为推荐（兼容旧逻辑）
UPDATE friends SET is_featured = true WHERE accessible = -1;
