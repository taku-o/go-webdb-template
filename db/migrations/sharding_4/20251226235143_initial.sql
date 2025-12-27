-- Create "posts_024" table
CREATE TABLE `posts_024` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_024_user_id" to table: "posts_024"
CREATE INDEX `idx_posts_024_user_id` ON `posts_024` (`user_id`);
-- Create index "idx_posts_024_created_at" to table: "posts_024"
CREATE INDEX `idx_posts_024_created_at` ON `posts_024` (`created_at`);
-- Create "posts_025" table
CREATE TABLE `posts_025` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_025_user_id" to table: "posts_025"
CREATE INDEX `idx_posts_025_user_id` ON `posts_025` (`user_id`);
-- Create index "idx_posts_025_created_at" to table: "posts_025"
CREATE INDEX `idx_posts_025_created_at` ON `posts_025` (`created_at`);
-- Create "posts_026" table
CREATE TABLE `posts_026` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_026_user_id" to table: "posts_026"
CREATE INDEX `idx_posts_026_user_id` ON `posts_026` (`user_id`);
-- Create index "idx_posts_026_created_at" to table: "posts_026"
CREATE INDEX `idx_posts_026_created_at` ON `posts_026` (`created_at`);
-- Create "posts_027" table
CREATE TABLE `posts_027` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_027_user_id" to table: "posts_027"
CREATE INDEX `idx_posts_027_user_id` ON `posts_027` (`user_id`);
-- Create index "idx_posts_027_created_at" to table: "posts_027"
CREATE INDEX `idx_posts_027_created_at` ON `posts_027` (`created_at`);
-- Create "posts_028" table
CREATE TABLE `posts_028` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_028_user_id" to table: "posts_028"
CREATE INDEX `idx_posts_028_user_id` ON `posts_028` (`user_id`);
-- Create index "idx_posts_028_created_at" to table: "posts_028"
CREATE INDEX `idx_posts_028_created_at` ON `posts_028` (`created_at`);
-- Create "posts_029" table
CREATE TABLE `posts_029` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_029_user_id" to table: "posts_029"
CREATE INDEX `idx_posts_029_user_id` ON `posts_029` (`user_id`);
-- Create index "idx_posts_029_created_at" to table: "posts_029"
CREATE INDEX `idx_posts_029_created_at` ON `posts_029` (`created_at`);
-- Create "posts_030" table
CREATE TABLE `posts_030` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_030_user_id" to table: "posts_030"
CREATE INDEX `idx_posts_030_user_id` ON `posts_030` (`user_id`);
-- Create index "idx_posts_030_created_at" to table: "posts_030"
CREATE INDEX `idx_posts_030_created_at` ON `posts_030` (`created_at`);
-- Create "posts_031" table
CREATE TABLE `posts_031` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_posts_031_user_id" to table: "posts_031"
CREATE INDEX `idx_posts_031_user_id` ON `posts_031` (`user_id`);
-- Create index "idx_posts_031_created_at" to table: "posts_031"
CREATE INDEX `idx_posts_031_created_at` ON `posts_031` (`created_at`);
-- Create "users_024" table
CREATE TABLE `users_024` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_024_email" to table: "users_024"
CREATE UNIQUE INDEX `idx_users_024_email` ON `users_024` (`email`);
-- Create "users_025" table
CREATE TABLE `users_025` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_025_email" to table: "users_025"
CREATE UNIQUE INDEX `idx_users_025_email` ON `users_025` (`email`);
-- Create "users_026" table
CREATE TABLE `users_026` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_026_email" to table: "users_026"
CREATE UNIQUE INDEX `idx_users_026_email` ON `users_026` (`email`);
-- Create "users_027" table
CREATE TABLE `users_027` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_027_email" to table: "users_027"
CREATE UNIQUE INDEX `idx_users_027_email` ON `users_027` (`email`);
-- Create "users_028" table
CREATE TABLE `users_028` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_028_email" to table: "users_028"
CREATE UNIQUE INDEX `idx_users_028_email` ON `users_028` (`email`);
-- Create "users_029" table
CREATE TABLE `users_029` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_029_email" to table: "users_029"
CREATE UNIQUE INDEX `idx_users_029_email` ON `users_029` (`email`);
-- Create "users_030" table
CREATE TABLE `users_030` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_030_email" to table: "users_030"
CREATE UNIQUE INDEX `idx_users_030_email` ON `users_030` (`email`);
-- Create "users_031" table
CREATE TABLE `users_031` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_031_email" to table: "users_031"
CREATE UNIQUE INDEX `idx_users_031_email` ON `users_031` (`email`);
