server/cmd/server-status/main.go の設計の修正

server/cmd/server-status/main.go を
usecase、serviceを利用して処理を行う仕組みに作り替える。

server/cmd/server-status/main.go
    入出力制御を担当、usecaseを呼び出す、usecaseからリストを受け取ってコンソール出力する
↓
server/internal/usecase/cli/server_status_usecase.go
    serviceに渡すパラメータを作る
    serviceから受け取ったリストをmain.goに渡す
↓
server/internal/service/server_status_service.go 
    list({name, host, port}) を受け取って、
    結果のリストを返す


/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/161
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0079-status-usecaseとしてください。"
think.

要件定義書を承認します。

/kiro:spec-design 0079-status-usecase

server/internal/usecase/cli/server_status_usecase.go
でserviceに渡すパラメータを作る。
サーバーの一覧は、usecaseで定義する。

設計書を承認します。

/kiro:spec-tasks 0079-status-usecase

このロジックって、並び順変わっちゃわない？
func (s *ServerStatusService) checkAllServers(servers []ServerInfo, timeout time.Duration) []ServerStatus {
	var wg sync.WaitGroup
	results := make([]ServerStatus, len(servers))

	for i, server := range servers {
		wg.Add(1)
		go func(index int, srv ServerInfo) {
			defer wg.Done()
			results[index] = s.checkServerStatus(srv, timeout)
		}(i, server)
	}

	wg.Wait()
	return results
}

設計書を承認します。


APIサーバー、
Adminサーバー、
CLIコード
どれも
controller -> usecase -> service -> repository -> DB
という作りなっているんだけど、
cc-sddのsteeringドキュメントあたりに、これらの機能を作るときには、
こういう設計をする、と記載して、今後の設計を制御することはできますか？


/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0079-status-usecase

タスクリストを承認します。

/kiro:spec-impl 0079-status-usecase

リファクタリングしたなら、
リファクタリングしたコードに合わせたテストに直して。
server/cmd/server-status/main_test.go

実装の削除は駄目だよ。

stagingに上がっている修正をcommitして、
https://github.com/taku-o/go-webdb-template/issues/161 に
対してpull requestを作成してください。

/review 162





