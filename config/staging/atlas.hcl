// ステージング環境用Atlas設定ファイル
// 本番環境と同様にPostgreSQL/MySQLを想定

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  // ステージング環境のデータベースURL
  // 環境変数から読み込む場合: url = getenv("ATLAS_MASTER_DB_URL")
  url = "postgres://user:password@localhost:5432/master_db_staging?sslmode=disable"
  dev = "postgres://user:password@localhost:5432/master_db_staging_dev?sslmode=disable"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングデータベース用環境

env "sharding_1" {
  src = "file://db/schema/sharding.hcl"
  url = "postgres://user:password@localhost:5432/sharding_db_1_staging?sslmode=disable"
  dev = "postgres://user:password@localhost:5432/sharding_db_1_staging_dev?sslmode=disable"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

env "sharding_2" {
  src = "file://db/schema/sharding.hcl"
  url = "postgres://user:password@localhost:5432/sharding_db_2_staging?sslmode=disable"
  dev = "postgres://user:password@localhost:5432/sharding_db_2_staging_dev?sslmode=disable"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

env "sharding_3" {
  src = "file://db/schema/sharding.hcl"
  url = "postgres://user:password@localhost:5432/sharding_db_3_staging?sslmode=disable"
  dev = "postgres://user:password@localhost:5432/sharding_db_3_staging_dev?sslmode=disable"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

env "sharding_4" {
  src = "file://db/schema/sharding.hcl"
  url = "postgres://user:password@localhost:5432/sharding_db_4_staging?sslmode=disable"
  dev = "postgres://user:password@localhost:5432/sharding_db_4_staging_dev?sslmode=disable"

  migration {
    dir = "file://db/migrations/sharding"
  }
}
