// テスト環境用Atlas設定ファイル (PostgreSQL)
// テスト用データベースは事前に作成し、マイグレーションを実行しておく必要がある
// マイグレーション: ./scripts/migrate-test.sh

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  url = "postgres://webdb:webdb@localhost:5432/webdb_master_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "postgres://webdb:webdb@localhost:5433/webdb_sharding_1_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_1"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "postgres://webdb:webdb@localhost:5434/webdb_sharding_2_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_2"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "postgres://webdb:webdb@localhost:5435/webdb_sharding_3_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_3"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "postgres://webdb:webdb@localhost:5436/webdb_sharding_4_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_4"
  }
}
