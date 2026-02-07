CREATE TABLE `groups` (
    `id` varchar(24) NOT NULL,
    `name` VARCHAR(24) NOT NULL COMMENT '群名称',
    `icon` VARCHAR(191) NOT NULL DEFAULT '' COMMENT '群头像URL',
    `status` TINYINT DEFAULT 1 COMMENT '群状态：1-正常 2-解散 3-禁言',
    `creator_uid` varchar(24) NOT NULL COMMENT '群创建者ID（关联users.id）',
    `group_type` TINYINT DEFAULT 1 COMMENT '群类型：1-普通群 2-企业群 3-兴趣群',
    `is_verify` TINYINT DEFAULT 1 COMMENT '加群是否验证：0-否 1-是',
    `notification` VARCHAR(191) DEFAULT '' COMMENT '群公告',
    `notification_uid` varchar(24) DEFAULT NULL COMMENT '发布公告的用户ID',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`),
    KEY `idx_creator_uid` (`creator_uid`) COMMENT '索引：查询某用户创建的群'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群基础信息表';