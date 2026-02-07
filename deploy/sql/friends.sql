CREATE TABLE `friends` (
    `id` varchar(24) NOT NULL,
    `user_id` varchar(24) NOT NULL COMMENT '当前用户ID（关联users.id）',
    `friend_uid` varchar(24) NOT NULL COMMENT '好友ID（关联users.id）',
    `remark` varchar(24) NOT NULL DEFAULT '' COMMENT '好友备注',
    `add_source` TINYINT DEFAULT 1 COMMENT '添加来源：1-搜索 2-群聊 3-推荐',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_user_friend` (`user_id`,`friend_uid`) COMMENT '唯一索引：避免重复添加好友',
    KEY `idx_friend_uid` (`friend_uid`) COMMENT '索引：查询某用户的好友'

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友关系表';