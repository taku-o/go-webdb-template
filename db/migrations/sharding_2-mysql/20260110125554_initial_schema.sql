-- Create "dm_posts_008" table
CREATE TABLE `dm_posts_008` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_008_created_at` (`created_at`),
  INDEX `idx_dm_posts_008_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_009" table
CREATE TABLE `dm_posts_009` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_009_created_at` (`created_at`),
  INDEX `idx_dm_posts_009_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_010" table
CREATE TABLE `dm_posts_010` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_010_created_at` (`created_at`),
  INDEX `idx_dm_posts_010_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_011" table
CREATE TABLE `dm_posts_011` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_011_created_at` (`created_at`),
  INDEX `idx_dm_posts_011_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_012" table
CREATE TABLE `dm_posts_012` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_012_created_at` (`created_at`),
  INDEX `idx_dm_posts_012_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_013" table
CREATE TABLE `dm_posts_013` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_013_created_at` (`created_at`),
  INDEX `idx_dm_posts_013_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_014" table
CREATE TABLE `dm_posts_014` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_014_created_at` (`created_at`),
  INDEX `idx_dm_posts_014_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_015" table
CREATE TABLE `dm_posts_015` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_015_created_at` (`created_at`),
  INDEX `idx_dm_posts_015_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_008" table
CREATE TABLE `dm_users_008` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_008_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_009" table
CREATE TABLE `dm_users_009` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_009_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_010" table
CREATE TABLE `dm_users_010` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_010_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_011" table
CREATE TABLE `dm_users_011` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_011_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_012" table
CREATE TABLE `dm_users_012` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_012_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_013" table
CREATE TABLE `dm_users_013` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_013_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_014" table
CREATE TABLE `dm_users_014` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_014_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_015" table
CREATE TABLE `dm_users_015` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_015_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
