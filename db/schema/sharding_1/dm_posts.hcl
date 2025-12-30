// dm_posts テーブル（sharding_db_1.db用: dm_posts_000 〜 dm_posts_007）

table "dm_posts_000" {
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
  index "idx_dm_posts_000_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_000_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_001" {
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
  index "idx_dm_posts_001_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_001_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_002" {
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
  index "idx_dm_posts_002_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_002_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_003" {
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
  index "idx_dm_posts_003_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_003_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_004" {
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
  index "idx_dm_posts_004_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_004_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_005" {
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
  index "idx_dm_posts_005_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_005_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_006" {
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
  index "idx_dm_posts_006_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_006_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_007" {
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
  index "idx_dm_posts_007_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_007_created_at" {
    columns = [column.created_at]
  }
}
