CREATE DATABASE IF NOT EXISTS `manager` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `manager`;

DROP TABLE IF EXISTS `user_manager`;
CREATE TABLE `user_manager` (
    `csd_id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '用户账号自增ID',
    `username` VARCHAR(255) NOT NULL UNIQUE COMMENT '用户名',
    `student_id` VARCHAR(20) NOT NULL UNIQUE COMMENT '学号',
    `real_name` VARCHAR(100) NOT NULL COMMENT '真实姓名',
    `phone_number` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `email` VARCHAR(255) DEFAULT NULL UNIQUE COMMENT '邮箱',
    `salt` CHAR(64) NOT NULL COMMENT '密码盐',
    `hash_password` CHAR(64) NOT NULL COMMENT '加密后的密码',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '账号创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '账号更新时间',
    `invite_by` INT UNSIGNED DEFAULT NULL COMMENT '邀请人ID',
    `status` ENUM('active', 'inactive', 'frozen') DEFAULT 'active' COMMENT '账号状态'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `admin_manager`;
CREATE TABLE `admin_manager` (
    `csd_id` INT UNSIGNED NOT NULL COMMENT '用户账号自增ID',
    `role_id` INT UNSIGNED NOT NULL COMMENT '身份的自增ID',
    `invite_by` INT UNSIGNED NOT NULL COMMENT '邀请人'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `role_manager`;
CREATE TABLE `role_manager` (
    `role_id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '身份的自增ID',
    `role_name` VARCHAR(255) NOT NULL UNIQUE COMMENT '身份名',
    `authorities` TEXT NOT NULL COMMENT '身份的自增ID的集合'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `authorities_manager`;
CREATE TABLE `authorities_manager` (
    `authorities_id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '权限的ID',
    `authorities_name` VARCHAR(255) NOT NULL UNIQUE COMMENT '权限名',
    `authorities_desc` VARCHAR(255) NOT NULL UNIQUE COMMENT '权限解释'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认的 root 用户
INSERT INTO `user_manager` (`username`, `student_id`, `real_name`, `phone_number`, `email`, `salt`, `hash_password`, `status`)
VALUES
('root', '0', '超级管理员', '0', 'root@example.com', 'random_salt_value', 'hashed_password_value', 'active');

-- 插入权限数据
INSERT INTO `authorities_manager` (`authorities_name`, `authorities_desc`) VALUES
('create_user', '创建用户权限'),
('delete_user', '删除用户权限'),
('create_role', '创建角色权限'),
('delete_role', '删除角色权限'),
('assign_permissions', '分配权限权限'),
('view_reports', '查看报告权限'),
('modify_settings', '修改设置权限'),
('access_all_data', '访问所有数据权限');

-- 插入 root 角色，并赋予所有权限
INSERT INTO `role_manager` (`role_name`, `authorities`) 
VALUES ('root', '1,2,3,4,5,6,7,8'); -- 假设权限ID是1到8

-- 插入 root 用户
INSERT INTO `admin_manager` (`csd_id`, `role_id`, `invite_by`) 
VALUES (1, 1, 1);