// dm_users テーブル（sharding_db_3.db用: dm_users_016 〜 dm_users_023）

table "dm_users_016" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_016_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_017" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_017_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_018" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_018_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_019" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_019_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_020" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_020_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_021" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_021_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_022" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_022_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_023" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
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
  index "idx_dm_users_023_email" {
    unique  = true
    columns = [column.email]
  }
}
