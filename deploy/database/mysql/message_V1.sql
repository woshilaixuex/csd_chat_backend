CREATE DATABASE IF NOT EXISTS `message` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `message`;

DROP TABLE IF EXISTS `im_user`;
DROP TABLE IF EXISTS `im_msg_content`;
DROP TABLE IF EXISTS `im_msg_relation`;
DROP TABLE IF EXISTS `im_msg_contact`;


-- 创建 IM_MSG_CONTENT 表
CREATE TABLE `im_msg_content` (
    `mid` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID',
    `content` VARCHAR(1024) NOT NULL COMMENT '消息内容',
    `sender_id` INT UNSIGNED NOT NULL COMMENT '发送者ID',
    `recipient_id` INT UNSIGNED NOT NULL COMMENT '接收者ID',
    `msg_type` INT NOT NULL COMMENT '消息类型（例如文本、图片等）',
    `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '消息发送时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建 IM_MSG_RELATION 表
CREATE TABLE `im_msg_relation` (
    `owner_uid` INT UNSIGNED NOT NULL COMMENT '拥有者用户ID',
    `other_uid` INT UNSIGNED NOT NULL COMMENT '另一个用户ID',
    `mid` INT UNSIGNED NOT NULL COMMENT '消息ID',
    `type` INT NOT NULL COMMENT '关系类型',
    `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '关系建立时间',
    PRIMARY KEY (`owner_uid`, `mid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建 IM_MSG_CONTACT 表
CREATE TABLE `im_msg_contact` (
    `owner_uid` INT UNSIGNED NOT NULL COMMENT '拥有者用户ID',
    `other_uid` INT UNSIGNED NOT NULL COMMENT '联系人ID',
    `mid` INT UNSIGNED NOT NULL COMMENT '消息ID',
    `type` INT NOT NULL COMMENT '联系类型',
    `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`owner_uid`, `other_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 为 IM_MSG_RELATION 表创建索引
CREATE INDEX `idx_owneruid_otheruid_msgid` ON `im_msg_relation`(`owner_uid`, `other_uid`, `mid`);