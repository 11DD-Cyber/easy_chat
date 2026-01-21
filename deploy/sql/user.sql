CREATE TABLE `users`(
    `id` varchar(24) NOT NULL,
    `avatar` VARCHAR(191) NOT NULL DEFAULT '',
    `nickname` VARCHAR(24) NOT NULL, 
    `phone` varchar(20) NOT NULL ,
    `password` VARCHAR(191) DEFAULT NULL,
    `status` TINYINT DEFAULT NULL,
    `sex` TINYINT DEFAULT NULL,
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;