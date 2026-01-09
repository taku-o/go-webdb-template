-- Create "dm_posts_000" table
CREATE TABLE "dm_posts_000" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_000_created_at" to table: "dm_posts_000"
CREATE INDEX "idx_dm_posts_000_created_at" ON "dm_posts_000" ("created_at");
-- Create index "idx_dm_posts_000_user_id" to table: "dm_posts_000"
CREATE INDEX "idx_dm_posts_000_user_id" ON "dm_posts_000" ("user_id");
-- Create "dm_posts_001" table
CREATE TABLE "dm_posts_001" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_001_created_at" to table: "dm_posts_001"
CREATE INDEX "idx_dm_posts_001_created_at" ON "dm_posts_001" ("created_at");
-- Create index "idx_dm_posts_001_user_id" to table: "dm_posts_001"
CREATE INDEX "idx_dm_posts_001_user_id" ON "dm_posts_001" ("user_id");
-- Create "dm_posts_002" table
CREATE TABLE "dm_posts_002" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_002_created_at" to table: "dm_posts_002"
CREATE INDEX "idx_dm_posts_002_created_at" ON "dm_posts_002" ("created_at");
-- Create index "idx_dm_posts_002_user_id" to table: "dm_posts_002"
CREATE INDEX "idx_dm_posts_002_user_id" ON "dm_posts_002" ("user_id");
-- Create "dm_posts_003" table
CREATE TABLE "dm_posts_003" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_003_created_at" to table: "dm_posts_003"
CREATE INDEX "idx_dm_posts_003_created_at" ON "dm_posts_003" ("created_at");
-- Create index "idx_dm_posts_003_user_id" to table: "dm_posts_003"
CREATE INDEX "idx_dm_posts_003_user_id" ON "dm_posts_003" ("user_id");
-- Create "dm_posts_004" table
CREATE TABLE "dm_posts_004" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_004_created_at" to table: "dm_posts_004"
CREATE INDEX "idx_dm_posts_004_created_at" ON "dm_posts_004" ("created_at");
-- Create index "idx_dm_posts_004_user_id" to table: "dm_posts_004"
CREATE INDEX "idx_dm_posts_004_user_id" ON "dm_posts_004" ("user_id");
-- Create "dm_posts_005" table
CREATE TABLE "dm_posts_005" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_005_created_at" to table: "dm_posts_005"
CREATE INDEX "idx_dm_posts_005_created_at" ON "dm_posts_005" ("created_at");
-- Create index "idx_dm_posts_005_user_id" to table: "dm_posts_005"
CREATE INDEX "idx_dm_posts_005_user_id" ON "dm_posts_005" ("user_id");
-- Create "dm_posts_006" table
CREATE TABLE "dm_posts_006" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_006_created_at" to table: "dm_posts_006"
CREATE INDEX "idx_dm_posts_006_created_at" ON "dm_posts_006" ("created_at");
-- Create index "idx_dm_posts_006_user_id" to table: "dm_posts_006"
CREATE INDEX "idx_dm_posts_006_user_id" ON "dm_posts_006" ("user_id");
-- Create "dm_posts_007" table
CREATE TABLE "dm_posts_007" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_007_created_at" to table: "dm_posts_007"
CREATE INDEX "idx_dm_posts_007_created_at" ON "dm_posts_007" ("created_at");
-- Create index "idx_dm_posts_007_user_id" to table: "dm_posts_007"
CREATE INDEX "idx_dm_posts_007_user_id" ON "dm_posts_007" ("user_id");
-- Create "dm_users_000" table
CREATE TABLE "dm_users_000" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_000_email" to table: "dm_users_000"
CREATE UNIQUE INDEX "idx_dm_users_000_email" ON "dm_users_000" ("email");
-- Create "dm_users_001" table
CREATE TABLE "dm_users_001" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_001_email" to table: "dm_users_001"
CREATE UNIQUE INDEX "idx_dm_users_001_email" ON "dm_users_001" ("email");
-- Create "dm_users_002" table
CREATE TABLE "dm_users_002" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_002_email" to table: "dm_users_002"
CREATE UNIQUE INDEX "idx_dm_users_002_email" ON "dm_users_002" ("email");
-- Create "dm_users_003" table
CREATE TABLE "dm_users_003" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_003_email" to table: "dm_users_003"
CREATE UNIQUE INDEX "idx_dm_users_003_email" ON "dm_users_003" ("email");
-- Create "dm_users_004" table
CREATE TABLE "dm_users_004" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_004_email" to table: "dm_users_004"
CREATE UNIQUE INDEX "idx_dm_users_004_email" ON "dm_users_004" ("email");
-- Create "dm_users_005" table
CREATE TABLE "dm_users_005" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_005_email" to table: "dm_users_005"
CREATE UNIQUE INDEX "idx_dm_users_005_email" ON "dm_users_005" ("email");
-- Create "dm_users_006" table
CREATE TABLE "dm_users_006" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_006_email" to table: "dm_users_006"
CREATE UNIQUE INDEX "idx_dm_users_006_email" ON "dm_users_006" ("email");
-- Create "dm_users_007" table
CREATE TABLE "dm_users_007" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_007_email" to table: "dm_users_007"
CREATE UNIQUE INDEX "idx_dm_users_007_email" ON "dm_users_007" ("email");
