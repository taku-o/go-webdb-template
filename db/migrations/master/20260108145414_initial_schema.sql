-- Create "dm_news" table
CREATE TABLE "dm_news" (
  "id" serial NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "author_id" integer NULL,
  "published_at" timestamp NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_news_author_id" to table: "dm_news"
CREATE INDEX "idx_dm_news_author_id" ON "dm_news" ("author_id");
-- Create index "idx_dm_news_published_at" to table: "dm_news"
CREATE INDEX "idx_dm_news_published_at" ON "dm_news" ("published_at");
-- Create "goadmin_menu" table
CREATE TABLE "goadmin_menu" (
  "id" serial NOT NULL,
  "parent_id" integer NOT NULL DEFAULT 0,
  "type" integer NOT NULL DEFAULT 0,
  "order" integer NOT NULL DEFAULT 0,
  "title" text NOT NULL,
  "icon" text NOT NULL,
  "uri" text NOT NULL DEFAULT '',
  "header" text NULL,
  "plugin_name" text NOT NULL DEFAULT '',
  "uuid" text NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create "goadmin_operation_log" table
CREATE TABLE "goadmin_operation_log" (
  "id" serial NOT NULL,
  "user_id" integer NOT NULL,
  "path" text NOT NULL,
  "method" text NOT NULL,
  "ip" text NOT NULL,
  "input" text NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create index "idx_goadmin_operation_log_user_id" to table: "goadmin_operation_log"
CREATE INDEX "idx_goadmin_operation_log_user_id" ON "goadmin_operation_log" ("user_id");
-- Create "goadmin_permissions" table
CREATE TABLE "goadmin_permissions" (
  "id" serial NOT NULL,
  "name" text NOT NULL,
  "slug" text NOT NULL,
  "http_method" text NULL,
  "http_path" text NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create index "idx_goadmin_permissions_slug" to table: "goadmin_permissions"
CREATE UNIQUE INDEX "idx_goadmin_permissions_slug" ON "goadmin_permissions" ("slug");
-- Create "goadmin_role_menu" table
CREATE TABLE "goadmin_role_menu" (
  "role_id" integer NOT NULL,
  "menu_id" integer NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("role_id", "menu_id")
);
-- Create "goadmin_role_permissions" table
CREATE TABLE "goadmin_role_permissions" (
  "role_id" integer NOT NULL,
  "permission_id" integer NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("role_id", "permission_id")
);
-- Create "goadmin_role_users" table
CREATE TABLE "goadmin_role_users" (
  "role_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("role_id", "user_id")
);
-- Create "goadmin_roles" table
CREATE TABLE "goadmin_roles" (
  "id" serial NOT NULL,
  "name" text NOT NULL,
  "slug" text NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create index "idx_goadmin_roles_slug" to table: "goadmin_roles"
CREATE UNIQUE INDEX "idx_goadmin_roles_slug" ON "goadmin_roles" ("slug");
-- Create "goadmin_session" table
CREATE TABLE "goadmin_session" (
  "id" serial NOT NULL,
  "sid" text NOT NULL,
  "values" text NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create "goadmin_site" table
CREATE TABLE "goadmin_site" (
  "id" serial NOT NULL,
  "key" text NULL,
  "value" text NULL,
  "description" text NULL,
  "state" integer NOT NULL DEFAULT 0,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create "goadmin_user_permissions" table
CREATE TABLE "goadmin_user_permissions" (
  "user_id" integer NOT NULL,
  "permission_id" integer NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("user_id", "permission_id")
);
-- Create "goadmin_users" table
CREATE TABLE "goadmin_users" (
  "id" serial NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "name" text NOT NULL,
  "avatar" text NULL,
  "remember_token" text NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create index "idx_goadmin_users_username" to table: "goadmin_users"
CREATE UNIQUE INDEX "idx_goadmin_users_username" ON "goadmin_users" ("username");
