CREATE TABLE `friend_requests` (
    `id` varchar(24) NOT NULL,
    `user_id` varchar(24) NOT NULL COMMENT '申请人ID（关联users.id）',
    `req_uid` varchar(24) NOT NULL COMMENT '被申请人ID（关联users.id）',
    `req_msg` VARCHAR(191) NOT NULL DEFAULT '' COMMENT '申请验证消息',
    `handle_result` TINYINT DEFAULT 0 COMMENT '处理结果：0-待处理 1-同意 2-拒绝',
    `handle_uid` varchar(24) DEFAULT NULL COMMENT '处理人ID（通常是req_uid）',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '处理时间',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_user_req` (`user_id`,`req_uid`) COMMENT '唯一索引：避免重复申请',
    KEY `idx_req_uid` (`req_uid`) COMMENT '索引：查询某用户收到的好友申请'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友申请表';