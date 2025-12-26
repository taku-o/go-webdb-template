// Sharding Database Schema
// シャーディングデータベースのスキーマ定義（users, postsテーブル 32分割）

schema "main" {
}

// Users テーブル（32分割: users_000 ～ users_031）

table "users_000" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_000_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_001" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_001_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_002" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_002_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_003" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_003_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_004" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_004_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_005" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_005_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_006" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_006_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_007" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_007_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_008" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_008_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_009" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_009_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_010" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_010_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_011" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_011_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_012" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_012_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_013" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_013_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_014" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_014_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_015" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_015_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_016" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_016_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_017" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_017_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_018" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_018_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_019" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_019_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_020" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_020_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_021" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_021_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_022" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_022_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_023" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_023_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_024" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_024_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_025" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_025_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_026" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_026_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_027" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_027_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_028" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_028_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_029" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_029_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_030" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_030_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_031" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_031_email" {
    unique  = true
    columns = [column.email]
  }
}

// Posts テーブル（32分割: posts_000 ～ posts_031）

table "posts_000" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_000_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_000.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_000_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_000_created_at" {
    columns = [column.created_at]
  }
}

table "posts_001" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_001_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_001.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_001_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_001_created_at" {
    columns = [column.created_at]
  }
}

table "posts_002" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_002_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_002.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_002_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_002_created_at" {
    columns = [column.created_at]
  }
}

table "posts_003" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_003_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_003.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_003_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_003_created_at" {
    columns = [column.created_at]
  }
}

table "posts_004" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_004_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_004.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_004_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_004_created_at" {
    columns = [column.created_at]
  }
}

table "posts_005" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_005_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_005.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_005_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_005_created_at" {
    columns = [column.created_at]
  }
}

table "posts_006" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_006_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_006.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_006_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_006_created_at" {
    columns = [column.created_at]
  }
}

table "posts_007" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_007_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_007.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_007_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_007_created_at" {
    columns = [column.created_at]
  }
}

table "posts_008" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_008_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_008.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_008_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_008_created_at" {
    columns = [column.created_at]
  }
}

table "posts_009" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_009_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_009.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_009_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_009_created_at" {
    columns = [column.created_at]
  }
}

table "posts_010" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_010_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_010.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_010_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_010_created_at" {
    columns = [column.created_at]
  }
}

table "posts_011" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_011_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_011.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_011_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_011_created_at" {
    columns = [column.created_at]
  }
}

table "posts_012" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_012_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_012.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_012_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_012_created_at" {
    columns = [column.created_at]
  }
}

table "posts_013" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_013_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_013.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_013_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_013_created_at" {
    columns = [column.created_at]
  }
}

table "posts_014" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_014_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_014.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_014_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_014_created_at" {
    columns = [column.created_at]
  }
}

table "posts_015" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_015_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_015.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_015_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_015_created_at" {
    columns = [column.created_at]
  }
}

table "posts_016" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_016_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_016.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_016_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_016_created_at" {
    columns = [column.created_at]
  }
}

table "posts_017" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_017_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_017.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_017_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_017_created_at" {
    columns = [column.created_at]
  }
}

table "posts_018" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_018_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_018.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_018_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_018_created_at" {
    columns = [column.created_at]
  }
}

table "posts_019" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_019_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_019.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_019_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_019_created_at" {
    columns = [column.created_at]
  }
}

table "posts_020" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_020_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_020.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_020_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_020_created_at" {
    columns = [column.created_at]
  }
}

table "posts_021" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_021_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_021.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_021_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_021_created_at" {
    columns = [column.created_at]
  }
}

table "posts_022" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_022_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_022.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_022_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_022_created_at" {
    columns = [column.created_at]
  }
}

table "posts_023" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_023_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_023.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_023_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_023_created_at" {
    columns = [column.created_at]
  }
}

table "posts_024" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_024_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_024.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_024_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_024_created_at" {
    columns = [column.created_at]
  }
}

table "posts_025" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_025_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_025.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_025_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_025_created_at" {
    columns = [column.created_at]
  }
}

table "posts_026" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_026_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_026.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_026_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_026_created_at" {
    columns = [column.created_at]
  }
}

table "posts_027" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_027_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_027.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_027_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_027_created_at" {
    columns = [column.created_at]
  }
}

table "posts_028" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_028_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_028.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_028_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_028_created_at" {
    columns = [column.created_at]
  }
}

table "posts_029" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_029_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_029.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_029_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_029_created_at" {
    columns = [column.created_at]
  }
}

table "posts_030" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_030_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_030.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_030_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_030_created_at" {
    columns = [column.created_at]
  }
}

table "posts_031" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_031_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users_031.column.id]
    on_delete   = CASCADE
  }
  index "idx_posts_031_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_031_created_at" {
    columns = [column.created_at]
  }
}
