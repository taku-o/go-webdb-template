-- Create "posts_000" table
CREATE TABLE `posts_000` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_000_user_id" to table: "posts_000"
CREATE INDEX `idx_posts_000_user_id` ON `posts_000` (`user_id`);
-- Create index "idx_posts_000_created_at" to table: "posts_000"
CREATE INDEX `idx_posts_000_created_at` ON `posts_000` (`created_at`);
-- Create "posts_001" table
CREATE TABLE `posts_001` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_001_user_id" to table: "posts_001"
CREATE INDEX `idx_posts_001_user_id` ON `posts_001` (`user_id`);
-- Create index "idx_posts_001_created_at" to table: "posts_001"
CREATE INDEX `idx_posts_001_created_at` ON `posts_001` (`created_at`);
-- Create "posts_002" table
CREATE TABLE `posts_002` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_002_user_id" to table: "posts_002"
CREATE INDEX `idx_posts_002_user_id` ON `posts_002` (`user_id`);
-- Create index "idx_posts_002_created_at" to table: "posts_002"
CREATE INDEX `idx_posts_002_created_at` ON `posts_002` (`created_at`);
-- Create "posts_003" table
CREATE TABLE `posts_003` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_003_user_id" to table: "posts_003"
CREATE INDEX `idx_posts_003_user_id` ON `posts_003` (`user_id`);
-- Create index "idx_posts_003_created_at" to table: "posts_003"
CREATE INDEX `idx_posts_003_created_at` ON `posts_003` (`created_at`);
-- Create "posts_004" table
CREATE TABLE `posts_004` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_004_user_id" to table: "posts_004"
CREATE INDEX `idx_posts_004_user_id` ON `posts_004` (`user_id`);
-- Create index "idx_posts_004_created_at" to table: "posts_004"
CREATE INDEX `idx_posts_004_created_at` ON `posts_004` (`created_at`);
-- Create "posts_005" table
CREATE TABLE `posts_005` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_005_user_id" to table: "posts_005"
CREATE INDEX `idx_posts_005_user_id` ON `posts_005` (`user_id`);
-- Create index "idx_posts_005_created_at" to table: "posts_005"
CREATE INDEX `idx_posts_005_created_at` ON `posts_005` (`created_at`);
-- Create "posts_006" table
CREATE TABLE `posts_006` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_006_user_id" to table: "posts_006"
CREATE INDEX `idx_posts_006_user_id` ON `posts_006` (`user_id`);
-- Create index "idx_posts_006_created_at" to table: "posts_006"
CREATE INDEX `idx_posts_006_created_at` ON `posts_006` (`created_at`);
-- Create "posts_007" table
CREATE TABLE `posts_007` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_007_user_id" to table: "posts_007"
CREATE INDEX `idx_posts_007_user_id` ON `posts_007` (`user_id`);
-- Create index "idx_posts_007_created_at" to table: "posts_007"
CREATE INDEX `idx_posts_007_created_at` ON `posts_007` (`created_at`);
-- Create "users_000" table
CREATE TABLE `users_000` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_000_email" to table: "users_000"
CREATE UNIQUE INDEX `idx_users_000_email` ON `users_000` (`email`);
-- Create "users_001" table
CREATE TABLE `users_001` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_001_email" to table: "users_001"
CREATE UNIQUE INDEX `idx_users_001_email` ON `users_001` (`email`);
-- Create "users_002" table
CREATE TABLE `users_002` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_002_email" to table: "users_002"
CREATE UNIQUE INDEX `idx_users_002_email` ON `users_002` (`email`);
-- Create "users_003" table
CREATE TABLE `users_003` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_003_email" to table: "users_003"
CREATE UNIQUE INDEX `idx_users_003_email` ON `users_003` (`email`);
-- Create "users_004" table
CREATE TABLE `users_004` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_004_email" to table: "users_004"
CREATE UNIQUE INDEX `idx_users_004_email` ON `users_004` (`email`);
-- Create "users_005" table
CREATE TABLE `users_005` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_005_email" to table: "users_005"
CREATE UNIQUE INDEX `idx_users_005_email` ON `users_005` (`email`);
-- Create "users_006" table
CREATE TABLE `users_006` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_006_email" to table: "users_006"
CREATE UNIQUE INDEX `idx_users_006_email` ON `users_006` (`email`);
-- Create "users_007" table
CREATE TABLE `users_007` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_007_email" to table: "users_007"
CREATE UNIQUE INDEX `idx_users_007_email` ON `users_007` (`email`);
