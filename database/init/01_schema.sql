-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS canteen_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE canteen_system;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码（加密）',
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    role ENUM('admin', 'chef', 'employee') NOT NULL DEFAULT 'employee' COMMENT '角色',
    avatar VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_users_role (role),
    INDEX idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 员工过敏原表
CREATE TABLE IF NOT EXISTS user_allergens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    allergen VARCHAR(50) NOT NULL COMMENT '过敏原类型',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_allergen (user_id, allergen),
    INDEX idx_user_allergens_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户过敏原表';

-- 菜品分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '分类名称',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_categories_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜品分类表';

-- 菜品库表
CREATE TABLE IF NOT EXISTS dishes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '菜品名称',
    category_id BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '价格',
    image VARCHAR(500) DEFAULT NULL COMMENT '图片URL',
    description TEXT COMMENT '描述',
    daily_limit INT NOT NULL DEFAULT 50 COMMENT '每日限量份数',
    allergens JSON DEFAULT NULL COMMENT '过敏原标签，如 ["花生", "海鲜"]',
    is_available BOOLEAN DEFAULT TRUE COMMENT '是否可用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_dishes_category (category_id),
    INDEX idx_dishes_available (is_available),
    FOREIGN KEY (category_id) REFERENCES categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜品库表';

-- 周菜单模板表
CREATE TABLE IF NOT EXISTS weekly_menu_templates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    description TEXT COMMENT '模板描述',
    is_active BOOLEAN DEFAULT FALSE COMMENT '是否为当前活跃模板',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_weekly_templates_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='周菜单模板表';

-- 周菜单模板详情表
CREATE TABLE IF NOT EXISTS weekly_menu_template_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    template_id BIGINT UNSIGNED NOT NULL COMMENT '模板ID',
    day_of_week TINYINT UNSIGNED NOT NULL COMMENT '星期几：1-7 代表周一到周日',
    meal_period ENUM('breakfast', 'lunch', 'dinner') NOT NULL COMMENT '时段：早餐/午餐/晚餐',
    dish_id BIGINT UNSIGNED NOT NULL COMMENT '菜品ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_template_day_period_dish (template_id, day_of_week, meal_period, dish_id),
    INDEX idx_template_items_template (template_id),
    INDEX idx_template_items_day_period (day_of_week, meal_period),
    FOREIGN KEY (template_id) REFERENCES weekly_menu_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='周菜单模板详情表';

-- 每日菜单表
CREATE TABLE IF NOT EXISTS daily_menus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    menu_date DATE NOT NULL COMMENT '菜单日期',
    meal_period ENUM('breakfast', 'lunch', 'dinner') NOT NULL COMMENT '时段',
    start_time TIME NOT NULL COMMENT '供应开始时间',
    end_time TIME NOT NULL COMMENT '供应结束时间',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_date_period (menu_date, meal_period),
    INDEX idx_daily_menus_date (menu_date),
    INDEX idx_daily_menus_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日菜单表';

-- 每日菜单菜品关联表
CREATE TABLE IF NOT EXISTS daily_menu_dishes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    daily_menu_id BIGINT UNSIGNED NOT NULL COMMENT '每日菜单ID',
    dish_id BIGINT UNSIGNED NOT NULL COMMENT '菜品ID',
    remaining_quantity INT NOT NULL COMMENT '剩余可订份数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_menu_dish (daily_menu_id, dish_id),
    INDEX idx_daily_menu_dishes_menu (daily_menu_id),
    INDEX idx_daily_menu_dishes_dish (dish_id),
    FOREIGN KEY (daily_menu_id) REFERENCES daily_menus(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日菜单菜品关联表';

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL UNIQUE COMMENT '订单号',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '订单总金额',
    status ENUM('pending_confirm', 'in_production', 'ready_for_pickup', 'picked_up', 'reviewed') NOT NULL DEFAULT 'pending_confirm' COMMENT '订单状态',
    pickup_time_start DATETIME NOT NULL COMMENT '取餐开始时间',
    pickup_time_end DATETIME NOT NULL COMMENT '取餐结束时间',
    pickup_code VARCHAR(6) DEFAULT NULL COMMENT '取餐码（备用）',
    qr_code_content TEXT DEFAULT NULL COMMENT '二维码内容（订单号）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_orders_user (user_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_created_at (created_at),
    INDEX idx_orders_pickup_time (pickup_time_start),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 订单明细表
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    dish_id BIGINT UNSIGNED NOT NULL COMMENT '菜品ID',
    dish_name VARCHAR(100) NOT NULL COMMENT '菜品名称（冗余）',
    dish_price DECIMAL(10, 2) NOT NULL COMMENT '菜品单价（冗余）',
    quantity INT NOT NULL COMMENT '数量',
    subtotal DECIMAL(10, 2) NOT NULL COMMENT '小计',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order_items_order (order_id),
    INDEX idx_order_items_dish (dish_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';

-- 评价表
CREATE TABLE IF NOT EXISTS reviews (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    dish_id BIGINT UNSIGNED NOT NULL COMMENT '菜品ID',
    rating TINYINT UNSIGNED NOT NULL COMMENT '评分：1-5',
    comment TEXT COMMENT '评价内容',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_order_dish_user (order_id, dish_id, user_id),
    INDEX idx_reviews_dish (dish_id),
    INDEX idx_reviews_rating (rating),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评价表';
