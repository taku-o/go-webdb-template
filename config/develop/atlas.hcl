// 開発環境用Atlas設定ファイル

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  url = "sqlite://server/data/master.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングデータベース用環境（共通設定）
// 各シャーディングDBに対して個別に適用する

env "sharding" {
  src = "file://db/schema/sharding.hcl"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding.hcl"
  url = "sqlite://server/data/sharding_db_1.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding.hcl"
  url = "sqlite://server/data/sharding_db_2.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding.hcl"
  url = "sqlite://server/data/sharding_db_3.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding.hcl"
  url = "sqlite://server/data/sharding_db_4.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding"
  }
}
