-- 添加手机端背景图片配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('blog.background_image_mobile', '', 'blog', TRUE)
ON CONFLICT (key) DO NOTHING;
