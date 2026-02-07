CREATE TABLE `group_members` (
    `id` varchar(24) NOT NULL,
    `group_id` varchar(24) NOT NULL COMMENT '群ID（关联groups.id）',
    `user_id` varchar(24) NOT NULL COMMENT '群成员ID（关联users.id）',
    `role_level` TINYINT DEFAULT 3 COMMENT '角色等级：1-群主 2-管理员 3-普通成员',
    `join_source` TINYINT DEFAULT 1 COMMENT '加入来源：1-邀请 2-主动申请 3-扫码',
    `inviter_uid` varchar(24) DEFAULT NULL COMMENT '邀请人ID（如有）',
    `operator_uid` varchar(24) DEFAULT NULL COMMENT '操作人ID（拉人/踢人的管理员）',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_group_user` (`group_id`,`user_id`) COMMENT '唯一索引：避免重复加群',
    KEY `idx_group_id` (`group_id`) COMMENT '索引：查询某群的所有成员',
    KEY `idx_user_id` (`user_id`) COMMENT '索引：查询某用户加入的所有群'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群成员表';