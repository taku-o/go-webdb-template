-- Create "dm_posts_016" table
CREATE TABLE `dm_posts_016` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_016_created_at` (`created_at`),
  INDEX `idx_dm_posts_016_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_017" table
CREATE TABLE `dm_posts_017` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_017_created_at` (`created_at`),
  INDEX `idx_dm_posts_017_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_018" table
CREATE TABLE `dm_posts_018` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_018_created_at` (`created_at`),
  INDEX `idx_dm_posts_018_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_019" table
CREATE TABLE `dm_posts_019` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_019_created_at` (`created_at`),
  INDEX `idx_dm_posts_019_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_020" table
CREATE TABLE `dm_posts_020` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_020_created_at` (`created_at`),
  INDEX `idx_dm_posts_020_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_021" table
CREATE TABLE `dm_posts_021` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_021_created_at` (`created_at`),
  INDEX `idx_dm_posts_021_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_022" table
CREATE TABLE `dm_posts_022` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_022_created_at` (`created_at`),
  INDEX `idx_dm_posts_022_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_023" table
CREATE TABLE `dm_posts_023` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_023_created_at` (`created_at`),
  INDEX `idx_dm_posts_023_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_016" table
CREATE TABLE `dm_users_016` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_016_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_017" table
CREATE TABLE `dm_users_017` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_017_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_018" table
CREATE TABLE `dm_users_018` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_018_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_019" table
CREATE TABLE `dm_users_019` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_019_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_020" table
CREATE TABLE `dm_users_020` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_020_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_021" table
CREATE TABLE `dm_users_021` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_021_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_022" table
CREATE TABLE `dm_users_022` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_022_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_023" table
CREATE TABLE `dm_users_023` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_023_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
