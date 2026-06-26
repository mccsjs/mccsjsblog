-- 更新第三方服务默认地址（项目开源用）
-- 执行时间：2026-06-26
-- 说明：这些是新默认值，已初始化过的数据库需要执行此脚本

UPDATE settings SET value = 'https://api.i-meto.com/meting/api' WHERE key = 'blog.meting_api';
UPDATE settings SET value = 'https://cravatar.cn/avatar/%s?s=200&d=robohash' WHERE key = 'blog.cravatar_url';
UPDATE settings SET value = 'http://ip-api.com/json/%s?lang=zh-CN' WHERE key = 'blog.ip_api_url';
UPDATE settings SET value = 'https://pixhub.flec.top' WHERE key = 'blog.cover_maker_api';
