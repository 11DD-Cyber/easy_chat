CREATE TABLE `group_requests` (
    `id` varchar(24) NOT NULL,
    `group_id` varchar(24) NOT NULL COMMENT '申请加入的群ID（关联groups.id）',
    `req_id` varchar(24) NOT NULL COMMENT '申请人ID（关联users.id）',
    `req_msg` VARCHAR(191) NOT NULL DEFAULT '' COMMENT '加群验证消息',
    `join_source` TINYINT DEFAULT 1 COMMENT '申请来源：1-搜索 2-好友邀请 3-扫码',
    `inviter_uid` varchar(24) DEFAULT NULL COMMENT '邀请人ID（如有）',
    `handle_uid` varchar(24) DEFAULT NULL COMMENT '处理人ID（群主/管理员）',
    `handle_result` TINYINT DEFAULT 0 COMMENT '处理结果：0-待处理 1-同意 2-拒绝',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '处理时间',
    PRIMARY KEY(`id`),
    KEY `idx_group_id` (`group_id`) COMMENT '索引：查询某群的加群申请',
    KEY `idx_req_id` (`req_id`) COMMENT '索引：查询某用户发起的加群申请'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群申请表';