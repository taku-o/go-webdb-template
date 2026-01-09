// dm_posts テーブル（sharding_db_3.db用: dm_posts_016 〜 dm_posts_023）

table "dm_posts_016" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_016_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_016_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_017" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_017_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_017_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_018" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_018_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_018_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_019" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_019_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_019_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_020" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_020_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_020_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_021" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_021_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_021_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_022" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_022_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_022_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_023" {
  schema = schema.public
  column "id" {
    null = false
    type = varchar(32)
  }
  column "user_id" {
    null = false
    type = varchar(32)
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
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_dm_posts_023_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_023_created_at" {
    columns = [column.created_at]
  }
}
