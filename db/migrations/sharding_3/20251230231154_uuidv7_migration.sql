-- UUIDv7 Migration: Change ID columns from bigint to varchar(32)
-- dm_users_016 to dm_users_023
-- dm_posts_016 to dm_posts_023

-- Delete existing data and recreate tables with new schema

-- dm_users_016
DROP TABLE IF EXISTS `dm_users_016`;
CREATE TABLE `dm_users_016` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_016_email` ON `dm_users_016` (`email`);

-- dm_users_017
DROP TABLE IF EXISTS `dm_users_017`;
CREATE TABLE `dm_users_017` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_017_email` ON `dm_users_017` (`email`);

-- dm_users_018
DROP TABLE IF EXISTS `dm_users_018`;
CREATE TABLE `dm_users_018` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_018_email` ON `dm_users_018` (`email`);

-- dm_users_019
DROP TABLE IF EXISTS `dm_users_019`;
CREATE TABLE `dm_users_019` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_019_email` ON `dm_users_019` (`email`);

-- dm_users_020
DROP TABLE IF EXISTS `dm_users_020`;
CREATE TABLE `dm_users_020` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_020_email` ON `dm_users_020` (`email`);

-- dm_users_021
DROP TABLE IF EXISTS `dm_users_021`;
CREATE TABLE `dm_users_021` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_021_email` ON `dm_users_021` (`email`);

-- dm_users_022
DROP TABLE IF EXISTS `dm_users_022`;
CREATE TABLE `dm_users_022` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_022_email` ON `dm_users_022` (`email`);

-- dm_users_023
DROP TABLE IF EXISTS `dm_users_023`;
CREATE TABLE `dm_users_023` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_023_email` ON `dm_users_023` (`email`);

-- dm_posts_016
DROP TABLE IF EXISTS `dm_posts_016`;
CREATE TABLE `dm_posts_016` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_016_user_id` ON `dm_posts_016` (`user_id`);
CREATE INDEX `idx_dm_posts_016_created_at` ON `dm_posts_016` (`created_at`);

-- dm_posts_017
DROP TABLE IF EXISTS `dm_posts_017`;
CREATE TABLE `dm_posts_017` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_017_user_id` ON `dm_posts_017` (`user_id`);
CREATE INDEX `idx_dm_posts_017_created_at` ON `dm_posts_017` (`created_at`);

-- dm_posts_018
DROP TABLE IF EXISTS `dm_posts_018`;
CREATE TABLE `dm_posts_018` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_018_user_id` ON `dm_posts_018` (`user_id`);
CREATE INDEX `idx_dm_posts_018_created_at` ON `dm_posts_018` (`created_at`);

-- dm_posts_019
DROP TABLE IF EXISTS `dm_posts_019`;
CREATE TABLE `dm_posts_019` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_019_user_id` ON `dm_posts_019` (`user_id`);
CREATE INDEX `idx_dm_posts_019_created_at` ON `dm_posts_019` (`created_at`);

-- dm_posts_020
DROP TABLE IF EXISTS `dm_posts_020`;
CREATE TABLE `dm_posts_020` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_020_user_id` ON `dm_posts_020` (`user_id`);
CREATE INDEX `idx_dm_posts_020_created_at` ON `dm_posts_020` (`created_at`);

-- dm_posts_021
DROP TABLE IF EXISTS `dm_posts_021`;
CREATE TABLE `dm_posts_021` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_021_user_id` ON `dm_posts_021` (`user_id`);
CREATE INDEX `idx_dm_posts_021_created_at` ON `dm_posts_021` (`created_at`);

-- dm_posts_022
DROP TABLE IF EXISTS `dm_posts_022`;
CREATE TABLE `dm_posts_022` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_022_user_id` ON `dm_posts_022` (`user_id`);
CREATE INDEX `idx_dm_posts_022_created_at` ON `dm_posts_022` (`created_at`);

-- dm_posts_023
DROP TABLE IF EXISTS `dm_posts_023`;
CREATE TABLE `dm_posts_023` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_023_user_id` ON `dm_posts_023` (`user_id`);
CREATE INDEX `idx_dm_posts_023_created_at` ON `dm_posts_023` (`created_at`);
