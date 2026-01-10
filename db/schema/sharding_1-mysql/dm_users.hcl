// dm_users テーブル（MySQL用: dm_users_000 〜 dm_users_007）

table "dm_users_000" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_000_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_001" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_001_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_002" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_002_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_003" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_003_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_004" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_004_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_005" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_005_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_006" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_006_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_007" {
  schema = schema.webdb_sharding_1
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
  index "idx_dm_users_007_email" {
    unique  = true
    columns = [column.email]
  }
}
