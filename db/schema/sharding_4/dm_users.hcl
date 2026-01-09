// dm_users テーブル（sharding_db_4.db用: dm_users_024 〜 dm_users_031）

table "dm_users_024" {
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
  index "idx_dm_users_024_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_025" {
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
  index "idx_dm_users_025_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_026" {
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
  index "idx_dm_users_026_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_027" {
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
  index "idx_dm_users_027_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_028" {
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
  index "idx_dm_users_028_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_029" {
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
  index "idx_dm_users_029_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_030" {
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
  index "idx_dm_users_030_email" {
    unique  = true
    columns = [column.email]
  }
}

table "dm_users_031" {
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
  index "idx_dm_users_031_email" {
    unique  = true
    columns = [column.email]
  }
}
