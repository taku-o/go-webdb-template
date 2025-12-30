-- UUIDv7 Migration: Change ID columns from bigint to varchar(32)
-- dm_users_024 to dm_users_031
-- dm_posts_024 to dm_posts_031

-- Delete existing data and recreate tables with new schema

-- dm_users_024
DROP TABLE IF EXISTS `dm_users_024`;
CREATE TABLE `dm_users_024` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_024_email` ON `dm_users_024` (`email`);

-- dm_users_025
DROP TABLE IF EXISTS `dm_users_025`;
CREATE TABLE `dm_users_025` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_025_email` ON `dm_users_025` (`email`);

-- dm_users_026
DROP TABLE IF EXISTS `dm_users_026`;
CREATE TABLE `dm_users_026` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_026_email` ON `dm_users_026` (`email`);

-- dm_users_027
DROP TABLE IF EXISTS `dm_users_027`;
CREATE TABLE `dm_users_027` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_027_email` ON `dm_users_027` (`email`);

-- dm_users_028
DROP TABLE IF EXISTS `dm_users_028`;
CREATE TABLE `dm_users_028` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_028_email` ON `dm_users_028` (`email`);

-- dm_users_029
DROP TABLE IF EXISTS `dm_users_029`;
CREATE TABLE `dm_users_029` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_029_email` ON `dm_users_029` (`email`);

-- dm_users_030
DROP TABLE IF EXISTS `dm_users_030`;
CREATE TABLE `dm_users_030` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_030_email` ON `dm_users_030` (`email`);

-- dm_users_031
DROP TABLE IF EXISTS `dm_users_031`;
CREATE TABLE `dm_users_031` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_031_email` ON `dm_users_031` (`email`);

-- dm_posts_024
DROP TABLE IF EXISTS `dm_posts_024`;
CREATE TABLE `dm_posts_024` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_024_user_id` ON `dm_posts_024` (`user_id`);
CREATE INDEX `idx_dm_posts_024_created_at` ON `dm_posts_024` (`created_at`);

-- dm_posts_025
DROP TABLE IF EXISTS `dm_posts_025`;
CREATE TABLE `dm_posts_025` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_025_user_id` ON `dm_posts_025` (`user_id`);
CREATE INDEX `idx_dm_posts_025_created_at` ON `dm_posts_025` (`created_at`);

-- dm_posts_026
DROP TABLE IF EXISTS `dm_posts_026`;
CREATE TABLE `dm_posts_026` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_026_user_id` ON `dm_posts_026` (`user_id`);
CREATE INDEX `idx_dm_posts_026_created_at` ON `dm_posts_026` (`created_at`);

-- dm_posts_027
DROP TABLE IF EXISTS `dm_posts_027`;
CREATE TABLE `dm_posts_027` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_027_user_id` ON `dm_posts_027` (`user_id`);
CREATE INDEX `idx_dm_posts_027_created_at` ON `dm_posts_027` (`created_at`);

-- dm_posts_028
DROP TABLE IF EXISTS `dm_posts_028`;
CREATE TABLE `dm_posts_028` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_028_user_id` ON `dm_posts_028` (`user_id`);
CREATE INDEX `idx_dm_posts_028_created_at` ON `dm_posts_028` (`created_at`);

-- dm_posts_029
DROP TABLE IF EXISTS `dm_posts_029`;
CREATE TABLE `dm_posts_029` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_029_user_id` ON `dm_posts_029` (`user_id`);
CREATE INDEX `idx_dm_posts_029_created_at` ON `dm_posts_029` (`created_at`);

-- dm_posts_030
DROP TABLE IF EXISTS `dm_posts_030`;
CREATE TABLE `dm_posts_030` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_030_user_id` ON `dm_posts_030` (`user_id`);
CREATE INDEX `idx_dm_posts_030_created_at` ON `dm_posts_030` (`created_at`);

-- dm_posts_031
DROP TABLE IF EXISTS `dm_posts_031`;
CREATE TABLE `dm_posts_031` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_031_user_id` ON `dm_posts_031` (`user_id`);
CREATE INDEX `idx_dm_posts_031_created_at` ON `dm_posts_031` (`created_at`);
