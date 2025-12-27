-- Create "posts_016" table
CREATE TABLE `posts_016` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_016_user_id" to table: "posts_016"
CREATE INDEX `idx_posts_016_user_id` ON `posts_016` (`user_id`);
-- Create index "idx_posts_016_created_at" to table: "posts_016"
CREATE INDEX `idx_posts_016_created_at` ON `posts_016` (`created_at`);
-- Create "posts_017" table
CREATE TABLE `posts_017` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_017_user_id" to table: "posts_017"
CREATE INDEX `idx_posts_017_user_id` ON `posts_017` (`user_id`);
-- Create index "idx_posts_017_created_at" to table: "posts_017"
CREATE INDEX `idx_posts_017_created_at` ON `posts_017` (`created_at`);
-- Create "posts_018" table
CREATE TABLE `posts_018` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_018_user_id" to table: "posts_018"
CREATE INDEX `idx_posts_018_user_id` ON `posts_018` (`user_id`);
-- Create index "idx_posts_018_created_at" to table: "posts_018"
CREATE INDEX `idx_posts_018_created_at` ON `posts_018` (`created_at`);
-- Create "posts_019" table
CREATE TABLE `posts_019` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_019_user_id" to table: "posts_019"
CREATE INDEX `idx_posts_019_user_id` ON `posts_019` (`user_id`);
-- Create index "idx_posts_019_created_at" to table: "posts_019"
CREATE INDEX `idx_posts_019_created_at` ON `posts_019` (`created_at`);
-- Create "posts_020" table
CREATE TABLE `posts_020` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_020_user_id" to table: "posts_020"
CREATE INDEX `idx_posts_020_user_id` ON `posts_020` (`user_id`);
-- Create index "idx_posts_020_created_at" to table: "posts_020"
CREATE INDEX `idx_posts_020_created_at` ON `posts_020` (`created_at`);
-- Create "posts_021" table
CREATE TABLE `posts_021` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_021_user_id" to table: "posts_021"
CREATE INDEX `idx_posts_021_user_id` ON `posts_021` (`user_id`);
-- Create index "idx_posts_021_created_at" to table: "posts_021"
CREATE INDEX `idx_posts_021_created_at` ON `posts_021` (`created_at`);
-- Create "posts_022" table
CREATE TABLE `posts_022` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_022_user_id" to table: "posts_022"
CREATE INDEX `idx_posts_022_user_id` ON `posts_022` (`user_id`);
-- Create index "idx_posts_022_created_at" to table: "posts_022"
CREATE INDEX `idx_posts_022_created_at` ON `posts_022` (`created_at`);
-- Create "posts_023" table
CREATE TABLE `posts_023` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_023_user_id" to table: "posts_023"
CREATE INDEX `idx_posts_023_user_id` ON `posts_023` (`user_id`);
-- Create index "idx_posts_023_created_at" to table: "posts_023"
CREATE INDEX `idx_posts_023_created_at` ON `posts_023` (`created_at`);
-- Create "users_016" table
CREATE TABLE `users_016` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_016_email" to table: "users_016"
CREATE UNIQUE INDEX `idx_users_016_email` ON `users_016` (`email`);
-- Create "users_017" table
CREATE TABLE `users_017` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_017_email" to table: "users_017"
CREATE UNIQUE INDEX `idx_users_017_email` ON `users_017` (`email`);
-- Create "users_018" table
CREATE TABLE `users_018` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_018_email" to table: "users_018"
CREATE UNIQUE INDEX `idx_users_018_email` ON `users_018` (`email`);
-- Create "users_019" table
CREATE TABLE `users_019` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_019_email" to table: "users_019"
CREATE UNIQUE INDEX `idx_users_019_email` ON `users_019` (`email`);
-- Create "users_020" table
CREATE TABLE `users_020` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_020_email" to table: "users_020"
CREATE UNIQUE INDEX `idx_users_020_email` ON `users_020` (`email`);
-- Create "users_021" table
CREATE TABLE `users_021` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_021_email" to table: "users_021"
CREATE UNIQUE INDEX `idx_users_021_email` ON `users_021` (`email`);
-- Create "users_022" table
CREATE TABLE `users_022` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_022_email" to table: "users_022"
CREATE UNIQUE INDEX `idx_users_022_email` ON `users_022` (`email`);
-- Create "users_023" table
CREATE TABLE `users_023` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_023_email" to table: "users_023"
CREATE UNIQUE INDEX `idx_users_023_email` ON `users_023` (`email`);
