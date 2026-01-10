-- Create "dm_news" table
CREATE TABLE `dm_news` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `author_id` int NULL,
  `published_at` timestamp NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_news_author_id` (`author_id`),
  INDEX `idx_dm_news_published_at` (`published_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_menu" table
CREATE TABLE `goadmin_menu` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT 0,
  `type` int NOT NULL DEFAULT 0,
  `order` int NOT NULL DEFAULT 0,
  `title` text NOT NULL,
  `icon` text NOT NULL,
  `uri` varchar(255) NOT NULL DEFAULT "",
  `header` text NULL,
  `plugin_name` varchar(255) NOT NULL DEFAULT "",
  `uuid` text NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_operation_log" table
CREATE TABLE `goadmin_operation_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `path` text NOT NULL,
  `method` text NOT NULL,
  `ip` text NOT NULL,
  `input` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_goadmin_operation_log_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_permissions" table
CREATE TABLE `goadmin_permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `slug` varchar(191) NOT NULL,
  `http_method` text NULL,
  `http_path` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_goadmin_permissions_slug` (`slug`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_role_menu" table
CREATE TABLE `goadmin_role_menu` (
  `role_id` int NOT NULL,
  `menu_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`role_id`, `menu_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_role_permissions" table
CREATE TABLE `goadmin_role_permissions` (
  `role_id` int NOT NULL,
  `permission_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`role_id`, `permission_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_role_users" table
CREATE TABLE `goadmin_role_users` (
  `role_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`role_id`, `user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_roles" table
CREATE TABLE `goadmin_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `slug` varchar(191) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_goadmin_roles_slug` (`slug`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_session" table
CREATE TABLE `goadmin_session` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sid` text NOT NULL,
  `values` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_site" table
CREATE TABLE `goadmin_site` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` text NULL,
  `value` text NULL,
  `description` text NULL,
  `state` int NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_user_permissions" table
CREATE TABLE `goadmin_user_permissions` (
  `user_id` int NOT NULL,
  `permission_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `permission_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "goadmin_users" table
CREATE TABLE `goadmin_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(191) NOT NULL,
  `password` text NOT NULL,
  `name` text NOT NULL,
  `avatar` text NULL,
  `remember_token` text NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_goadmin_users_username` (`username`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
