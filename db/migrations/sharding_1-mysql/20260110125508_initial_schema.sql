-- Create "dm_posts_000" table
CREATE TABLE `dm_posts_000` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_000_created_at` (`created_at`),
  INDEX `idx_dm_posts_000_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_001" table
CREATE TABLE `dm_posts_001` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_001_created_at` (`created_at`),
  INDEX `idx_dm_posts_001_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_002" table
CREATE TABLE `dm_posts_002` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_002_created_at` (`created_at`),
  INDEX `idx_dm_posts_002_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_003" table
CREATE TABLE `dm_posts_003` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_003_created_at` (`created_at`),
  INDEX `idx_dm_posts_003_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_004" table
CREATE TABLE `dm_posts_004` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_004_created_at` (`created_at`),
  INDEX `idx_dm_posts_004_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_005" table
CREATE TABLE `dm_posts_005` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_005_created_at` (`created_at`),
  INDEX `idx_dm_posts_005_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_006" table
CREATE TABLE `dm_posts_006` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_006_created_at` (`created_at`),
  INDEX `idx_dm_posts_006_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_007" table
CREATE TABLE `dm_posts_007` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_007_created_at` (`created_at`),
  INDEX `idx_dm_posts_007_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_000" table
CREATE TABLE `dm_users_000` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_000_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_001" table
CREATE TABLE `dm_users_001` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_001_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_002" table
CREATE TABLE `dm_users_002` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_002_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_003" table
CREATE TABLE `dm_users_003` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_003_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_004" table
CREATE TABLE `dm_users_004` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_004_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_005" table
CREATE TABLE `dm_users_005` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_005_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_006" table
CREATE TABLE `dm_users_006` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_006_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_007" table
CREATE TABLE `dm_users_007` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_007_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
