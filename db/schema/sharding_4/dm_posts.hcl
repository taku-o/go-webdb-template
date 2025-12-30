// dm_posts テーブル（sharding_db_4.db用: dm_posts_024 〜 dm_posts_031）

table "dm_posts_024" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_024_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_024_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_025" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_025_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_025_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_026" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_026_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_026_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_027" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_027_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_027_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_028" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_028_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_028_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_029" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_029_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_029_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_030" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_030_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_030_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_031" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  column "user_id" {
    null     = false
    type     = bigint
    unsigned = true
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
  index "idx_dm_posts_031_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_031_created_at" {
    columns = [column.created_at]
  }
}
