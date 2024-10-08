CREATE DATABASE IF NOT EXISTS `annotation` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

USE `annotation`;

DROP TABLE IF EXISTS `task`;
CREATE TABLE IF NOT EXISTS `task` (
    `id` INT(0) UNSIGNED NOT NULL AUTO_INCREMENT,
    `url` VARCHAR(255) NOT NULL,
    `status` INT(0) NOT NULL,
    `result` TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_task_status` (`status`)
);

