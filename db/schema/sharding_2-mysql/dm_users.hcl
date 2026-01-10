// dm_users テーブル（MySQL用: dm_users_008 〜 dm_users_015）

table "dm_users_008" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_008_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_009" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_009_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_010" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_010_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_011" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_011_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_012" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_012_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_013" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_013_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_014" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_014_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_015" {
  schema = schema.webdb_sharding_2
  column "id" {
    null = false
    type = varchar(32)
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = varchar(191)
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
  index "idx_dm_users_015_email" {
    unique  = true
    columns = [column.email]
  }
}
