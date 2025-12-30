-- Create "dm_posts_008" table
CREATE TABLE `dm_posts_008` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_008_user_id" to table: "dm_posts_008"
CREATE INDEX `idx_dm_posts_008_user_id` ON `dm_posts_008` (`user_id`);
-- Create index "idx_dm_posts_008_created_at" to table: "dm_posts_008"
CREATE INDEX `idx_dm_posts_008_created_at` ON `dm_posts_008` (`created_at`);
-- Create "dm_posts_009" table
CREATE TABLE `dm_posts_009` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_009_user_id" to table: "dm_posts_009"
CREATE INDEX `idx_dm_posts_009_user_id` ON `dm_posts_009` (`user_id`);
-- Create index "idx_dm_posts_009_created_at" to table: "dm_posts_009"
CREATE INDEX `idx_dm_posts_009_created_at` ON `dm_posts_009` (`created_at`);
-- Create "dm_posts_010" table
CREATE TABLE `dm_posts_010` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_010_user_id" to table: "dm_posts_010"
CREATE INDEX `idx_dm_posts_010_user_id` ON `dm_posts_010` (`user_id`);
-- Create index "idx_dm_posts_010_created_at" to table: "dm_posts_010"
CREATE INDEX `idx_dm_posts_010_created_at` ON `dm_posts_010` (`created_at`);
-- Create "dm_posts_011" table
CREATE TABLE `dm_posts_011` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_011_user_id" to table: "dm_posts_011"
CREATE INDEX `idx_dm_posts_011_user_id` ON `dm_posts_011` (`user_id`);
-- Create index "idx_dm_posts_011_created_at" to table: "dm_posts_011"
CREATE INDEX `idx_dm_posts_011_created_at` ON `dm_posts_011` (`created_at`);
-- Create "dm_posts_012" table
CREATE TABLE `dm_posts_012` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_012_user_id" to table: "dm_posts_012"
CREATE INDEX `idx_dm_posts_012_user_id` ON `dm_posts_012` (`user_id`);
-- Create index "idx_dm_posts_012_created_at" to table: "dm_posts_012"
CREATE INDEX `idx_dm_posts_012_created_at` ON `dm_posts_012` (`created_at`);
-- Create "dm_posts_013" table
CREATE TABLE `dm_posts_013` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_013_user_id" to table: "dm_posts_013"
CREATE INDEX `idx_dm_posts_013_user_id` ON `dm_posts_013` (`user_id`);
-- Create index "idx_dm_posts_013_created_at" to table: "dm_posts_013"
CREATE INDEX `idx_dm_posts_013_created_at` ON `dm_posts_013` (`created_at`);
-- Create "dm_posts_014" table
CREATE TABLE `dm_posts_014` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_014_user_id" to table: "dm_posts_014"
CREATE INDEX `idx_dm_posts_014_user_id` ON `dm_posts_014` (`user_id`);
-- Create index "idx_dm_posts_014_created_at" to table: "dm_posts_014"
CREATE INDEX `idx_dm_posts_014_created_at` ON `dm_posts_014` (`created_at`);
-- Create "dm_posts_015" table
CREATE TABLE `dm_posts_015` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_015_user_id" to table: "dm_posts_015"
CREATE INDEX `idx_dm_posts_015_user_id` ON `dm_posts_015` (`user_id`);
-- Create index "idx_dm_posts_015_created_at" to table: "dm_posts_015"
CREATE INDEX `idx_dm_posts_015_created_at` ON `dm_posts_015` (`created_at`);
-- Create "dm_users_008" table
CREATE TABLE `dm_users_008` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_008_email" to table: "dm_users_008"
CREATE UNIQUE INDEX `idx_dm_users_008_email` ON `dm_users_008` (`email`);
-- Create "dm_users_009" table
CREATE TABLE `dm_users_009` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_009_email" to table: "dm_users_009"
CREATE UNIQUE INDEX `idx_dm_users_009_email` ON `dm_users_009` (`email`);
-- Create "dm_users_010" table
CREATE TABLE `dm_users_010` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_010_email" to table: "dm_users_010"
CREATE UNIQUE INDEX `idx_dm_users_010_email` ON `dm_users_010` (`email`);
-- Create "dm_users_011" table
CREATE TABLE `dm_users_011` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_011_email" to table: "dm_users_011"
CREATE UNIQUE INDEX `idx_dm_users_011_email` ON `dm_users_011` (`email`);
-- Create "dm_users_012" table
CREATE TABLE `dm_users_012` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_012_email" to table: "dm_users_012"
CREATE UNIQUE INDEX `idx_dm_users_012_email` ON `dm_users_012` (`email`);
-- Create "dm_users_013" table
CREATE TABLE `dm_users_013` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_013_email" to table: "dm_users_013"
CREATE UNIQUE INDEX `idx_dm_users_013_email` ON `dm_users_013` (`email`);
-- Create "dm_users_014" table
CREATE TABLE `dm_users_014` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_014_email" to table: "dm_users_014"
CREATE UNIQUE INDEX `idx_dm_users_014_email` ON `dm_users_014` (`email`);
-- Create "dm_users_015" table
CREATE TABLE `dm_users_015` (
  `id` bigint NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_015_email" to table: "dm_users_015"
CREATE UNIQUE INDEX `idx_dm_users_015_email` ON `dm_users_015` (`email`);
