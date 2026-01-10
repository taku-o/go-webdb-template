// 開発環境用Atlas設定ファイル (MySQL)

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master-mysql"
  url = "mysql://webdb:webdb@localhost:3306/webdb_master"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/master-mysql"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1-mysql"
  url = "mysql://webdb:webdb@localhost:3307/webdb_sharding_1"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_1-mysql"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding_2-mysql"
  url = "mysql://webdb:webdb@localhost:3308/webdb_sharding_2"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_2-mysql"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding_3-mysql"
  url = "mysql://webdb:webdb@localhost:3309/webdb_sharding_3"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_3-mysql"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding_4-mysql"
  url = "mysql://webdb:webdb@localhost:3310/webdb_sharding_4"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_4-mysql"
  }
}
