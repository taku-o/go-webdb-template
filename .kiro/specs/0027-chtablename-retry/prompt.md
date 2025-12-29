/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/54 に対応してください。
一度対応して修正が入っていますが、特に変数名の修正のなおし損ねなどで、
やり残しの作業が多数ありました。

server/cmd/generate-sample-data/main.go
    generateUsers -> generateDmUsers
    insertUsersBatch -> insertDmUsersBatch
    fetchUserIDs -> fetchDmUserIDs

server/cmd/admin/main.go
	pages.UserRegisterPage(goadminContext.NewContext(ctx.Request), groupManager) -> pages.DmUserRegisterPage(goadminContext.NewContext(ctx.Request), groupManager)
    pages.UserRegisterCompletePage(goadminContext.NewContext(ctx.Request), conn) -> pages.DmUserRegisterCompletePage(goadminContext.NewContext(ctx.Request), conn)

server/cmd/server/main.go
    userService := service.NewDmUserService(userRepo) -> dmUserService := service.NewDmUserService(userRepo)
    userRepo := repository.NewDmUserRepositoryGORM(groupManager) -> dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
    userHandler := handler.NewDmUserHandler(userService) -> dmUserHandler := handler.NewDmUserHandler(userService)

server/cmd/list-dm-users/main.go
    func printUsersTSV(dmUsers []*model.DmUser) -> func printDmUsersTSV(dmUsers []*model.DmUser)

server/internal/repository/dm_post_repository_gorm.go
	var tableUserPosts []*model.DmUserPost -> var tableDmUserPosts []*model.DmUserPost

cc-sddのfeature名は0027-chtablename-retryとしてください。"
think.


リストしたのは一例に過ぎなくて、
まだまだ他にも大量にある。
探して修正して。

server/cmd/generate-sample-data/main.go
server/cmd/admin/main.go
server/cmd/server/main.go
server/cmd/list-dm-users/main.go
server/cmd/list-dm-users/main_test.go
server/test/integration/dm_user_flow_gorm_test.go
server/test/integration/dm_user_flow_test.go
server/test/integration/dm_post_flow_test.go
server/test/integration/sharding_test.go
server/test/fixtures/dm_users.go
server/test/e2e/api_test.go
server/internal/repository/interfaces.go
server/internal/repository/dm_post_repository_gorm.go
server/internal/repository/dm_user_repository_gorm_test.go
server/internal/repository/dm_post_repository.go
server/internal/repository/dm_user_repository_test.go
server/internal/admin/sharding.go



任意ではない。そのままの名前だと邪魔だから、変更は必須です。
> 注意**: テスト関数名の変更は任意だが、一貫性のため推奨


userだけじゃなくて、postとnewsも対象だよ。
そのままの名前が残ると邪魔だからね。

良さそうだ。
要件定義書を承認します。

/kiro:spec-design 0027-chtablename-retry

設計書を承認します。

/kiro:spec-tasks 0027-chtablename-retry

タスクリストを承認します。

作業量が多く、おそらく作業の途中でコンテキストが尽きてしまうでしょう。
タスクの進捗を管理する方法が必要です。
.kiro/specs/0027-chtablename-retry/progress.md
を作成してください。


/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0027-chtablename-retry






