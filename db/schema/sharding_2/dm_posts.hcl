// dm_posts テーブル（sharding_db_2.db用: dm_posts_008 〜 dm_posts_015）

table "dm_posts_008" {
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
  index "idx_dm_posts_008_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_008_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_009" {
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
  index "idx_dm_posts_009_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_009_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_010" {
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
  index "idx_dm_posts_010_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_010_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_011" {
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
  index "idx_dm_posts_011_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_011_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_012" {
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
  index "idx_dm_posts_012_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_012_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_013" {
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
  index "idx_dm_posts_013_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_013_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_014" {
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
  index "idx_dm_posts_014_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_014_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_015" {
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
  index "idx_dm_posts_015_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_015_created_at" {
    columns = [column.created_at]
  }
}
