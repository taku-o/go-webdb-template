-- Create "dm_posts_024" table
CREATE TABLE "dm_posts_024" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_024_created_at" to table: "dm_posts_024"
CREATE INDEX "idx_dm_posts_024_created_at" ON "dm_posts_024" ("created_at");
-- Create index "idx_dm_posts_024_user_id" to table: "dm_posts_024"
CREATE INDEX "idx_dm_posts_024_user_id" ON "dm_posts_024" ("user_id");
-- Create "dm_posts_025" table
CREATE TABLE "dm_posts_025" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_025_created_at" to table: "dm_posts_025"
CREATE INDEX "idx_dm_posts_025_created_at" ON "dm_posts_025" ("created_at");
-- Create index "idx_dm_posts_025_user_id" to table: "dm_posts_025"
CREATE INDEX "idx_dm_posts_025_user_id" ON "dm_posts_025" ("user_id");
-- Create "dm_posts_026" table
CREATE TABLE "dm_posts_026" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_026_created_at" to table: "dm_posts_026"
CREATE INDEX "idx_dm_posts_026_created_at" ON "dm_posts_026" ("created_at");
-- Create index "idx_dm_posts_026_user_id" to table: "dm_posts_026"
CREATE INDEX "idx_dm_posts_026_user_id" ON "dm_posts_026" ("user_id");
-- Create "dm_posts_027" table
CREATE TABLE "dm_posts_027" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_027_created_at" to table: "dm_posts_027"
CREATE INDEX "idx_dm_posts_027_created_at" ON "dm_posts_027" ("created_at");
-- Create index "idx_dm_posts_027_user_id" to table: "dm_posts_027"
CREATE INDEX "idx_dm_posts_027_user_id" ON "dm_posts_027" ("user_id");
-- Create "dm_posts_028" table
CREATE TABLE "dm_posts_028" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_028_created_at" to table: "dm_posts_028"
CREATE INDEX "idx_dm_posts_028_created_at" ON "dm_posts_028" ("created_at");
-- Create index "idx_dm_posts_028_user_id" to table: "dm_posts_028"
CREATE INDEX "idx_dm_posts_028_user_id" ON "dm_posts_028" ("user_id");
-- Create "dm_posts_029" table
CREATE TABLE "dm_posts_029" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_029_created_at" to table: "dm_posts_029"
CREATE INDEX "idx_dm_posts_029_created_at" ON "dm_posts_029" ("created_at");
-- Create index "idx_dm_posts_029_user_id" to table: "dm_posts_029"
CREATE INDEX "idx_dm_posts_029_user_id" ON "dm_posts_029" ("user_id");
-- Create "dm_posts_030" table
CREATE TABLE "dm_posts_030" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_030_created_at" to table: "dm_posts_030"
CREATE INDEX "idx_dm_posts_030_created_at" ON "dm_posts_030" ("created_at");
-- Create index "idx_dm_posts_030_user_id" to table: "dm_posts_030"
CREATE INDEX "idx_dm_posts_030_user_id" ON "dm_posts_030" ("user_id");
-- Create "dm_posts_031" table
CREATE TABLE "dm_posts_031" (
  "id" character varying(32) NOT NULL,
  "user_id" character varying(32) NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_posts_031_created_at" to table: "dm_posts_031"
CREATE INDEX "idx_dm_posts_031_created_at" ON "dm_posts_031" ("created_at");
-- Create index "idx_dm_posts_031_user_id" to table: "dm_posts_031"
CREATE INDEX "idx_dm_posts_031_user_id" ON "dm_posts_031" ("user_id");
-- Create "dm_users_024" table
CREATE TABLE "dm_users_024" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_024_email" to table: "dm_users_024"
CREATE UNIQUE INDEX "idx_dm_users_024_email" ON "dm_users_024" ("email");
-- Create "dm_users_025" table
CREATE TABLE "dm_users_025" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_025_email" to table: "dm_users_025"
CREATE UNIQUE INDEX "idx_dm_users_025_email" ON "dm_users_025" ("email");
-- Create "dm_users_026" table
CREATE TABLE "dm_users_026" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_026_email" to table: "dm_users_026"
CREATE UNIQUE INDEX "idx_dm_users_026_email" ON "dm_users_026" ("email");
-- Create "dm_users_027" table
CREATE TABLE "dm_users_027" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_027_email" to table: "dm_users_027"
CREATE UNIQUE INDEX "idx_dm_users_027_email" ON "dm_users_027" ("email");
-- Create "dm_users_028" table
CREATE TABLE "dm_users_028" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_028_email" to table: "dm_users_028"
CREATE UNIQUE INDEX "idx_dm_users_028_email" ON "dm_users_028" ("email");
-- Create "dm_users_029" table
CREATE TABLE "dm_users_029" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_029_email" to table: "dm_users_029"
CREATE UNIQUE INDEX "idx_dm_users_029_email" ON "dm_users_029" ("email");
-- Create "dm_users_030" table
CREATE TABLE "dm_users_030" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_030_email" to table: "dm_users_030"
CREATE UNIQUE INDEX "idx_dm_users_030_email" ON "dm_users_030" ("email");
-- Create "dm_users_031" table
CREATE TABLE "dm_users_031" (
  "id" character varying(32) NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_dm_users_031_email" to table: "dm_users_031"
CREATE UNIQUE INDEX "idx_dm_users_031_email" ON "dm_users_031" ("email");
