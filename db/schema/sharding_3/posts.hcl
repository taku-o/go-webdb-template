// Posts テーブル（sharding_db_3.db用: posts_016 〜 posts_023）

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
  index "idx_posts_023_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_023_created_at" {
    columns = [column.created_at]
  }
}
