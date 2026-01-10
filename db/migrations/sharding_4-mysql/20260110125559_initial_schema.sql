-- Create "dm_posts_024" table
CREATE TABLE `dm_posts_024` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_024_created_at` (`created_at`),
  INDEX `idx_dm_posts_024_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_025" table
CREATE TABLE `dm_posts_025` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_025_created_at` (`created_at`),
  INDEX `idx_dm_posts_025_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_026" table
CREATE TABLE `dm_posts_026` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_026_created_at` (`created_at`),
  INDEX `idx_dm_posts_026_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_027" table
CREATE TABLE `dm_posts_027` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_027_created_at` (`created_at`),
  INDEX `idx_dm_posts_027_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_028" table
CREATE TABLE `dm_posts_028` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_028_created_at` (`created_at`),
  INDEX `idx_dm_posts_028_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_029" table
CREATE TABLE `dm_posts_029` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_029_created_at` (`created_at`),
  INDEX `idx_dm_posts_029_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_030" table
CREATE TABLE `dm_posts_030` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_030_created_at` (`created_at`),
  INDEX `idx_dm_posts_030_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_posts_031" table
CREATE TABLE `dm_posts_031` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_dm_posts_031_created_at` (`created_at`),
  INDEX `idx_dm_posts_031_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_024" table
CREATE TABLE `dm_users_024` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_024_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_025" table
CREATE TABLE `dm_users_025` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_025_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_026" table
CREATE TABLE `dm_users_026` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_026_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_027" table
CREATE TABLE `dm_users_027` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_027_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_028" table
CREATE TABLE `dm_users_028` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_028_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_029" table
CREATE TABLE `dm_users_029` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_029_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_030" table
CREATE TABLE `dm_users_030` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_030_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "dm_users_031" table
CREATE TABLE `dm_users_031` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` varchar(191) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dm_users_031_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
