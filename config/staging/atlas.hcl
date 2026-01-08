// ステージング環境用Atlas設定ファイル
// 本番環境と同様にPostgreSQL/MySQLを想定

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  // ステージング環境のデータベースURL
  // 環境変数から読み込む場合: url = getenv("ATLAS_MASTER_DB_URL")
  url = "postgres://user:password@localhost:5432/master_db_staging?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングデータベース用環境

env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "postgres://user:password@localhost:5432/sharding_db_1_staging?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_1"
  }
}

env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "postgres://user:password@localhost:5432/sharding_db_2_staging?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_2"
  }
}

env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "postgres://user:password@localhost:5432/sharding_db_3_staging?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_3"
  }
}

env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "postgres://user:password@localhost:5432/sharding_db_4_staging?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_4"
  }
}
