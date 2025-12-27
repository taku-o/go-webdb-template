-- Create "posts_008" table
CREATE TABLE `posts_008` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_008_user_id" to table: "posts_008"
CREATE INDEX `idx_posts_008_user_id` ON `posts_008` (`user_id`);
-- Create index "idx_posts_008_created_at" to table: "posts_008"
CREATE INDEX `idx_posts_008_created_at` ON `posts_008` (`created_at`);
-- Create "posts_009" table
CREATE TABLE `posts_009` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_009_user_id" to table: "posts_009"
CREATE INDEX `idx_posts_009_user_id` ON `posts_009` (`user_id`);
-- Create index "idx_posts_009_created_at" to table: "posts_009"
CREATE INDEX `idx_posts_009_created_at` ON `posts_009` (`created_at`);
-- Create "posts_010" table
CREATE TABLE `posts_010` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_010_user_id" to table: "posts_010"
CREATE INDEX `idx_posts_010_user_id` ON `posts_010` (`user_id`);
-- Create index "idx_posts_010_created_at" to table: "posts_010"
CREATE INDEX `idx_posts_010_created_at` ON `posts_010` (`created_at`);
-- Create "posts_011" table
CREATE TABLE `posts_011` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_011_user_id" to table: "posts_011"
CREATE INDEX `idx_posts_011_user_id` ON `posts_011` (`user_id`);
-- Create index "idx_posts_011_created_at" to table: "posts_011"
CREATE INDEX `idx_posts_011_created_at` ON `posts_011` (`created_at`);
-- Create "posts_012" table
CREATE TABLE `posts_012` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_012_user_id" to table: "posts_012"
CREATE INDEX `idx_posts_012_user_id` ON `posts_012` (`user_id`);
-- Create index "idx_posts_012_created_at" to table: "posts_012"
CREATE INDEX `idx_posts_012_created_at` ON `posts_012` (`created_at`);
-- Create "posts_013" table
CREATE TABLE `posts_013` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_013_user_id" to table: "posts_013"
CREATE INDEX `idx_posts_013_user_id` ON `posts_013` (`user_id`);
-- Create index "idx_posts_013_created_at" to table: "posts_013"
CREATE INDEX `idx_posts_013_created_at` ON `posts_013` (`created_at`);
-- Create "posts_014" table
CREATE TABLE `posts_014` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_014_user_id" to table: "posts_014"
CREATE INDEX `idx_posts_014_user_id` ON `posts_014` (`user_id`);
-- Create index "idx_posts_014_created_at" to table: "posts_014"
CREATE INDEX `idx_posts_014_created_at` ON `posts_014` (`created_at`);
-- Create "posts_015" table
CREATE TABLE `posts_015` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_015_user_id" to table: "posts_015"
CREATE INDEX `idx_posts_015_user_id` ON `posts_015` (`user_id`);
-- Create index "idx_posts_015_created_at" to table: "posts_015"
CREATE INDEX `idx_posts_015_created_at` ON `posts_015` (`created_at`);
-- Create "users_008" table
CREATE TABLE `users_008` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_008_email" to table: "users_008"
CREATE UNIQUE INDEX `idx_users_008_email` ON `users_008` (`email`);
-- Create "users_009" table
CREATE TABLE `users_009` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_009_email" to table: "users_009"
CREATE UNIQUE INDEX `idx_users_009_email` ON `users_009` (`email`);
-- Create "users_010" table
CREATE TABLE `users_010` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_010_email" to table: "users_010"
CREATE UNIQUE INDEX `idx_users_010_email` ON `users_010` (`email`);
-- Create "users_011" table
CREATE TABLE `users_011` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_011_email" to table: "users_011"
CREATE UNIQUE INDEX `idx_users_011_email` ON `users_011` (`email`);
-- Create "users_012" table
CREATE TABLE `users_012` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_012_email" to table: "users_012"
CREATE UNIQUE INDEX `idx_users_012_email` ON `users_012` (`email`);
-- Create "users_013" table
CREATE TABLE `users_013` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_013_email" to table: "users_013"
CREATE UNIQUE INDEX `idx_users_013_email` ON `users_013` (`email`);
-- Create "users_014" table
CREATE TABLE `users_014` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_014_email" to table: "users_014"
CREATE UNIQUE INDEX `idx_users_014_email` ON `users_014` (`email`);
-- Create "users_015" table
CREATE TABLE `users_015` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_015_email" to table: "users_015"
CREATE UNIQUE INDEX `idx_users_015_email` ON `users_015` (`email`);
