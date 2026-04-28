USE canteen_system;

-- 插入默认菜品分类
INSERT INTO categories (name, sort_order) VALUES 
('主食', 1),
('热菜', 2),
('凉菜', 3),
('汤品', 4),
('饮品', 5);

-- 插入示例用户（密码使用 bcrypt 加密，密码均为：123456）
-- 管理员
INSERT INTO users (username, password, name, role) VALUES 
('admin', '$2b$10$gHMgMATUXBXg4enJ1cjUy./tn0tn5rbqg5tU8ozsRe3LXdilIOAqC', '张管理员', 'admin');

-- 厨师
INSERT INTO users (username, password, name, role) VALUES 
('chef', '$2b$10$gHMgMATUXBXg4enJ1cjUy./tn0tn5rbqg5tU8ozsRe3LXdilIOAqC', '李大厨', 'chef'),
('chef2', '$2b$10$gHMgMATUXBXg4enJ1cjUy./tn0tn5rbqg5tU8ozsRe3LXdilIOAqC', '王大厨', 'chef');

-- 员工
INSERT INTO users (username, password, name, role) VALUES 
('employee1', '$2b$10$gHMgMATUXBXg4enJ1cjUy./tn0tn5rbqg5tU8ozsRe3LXdilIOAqC', '员工张三', 'employee'),
('employee2', '$2b$10$gHMgMATUXBXg4enJ1cjUy./tn0tn5rbqg5tU8ozsRe3LXdilIOAqC', '员工李四', 'employee'),
('employee3', '$2b$10$gHMgMATUXBXg4enJ1cjUy./tn0tn5rbqg5tU8ozsRe3LXdilIOAqC', '员工王五', 'employee');

-- 插入员工过敏原示例
INSERT INTO user_allergens (user_id, allergen) VALUES 
(4, '花生'),
(4, '海鲜'),
(5, '乳制品');

-- 插入示例菜品
INSERT INTO dishes (name, category_id, price, image, description, daily_limit, allergens, is_available) VALUES 
('白米饭', 1, 2.00, NULL, '精选东北大米，软糯可口', 100, '[]', TRUE),
('馒头', 1, 1.00, NULL, '手工现蒸白面馒头', 80, '["麸质"]', TRUE),
('小笼包', 1, 8.00, NULL, '鲜肉小笼包，皮薄馅大', 50, '["麸质", "肉类"]', TRUE),

('红烧肉', 2, 18.00, NULL, '精选五花肉，肥而不腻', 40, '["肉类"]', TRUE),
('宫保鸡丁', 2, 16.00, NULL, '经典川菜，麻辣鲜香', 45, '["花生", "肉类"]', TRUE),
('清蒸鲈鱼', 2, 25.00, NULL, '新鲜鲈鱼，肉质鲜嫩', 20, '["海鲜"]', TRUE),
('番茄炒蛋', 2, 12.00, NULL, '家常菜品，酸甜可口', 50, '["鸡蛋"]', TRUE),

('凉拌黄瓜', 3, 8.00, NULL, '清脆爽口，开胃小菜', 60, '[]', TRUE),
('凉拌木耳', 3, 10.00, NULL, '营养丰富，口感爽脆', 40, '[]', TRUE),
('拍黄瓜', 3, 6.00, NULL, '蒜香浓郁，开胃解腻', 60, '["大蒜"]', TRUE),

('紫菜蛋花汤', 4, 5.00, NULL, '鲜美可口，营养丰富', 80, '["鸡蛋", "海鲜"]', TRUE),
('西红柿蛋汤', 4, 4.00, NULL, '酸甜开胃，老少皆宜', 80, '["鸡蛋"]', TRUE),
('玉米排骨汤', 4, 12.00, NULL, '滋补养生，汤鲜味美', 30, '["肉类"]', TRUE),

('鲜榨橙汁', 5, 8.00, NULL, '新鲜橙子现榨，维C满满', 40, '[]', TRUE),
('酸梅汤', 5, 5.00, NULL, '自制酸梅汤，酸甜解渴', 60, '[]', TRUE),
('牛奶', 5, 4.00, NULL, '新鲜纯牛奶', 50, '["乳制品"]', TRUE);

-- 创建一个默认周菜单模板
INSERT INTO weekly_menu_templates (name, description, is_active) VALUES 
('标准周菜单', '公司食堂标准周菜单配置', TRUE);
