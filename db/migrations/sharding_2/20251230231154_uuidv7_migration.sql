-- UUIDv7 Migration: Change ID columns from bigint to varchar(32)
-- dm_users_008 to dm_users_015
-- dm_posts_008 to dm_posts_015

-- Delete existing data and recreate tables with new schema

-- dm_users_008
DROP TABLE IF EXISTS `dm_users_008`;
CREATE TABLE `dm_users_008` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_008_email` ON `dm_users_008` (`email`);

-- dm_users_009
DROP TABLE IF EXISTS `dm_users_009`;
CREATE TABLE `dm_users_009` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_009_email` ON `dm_users_009` (`email`);

-- dm_users_010
DROP TABLE IF EXISTS `dm_users_010`;
CREATE TABLE `dm_users_010` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_010_email` ON `dm_users_010` (`email`);

-- dm_users_011
DROP TABLE IF EXISTS `dm_users_011`;
CREATE TABLE `dm_users_011` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_011_email` ON `dm_users_011` (`email`);

-- dm_users_012
DROP TABLE IF EXISTS `dm_users_012`;
CREATE TABLE `dm_users_012` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_012_email` ON `dm_users_012` (`email`);

-- dm_users_013
DROP TABLE IF EXISTS `dm_users_013`;
CREATE TABLE `dm_users_013` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_013_email` ON `dm_users_013` (`email`);

-- dm_users_014
DROP TABLE IF EXISTS `dm_users_014`;
CREATE TABLE `dm_users_014` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_014_email` ON `dm_users_014` (`email`);

-- dm_users_015
DROP TABLE IF EXISTS `dm_users_015`;
CREATE TABLE `dm_users_015` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_015_email` ON `dm_users_015` (`email`);

-- dm_posts_008
DROP TABLE IF EXISTS `dm_posts_008`;
CREATE TABLE `dm_posts_008` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_008_user_id` ON `dm_posts_008` (`user_id`);
CREATE INDEX `idx_dm_posts_008_created_at` ON `dm_posts_008` (`created_at`);

-- dm_posts_009
DROP TABLE IF EXISTS `dm_posts_009`;
CREATE TABLE `dm_posts_009` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_009_user_id` ON `dm_posts_009` (`user_id`);
CREATE INDEX `idx_dm_posts_009_created_at` ON `dm_posts_009` (`created_at`);

-- dm_posts_010
DROP TABLE IF EXISTS `dm_posts_010`;
CREATE TABLE `dm_posts_010` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_010_user_id` ON `dm_posts_010` (`user_id`);
CREATE INDEX `idx_dm_posts_010_created_at` ON `dm_posts_010` (`created_at`);

-- dm_posts_011
DROP TABLE IF EXISTS `dm_posts_011`;
CREATE TABLE `dm_posts_011` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_011_user_id` ON `dm_posts_011` (`user_id`);
CREATE INDEX `idx_dm_posts_011_created_at` ON `dm_posts_011` (`created_at`);

-- dm_posts_012
DROP TABLE IF EXISTS `dm_posts_012`;
CREATE TABLE `dm_posts_012` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_012_user_id` ON `dm_posts_012` (`user_id`);
CREATE INDEX `idx_dm_posts_012_created_at` ON `dm_posts_012` (`created_at`);

-- dm_posts_013
DROP TABLE IF EXISTS `dm_posts_013`;
CREATE TABLE `dm_posts_013` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_013_user_id` ON `dm_posts_013` (`user_id`);
CREATE INDEX `idx_dm_posts_013_created_at` ON `dm_posts_013` (`created_at`);

-- dm_posts_014
DROP TABLE IF EXISTS `dm_posts_014`;
CREATE TABLE `dm_posts_014` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_014_user_id` ON `dm_posts_014` (`user_id`);
CREATE INDEX `idx_dm_posts_014_created_at` ON `dm_posts_014` (`created_at`);

-- dm_posts_015
DROP TABLE IF EXISTS `dm_posts_015`;
CREATE TABLE `dm_posts_015` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_015_user_id` ON `dm_posts_015` (`user_id`);
CREATE INDEX `idx_dm_posts_015_created_at` ON `dm_posts_015` (`created_at`);
