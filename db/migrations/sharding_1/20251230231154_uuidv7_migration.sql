-- UUIDv7 Migration: Change ID columns from bigint to varchar(32)
-- dm_users_000 to dm_users_007
-- dm_posts_000 to dm_posts_007

-- Delete existing data and recreate tables with new schema

-- dm_users_000
DROP TABLE IF EXISTS `dm_users_000`;
CREATE TABLE `dm_users_000` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_000_email` ON `dm_users_000` (`email`);

-- dm_users_001
DROP TABLE IF EXISTS `dm_users_001`;
CREATE TABLE `dm_users_001` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_001_email` ON `dm_users_001` (`email`);

-- dm_users_002
DROP TABLE IF EXISTS `dm_users_002`;
CREATE TABLE `dm_users_002` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_002_email` ON `dm_users_002` (`email`);

-- dm_users_003
DROP TABLE IF EXISTS `dm_users_003`;
CREATE TABLE `dm_users_003` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_003_email` ON `dm_users_003` (`email`);

-- dm_users_004
DROP TABLE IF EXISTS `dm_users_004`;
CREATE TABLE `dm_users_004` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_004_email` ON `dm_users_004` (`email`);

-- dm_users_005
DROP TABLE IF EXISTS `dm_users_005`;
CREATE TABLE `dm_users_005` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_005_email` ON `dm_users_005` (`email`);

-- dm_users_006
DROP TABLE IF EXISTS `dm_users_006`;
CREATE TABLE `dm_users_006` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_006_email` ON `dm_users_006` (`email`);

-- dm_users_007
DROP TABLE IF EXISTS `dm_users_007`;
CREATE TABLE `dm_users_007` (
  `id` varchar(32) NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_dm_users_007_email` ON `dm_users_007` (`email`);

-- dm_posts_000
DROP TABLE IF EXISTS `dm_posts_000`;
CREATE TABLE `dm_posts_000` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_000_user_id` ON `dm_posts_000` (`user_id`);
CREATE INDEX `idx_dm_posts_000_created_at` ON `dm_posts_000` (`created_at`);

-- dm_posts_001
DROP TABLE IF EXISTS `dm_posts_001`;
CREATE TABLE `dm_posts_001` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_001_user_id` ON `dm_posts_001` (`user_id`);
CREATE INDEX `idx_dm_posts_001_created_at` ON `dm_posts_001` (`created_at`);

-- dm_posts_002
DROP TABLE IF EXISTS `dm_posts_002`;
CREATE TABLE `dm_posts_002` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_002_user_id` ON `dm_posts_002` (`user_id`);
CREATE INDEX `idx_dm_posts_002_created_at` ON `dm_posts_002` (`created_at`);

-- dm_posts_003
DROP TABLE IF EXISTS `dm_posts_003`;
CREATE TABLE `dm_posts_003` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_003_user_id` ON `dm_posts_003` (`user_id`);
CREATE INDEX `idx_dm_posts_003_created_at` ON `dm_posts_003` (`created_at`);

-- dm_posts_004
DROP TABLE IF EXISTS `dm_posts_004`;
CREATE TABLE `dm_posts_004` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_004_user_id` ON `dm_posts_004` (`user_id`);
CREATE INDEX `idx_dm_posts_004_created_at` ON `dm_posts_004` (`created_at`);

-- dm_posts_005
DROP TABLE IF EXISTS `dm_posts_005`;
CREATE TABLE `dm_posts_005` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_005_user_id` ON `dm_posts_005` (`user_id`);
CREATE INDEX `idx_dm_posts_005_created_at` ON `dm_posts_005` (`created_at`);

-- dm_posts_006
DROP TABLE IF EXISTS `dm_posts_006`;
CREATE TABLE `dm_posts_006` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_006_user_id` ON `dm_posts_006` (`user_id`);
CREATE INDEX `idx_dm_posts_006_created_at` ON `dm_posts_006` (`created_at`);

-- dm_posts_007
DROP TABLE IF EXISTS `dm_posts_007`;
CREATE TABLE `dm_posts_007` (
  `id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_dm_posts_007_user_id` ON `dm_posts_007` (`user_id`);
CREATE INDEX `idx_dm_posts_007_created_at` ON `dm_posts_007` (`created_at`);
