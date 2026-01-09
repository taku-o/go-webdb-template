/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/85 のsub issue
https://github.com/taku-o/go-webdb-template/issues/90 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0047-pgmain-dockerとしてください。

issue 90の修正は、最終的に switch-to-postgresqlブランチに修正を取り込みます。"
think.

要件定義書を承認します。

/kiro:spec-design 0047-pgmain-docker

設計書を承認します。

/kiro:spec-tasks 0047-pgmain-docker

adminサーバーはhealthが今のところ無いので、
何かを代替にしてください。
> http://localhost:8081/health

タスクリストを承認します。

この要件の作業用のgitブランチをswitch-to-postgresqlブランチから切ってください。
ここまでの作業をcommitしてください。
そこまで作業したら、いったんユーザーに応答を返してください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0047-pgmain-docker

Docker Desktop再起動する。
少し待ってて。

お待たせしました。
Docker Desktop再起動した。



これはMySQLの例だけど、docker-composeにこんな感じに指定して。
> ⏺ 設定ファイルではローカル開発用にlocalhostを使用しています。Dockerコンテナ間の通信は設定の問題であり、Dockerfileの変更タスクのスコープ外です。

services:
  api:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: api-develop
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=develop
      - REDIS_JOBQUEUE_ADDR=redis:6379
      # Database DSNs for Docker environment
      - DB_MASTER_WRITER_DSN=webdb:webdb@tcp(mysql-master:3306)/master_db?parseTime=true
      - DB_MASTER_READER_DSN=webdb:webdb@tcp(mysql-master:3306)/master_db?parseTime=true
      - DB_SHARD1_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-1:3306)/sharding_db_1?parseTime=true
      - DB_SHARD1_READER_DSN=webdb:webdb@tcp(mysql-sharding-1:3306)/sharding_db_1?parseTime=true
      - DB_SHARD2_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-1:3306)/sharding_db_1?parseTime=true
      - DB_SHARD2_READER_DSN=webdb:webdb@tcp(mysql-sharding-1:3306)/sharding_db_1?parseTime=true
      - DB_SHARD3_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-2:3306)/sharding_db_2?parseTime=true
      - DB_SHARD3_READER_DSN=webdb:webdb@tcp(mysql-sharding-2:3306)/sharding_db_2?parseTime=true
      - DB_SHARD4_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-2:3306)/sharding_db_2?parseTime=true
      - DB_SHARD4_READER_DSN=webdb:webdb@tcp(mysql-sharding-2:3306)/sharding_db_2?parseTime=true
      - DB_SHARD5_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-3:3306)/sharding_db_3?parseTime=true
      - DB_SHARD5_READER_DSN=webdb:webdb@tcp(mysql-sharding-3:3306)/sharding_db_3?parseTime=true
      - DB_SHARD6_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-3:3306)/sharding_db_3?parseTime=true
      - DB_SHARD6_READER_DSN=webdb:webdb@tcp(mysql-sharding-3:3306)/sharding_db_3?parseTime=true
      - DB_SHARD7_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-4:3306)/sharding_db_4?parseTime=true
      - DB_SHARD7_READER_DSN=webdb:webdb@tcp(mysql-sharding-4:3306)/sharding_db_4?parseTime=true
      - DB_SHARD8_WRITER_DSN=webdb:webdb@tcp(mysql-sharding-4:3306)/sharding_db_4?parseTime=true
      - DB_SHARD8_READER_DSN=webdb:webdb@tcp(mysql-sharding-4:3306)/sharding_db_4?parseTime=true
    volumes:
      - ./config/develop:/app/config/develop:ro
      - ./logs:/app/logs
    networks:
      - mysql-network
      - redis-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  mysql-network:
    external: true
    name: mysql-network
  redis-network:
    external: true
    name: redis-network

これらにもDB_MASTER_WRITER_DSNあたりの設定を入れておいて。
動作確認まではしなくてもいい。
* docker-compose.admin.production.yml
* docker-compose.admin.staging.yml
* docker-compose.api.production.yml
* docker-compose.api.staging.yml


server/internal/config/config.goに
こんな感じの修正を入れる。

	// Docker環境用: 環境変数でデータベースDSNを上書き
	// Master
	if len(cfg.Database.Groups.Master) > 0 {
		if writerDSN := os.Getenv("DB_MASTER_WRITER_DSN"); writerDSN != "" {
			cfg.Database.Groups.Master[0].WriterDSN = writerDSN
		}
		if readerDSN := os.Getenv("DB_MASTER_READER_DSN"); readerDSN != "" {
			cfg.Database.Groups.Master[0].ReaderDSNs = []string{readerDSN}
		}
		// GoAdmin用: Host/Portの上書き
		if host := os.Getenv("DB_MASTER_HOST"); host != "" {
			cfg.Database.Groups.Master[0].Host = host
		}
		if portStr := os.Getenv("DB_MASTER_PORT"); portStr != "" {
			if port, err := strconv.Atoi(portStr); err != nil {
				log.Printf("Warning: DB_MASTER_PORT '%s' is not a valid integer, using default value", portStr)
			} else {
				cfg.Database.Groups.Master[0].Port = port
			}
		}
	}
	// Sharding databases
	for i := range cfg.Database.Groups.Sharding.Databases {
		envKeyWriter := fmt.Sprintf("DB_SHARD%d_WRITER_DSN", i+1)
		envKeyReader := fmt.Sprintf("DB_SHARD%d_READER_DSN", i+1)
		if writerDSN := os.Getenv(envKeyWriter); writerDSN != "" {
			cfg.Database.Groups.Sharding.Databases[i].WriterDSN = writerDSN
		}
		if readerDSN := os.Getenv(envKeyReader); readerDSN != "" {
			cfg.Database.Groups.Sharding.Databases[i].ReaderDSNs = []string{readerDSN}
		}
	}


タスクは全部終わってる？
タスク外と判断してスキップした作業とかある？

この作業お願い。
>  - ヘルスチェック動作確認

APIサーバー、クライアントサーバー、Adminサーバーを
Docker版で起動して

クライアントサーバーを止めて、Docker版を起動して

docker-compose.admin.develop.yml に入れた修正を
docker-compose.admin.staging.yml
docker-compose.admin.production.yml
にも入れて。

PostgreSQLサーバー、APIサーバー、クライアントサーバー、Adminサーバーを
Docker版で起動して

OK。
PostgreSQLサーバー、APIサーバー、クライアントサーバー、Adminサーバーを停止して。

cloudbeaverも止めちゃおう

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/90 に対して
pull requestを作成してください。

stagingされている修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/90 に対して
pull requestを作成してください。






