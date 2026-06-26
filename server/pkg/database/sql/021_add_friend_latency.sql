-- 添加友链延迟检测字段
ALTER TABLE friends ADD COLUMN IF NOT EXISTS latency integer DEFAULT 0;
COMMENT ON COLUMN friends.latency IS '最后检测延迟(毫秒), 0=未检测';
