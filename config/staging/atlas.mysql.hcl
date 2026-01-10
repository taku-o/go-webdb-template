// ステージング環境用Atlas設定ファイル (MySQL)
// 本番環境と同様にMySQLを想定

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  // ステージング環境のデータベースURL
  // 環境変数から読み込む場合: url = getenv("ATLAS_MASTER_DB_URL")
  url = "mysql://user:password@localhost:3306/master_db_staging"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/master-mysql"
  }
}

// シャーディングデータベース用環境

env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "mysql://user:password@localhost:3307/sharding_db_1_staging"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_1-mysql"
  }
}

env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "mysql://user:password@localhost:3308/sharding_db_2_staging"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_2-mysql"
  }
}

env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "mysql://user:password@localhost:3309/sharding_db_3_staging"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_3-mysql"
  }
}

env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "mysql://user:password@localhost:3310/sharding_db_4_staging"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_4-mysql"
  }
}
