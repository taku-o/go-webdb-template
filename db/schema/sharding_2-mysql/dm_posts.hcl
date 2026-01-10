// dm_posts テーブル（MySQL用: dm_posts_008 〜 dm_posts_015）

table "dm_posts_008" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_008_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_008_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_009" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_009_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_009_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_010" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_010_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_010_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_011" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_011_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_011_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_012" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_012_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_012_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_013" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_013_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_013_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_014" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_014_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_014_created_at" {
    columns = [column.created_at]
  }
}

table "dm_posts_015" {
  schema = schema.webdb_sharding_2
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
  index "idx_dm_posts_015_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_015_created_at" {
    columns = [column.created_at]
  }
}
