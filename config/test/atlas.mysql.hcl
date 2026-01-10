// テスト環境用Atlas設定ファイル (MySQL)
// テスト用データベースは事前に作成し、マイグレーションを実行しておく必要がある
// マイグレーション: ./scripts/migrate-test-mysql.sh

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  url = "mysql://webdb:webdb@localhost:3306/webdb_master_test"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/master-mysql"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "mysql://webdb:webdb@localhost:3307/webdb_sharding_1_test"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_1-mysql"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "mysql://webdb:webdb@localhost:3308/webdb_sharding_2_test"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_2-mysql"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "mysql://webdb:webdb@localhost:3309/webdb_sharding_3_test"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_3-mysql"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "mysql://webdb:webdb@localhost:3310/webdb_sharding_4_test"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_4-mysql"
  }
}
