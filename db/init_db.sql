SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `pan`;
USE `pan`;

-- ----------------------------
--  Table structure for `tbl_file`
-- ----------------------------
DROP TABLE IF EXISTS `tbl_file`;
CREATE TABLE `tbl_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_qetag` char(40) DEFAULT '' COMMENT '文件qetag',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` int(11) DEFAULT NULL COMMENT '状态(1可用/2禁用/3已删除)',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_qetag`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `tbl_user`
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `email` varchar(64) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(128) DEFAULT '' COMMENT '手机号',
  `email_validated` tinyint(1) DEFAULT '0' COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) DEFAULT '0' COMMENT '手机号是否已验证',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `profile` text COMMENT '用户属性',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `tbl_user_file`
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_file`;
CREATE TABLE `tbl_user_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL,
  `file_qetag` char(40) DEFAULT '' COMMENT '文件qetag',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `is_dir` int(11) DEFAULT '0' COMMENT '1是目录0是文件',
  `parent_dir` int(11) DEFAULT '0' COMMENT '父目录',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `status` int(11) DEFAULT '1' COMMENT '状态(1可用/2禁用/3已删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_file_dir` (`user_name`,`file_qetag`,`parent_dir`),
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `tbl_user_share_file`
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_share_file`;
CREATE TABLE `tbl_user_share_file` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `share_id` char(20) DEFAULT '' COMMENT '文件share id(主键的62进制)',
  `user_file_id` int(11) NOT NULL COMMENT '用户文件主键',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '分享时间',
  `share_pwd` varchar(50) NOT NULL DEFAULT '' COMMENT '文件分享密码',
  `share_time` int(11) DEFAULT '7' COMMENT '文件分享时长(1-1天/7-7天/0-永久)',
  PRIMARY KEY (`id`),
  KEY `share_id` (`share_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000000000 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;




USE mysql;
create user 'gocloud'@'%' identified by 'gocloud';
grant all privileges on pan.* to gocloud@'%' identified by 'gocloud';
FLUSH PRIVILEGES;