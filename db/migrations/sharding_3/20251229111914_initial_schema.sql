-- Create "dm_posts_016" table
CREATE TABLE `dm_posts_016` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_016_user_id" to table: "dm_posts_016"
CREATE INDEX `idx_dm_posts_016_user_id` ON `dm_posts_016` (`user_id`);
-- Create index "idx_dm_posts_016_created_at" to table: "dm_posts_016"
CREATE INDEX `idx_dm_posts_016_created_at` ON `dm_posts_016` (`created_at`);
-- Create "dm_posts_017" table
CREATE TABLE `dm_posts_017` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_017_user_id" to table: "dm_posts_017"
CREATE INDEX `idx_dm_posts_017_user_id` ON `dm_posts_017` (`user_id`);
-- Create index "idx_dm_posts_017_created_at" to table: "dm_posts_017"
CREATE INDEX `idx_dm_posts_017_created_at` ON `dm_posts_017` (`created_at`);
-- Create "dm_posts_018" table
CREATE TABLE `dm_posts_018` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_018_user_id" to table: "dm_posts_018"
CREATE INDEX `idx_dm_posts_018_user_id` ON `dm_posts_018` (`user_id`);
-- Create index "idx_dm_posts_018_created_at" to table: "dm_posts_018"
CREATE INDEX `idx_dm_posts_018_created_at` ON `dm_posts_018` (`created_at`);
-- Create "dm_posts_019" table
CREATE TABLE `dm_posts_019` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_019_user_id" to table: "dm_posts_019"
CREATE INDEX `idx_dm_posts_019_user_id` ON `dm_posts_019` (`user_id`);
-- Create index "idx_dm_posts_019_created_at" to table: "dm_posts_019"
CREATE INDEX `idx_dm_posts_019_created_at` ON `dm_posts_019` (`created_at`);
-- Create "dm_posts_020" table
CREATE TABLE `dm_posts_020` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_020_user_id" to table: "dm_posts_020"
CREATE INDEX `idx_dm_posts_020_user_id` ON `dm_posts_020` (`user_id`);
-- Create index "idx_dm_posts_020_created_at" to table: "dm_posts_020"
CREATE INDEX `idx_dm_posts_020_created_at` ON `dm_posts_020` (`created_at`);
-- Create "dm_posts_021" table
CREATE TABLE `dm_posts_021` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_021_user_id" to table: "dm_posts_021"
CREATE INDEX `idx_dm_posts_021_user_id` ON `dm_posts_021` (`user_id`);
-- Create index "idx_dm_posts_021_created_at" to table: "dm_posts_021"
CREATE INDEX `idx_dm_posts_021_created_at` ON `dm_posts_021` (`created_at`);
-- Create "dm_posts_022" table
CREATE TABLE `dm_posts_022` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_022_user_id" to table: "dm_posts_022"
CREATE INDEX `idx_dm_posts_022_user_id` ON `dm_posts_022` (`user_id`);
-- Create index "idx_dm_posts_022_created_at" to table: "dm_posts_022"
CREATE INDEX `idx_dm_posts_022_created_at` ON `dm_posts_022` (`created_at`);
-- Create "dm_posts_023" table
CREATE TABLE `dm_posts_023` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_posts_023_user_id" to table: "dm_posts_023"
CREATE INDEX `idx_dm_posts_023_user_id` ON `dm_posts_023` (`user_id`);
-- Create index "idx_dm_posts_023_created_at" to table: "dm_posts_023"
CREATE INDEX `idx_dm_posts_023_created_at` ON `dm_posts_023` (`created_at`);
-- Create "dm_users_016" table
CREATE TABLE `dm_users_016` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_016_email" to table: "dm_users_016"
CREATE UNIQUE INDEX `idx_dm_users_016_email` ON `dm_users_016` (`email`);
-- Create "dm_users_017" table
CREATE TABLE `dm_users_017` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_017_email" to table: "dm_users_017"
CREATE UNIQUE INDEX `idx_dm_users_017_email` ON `dm_users_017` (`email`);
-- Create "dm_users_018" table
CREATE TABLE `dm_users_018` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_018_email" to table: "dm_users_018"
CREATE UNIQUE INDEX `idx_dm_users_018_email` ON `dm_users_018` (`email`);
-- Create "dm_users_019" table
CREATE TABLE `dm_users_019` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_019_email" to table: "dm_users_019"
CREATE UNIQUE INDEX `idx_dm_users_019_email` ON `dm_users_019` (`email`);
-- Create "dm_users_020" table
CREATE TABLE `dm_users_020` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_020_email" to table: "dm_users_020"
CREATE UNIQUE INDEX `idx_dm_users_020_email` ON `dm_users_020` (`email`);
-- Create "dm_users_021" table
CREATE TABLE `dm_users_021` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_021_email" to table: "dm_users_021"
CREATE UNIQUE INDEX `idx_dm_users_021_email` ON `dm_users_021` (`email`);
-- Create "dm_users_022" table
CREATE TABLE `dm_users_022` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_022_email" to table: "dm_users_022"
CREATE UNIQUE INDEX `idx_dm_users_022_email` ON `dm_users_022` (`email`);
-- Create "dm_users_023" table
CREATE TABLE `dm_users_023` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_dm_users_023_email" to table: "dm_users_023"
CREATE UNIQUE INDEX `idx_dm_users_023_email` ON `dm_users_023` (`email`);
