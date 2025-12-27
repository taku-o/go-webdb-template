// Users テーブル（sharding_db_2.db用: users_008 〜 users_015）

table "users_008" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_008_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_009" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_009_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_010" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_010_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_011" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_011_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_012" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_012_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_013" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_013_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_014" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_014_email" {
    unique  = true
    columns = [column.email]
  }
}

table "users_015" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
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
  index "idx_users_015_email" {
    unique  = true
    columns = [column.email]
  }
}
