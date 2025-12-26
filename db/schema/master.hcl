// Master Database Schema
// マスターデータベースのスキーマ定義（newsテーブル、GoAdmin関連テーブル）

schema "main" {
}

// News テーブル
table "news" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "author_id" {
    null = true
    type = integer
  }
  column "published_at" {
    null = true
    type = datetime
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
  index "idx_news_published_at" {
    columns = [column.published_at]
  }
  index "idx_news_author_id" {
    columns = [column.author_id]
  }
}

// GoAdmin メニューテーブル
table "goadmin_menu" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "parent_id" {
    null    = false
    type    = integer
    default = 0
  }
  column "type" {
    null    = false
    type    = integer
    default = 0
  }
  column "order" {
    null    = false
    type    = integer
    default = 0
  }
  column "title" {
    null = false
    type = text
  }
  column "icon" {
    null = false
    type = text
  }
  column "uri" {
    null    = false
    type    = text
    default = ""
  }
  column "header" {
    null = true
    type = text
  }
  column "plugin_name" {
    null    = false
    type    = text
    default = ""
  }
  column "uuid" {
    null = true
    type = text
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
}

// GoAdmin 操作ログテーブル
table "goadmin_operation_log" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "path" {
    null = false
    type = text
  }
  column "method" {
    null = false
    type = text
  }
  column "ip" {
    null = false
    type = text
  }
  column "input" {
    null = false
    type = text
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_goadmin_operation_log_user_id" {
    columns = [column.user_id]
  }
}

// GoAdmin サイト設定テーブル
table "goadmin_site" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "key" {
    null = true
    type = text
  }
  column "value" {
    null = true
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "state" {
    null    = false
    type    = integer
    default = 0
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
}

// GoAdmin 権限テーブル
table "goadmin_permissions" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "name" {
    null = false
    type = text
  }
  column "slug" {
    null = false
    type = text
  }
  column "http_method" {
    null = true
    type = text
  }
  column "http_path" {
    null = false
    type = text
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_goadmin_permissions_slug" {
    unique  = true
    columns = [column.slug]
  }
}

// GoAdmin ロールテーブル
table "goadmin_roles" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "name" {
    null = false
    type = text
  }
  column "slug" {
    null = false
    type = text
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_goadmin_roles_slug" {
    unique  = true
    columns = [column.slug]
  }
}

// GoAdmin ロール-メニュー関連テーブル
table "goadmin_role_menu" {
  schema = schema.main
  column "role_id" {
    null = false
    type = integer
  }
  column "menu_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.role_id, column.menu_id]
  }
}

// GoAdmin ロール-権限関連テーブル
table "goadmin_role_permissions" {
  schema = schema.main
  column "role_id" {
    null = false
    type = integer
  }
  column "permission_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.role_id, column.permission_id]
  }
}

// GoAdmin ロール-ユーザー関連テーブル
table "goadmin_role_users" {
  schema = schema.main
  column "role_id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.role_id, column.user_id]
  }
}

// GoAdmin セッションテーブル
table "goadmin_session" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "sid" {
    null = false
    type = text
  }
  column "values" {
    null = false
    type = text
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
}

// GoAdmin ユーザー-権限関連テーブル
table "goadmin_user_permissions" {
  schema = schema.main
  column "user_id" {
    null = false
    type = integer
  }
  column "permission_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.user_id, column.permission_id]
  }
}

// GoAdmin 管理者ユーザーテーブル
table "goadmin_users" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "username" {
    null = false
    type = text
  }
  column "password" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "avatar" {
    null = true
    type = text
  }
  column "remember_token" {
    null = true
    type = text
  }
  column "created_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_goadmin_users_username" {
    unique  = true
    columns = [column.username]
  }
}

