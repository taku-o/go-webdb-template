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

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "sqlite://server/data/sharding_db_1.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_1"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "sqlite://server/data/sharding_db_2.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_2"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "sqlite://server/data/sharding_db_3.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_3"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "sqlite://server/data/sharding_db_4.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_4"
  }
}
