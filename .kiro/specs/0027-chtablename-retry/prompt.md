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
progress.mdに進捗を記録しながら､作業を進めてください。


発見ありがとう。これは未着手？
>  - 旧命名規則残存確認: 追加で14件のテスト関数名修正漏れを発見・修正
>    - dm_post_repository_test.go: 6件
>    - dm_post_repository_gorm_test.go: 8件


新しいファイルで、
	server/generate-sample-data
	server/main
というのが増えているけど、これは不要ファイル？

削除してください。


APIサーバー、クライアントサーバー、GoAdminサーバーを
再起動してください。


GoAdmin管理画面で、ソースマップ読込エラーが出ている。
JavaScriptをビルドとかした？
[Error] Failed to load resource: the server responded with a status of 404 (Not Found) (bootstrap-select.min.js.map, line 0)
[Error] Failed to load resource: the server responded with a status of 404 (Not Found) (bootstrap-select.min.js.map, line 0)
[Error] Failed to load resource: the server responded with a status of 404 (Not Found) (wangEditor.min.js.map, line 0)
[Error] Failed to load resource: the server responded with a status of 404 (Not Found) (toastr.js.map, line 0)


記憶にない。が、放っておこう。


この辺りも直したい。

server/internal/admin/tables.go:func GetNewsTable(ctx *context.Context) table.Table {
server/internal/admin/tables.go:	"dm-news": GetNewsTable,
server/internal/model/dm_news.go:type CreateNewsRequest struct {
server/internal/model/dm_news.go:type UpdateNewsRequest struct {

server/internal/repository/dm_post_repository.go:func (r *DmPostRepository) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
server/internal/repository/dm_post_repository.go:	userPosts := make([]*model.DmUserPost, 0)
server/internal/repository/dm_post_repository.go:			userPosts = append(userPosts, &up)
server/internal/repository/dm_post_repository.go:	return userPosts, nil
server/internal/api/handler/dm_post_handler.go:		resp := &humaapi.UserPostsOutput{}
server/internal/api/huma/inputs.go:type CreatePostInput struct {
server/internal/api/huma/inputs.go:type GetPostInput struct {
server/internal/api/huma/inputs.go:type ListPostsInput struct {
server/internal/api/huma/inputs.go:type UpdatePostInput struct {
server/internal/api/huma/inputs.go:type DeletePostInput struct {
server/internal/api/huma/inputs.go:type GetUserPostsInput struct {
server/internal/api/huma/outputs.go:type PostOutput struct {
server/internal/api/huma/outputs.go:type PostsOutput struct {
server/internal/api/huma/outputs.go:type UserPostsOutput struct {
server/internal/api/huma/outputs.go:type DeletePostOutput struct {

server/internal/api/huma/huma_test.go:func TestCreatePostInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := CreatePostInput{}
server/internal/api/huma/huma_test.go:func TestGetPostInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := GetPostInput{}
server/internal/api/huma/huma_test.go:func TestListPostsInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := ListPostsInput{}
server/internal/api/huma/huma_test.go:func TestGetUserPostsInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := GetUserPostsInput{}
server/internal/api/huma/huma_test.go:func TestPostOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	output := PostOutput{}
server/internal/api/huma/huma_test.go:func TestPostsOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	output := PostsOutput{}
server/internal/api/huma/huma_test.go:func TestUserPostsOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	output := UserPostsOutput{}
server/internal/api/huma/huma_test.go:func TestDeletePostOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	_ = DeletePostOutput{}
server/internal/api/router/router_test.go:func TestRegisterPostEndpointsIntegration(t *testing.T) {

server/internal/api/router/router.go:func NewRouter(userHandler *handler.DmUserHandler, postHandler *handler.DmPostHandler, todayHandler *handler.TodayHandler, cfg *config.Config) *echo.Echo {
server/internal/api/router/router.go:	handler.RegisterPostEndpoints(humaAPI, postHandler)

server/internal/api/handler/dm_user_handler_huma_test.go:func TestRegisterUserEndpointsExists(t *testing.T) {
server/internal/api/handler/dm_user_handler_huma_test.go:	var _ func(api huma.API, h *DmUserHandler) = RegisterUserEndpoints

server/internal/api/handler/dm_user_handler.go:func RegisterUserEndpoints(api huma.API, h *DmUserHandler) {
server/internal/api/handler/dm_user_handler.go:	}, func(ctx context.Context, input *humaapi.CreateUserInput) (*humaapi.UserOutput, error) {
server/internal/api/handler/dm_user_handler.go:		resp := &humaapi.UserOutput{}
server/internal/api/handler/dm_user_handler.go:	}, func(ctx context.Context, input *humaapi.GetUserInput) (*humaapi.UserOutput, error) {
server/internal/api/handler/dm_user_handler.go:		resp := &humaapi.UserOutput{}
server/internal/api/handler/dm_user_handler.go:	}, func(ctx context.Context, input *humaapi.ListUsersInput) (*humaapi.UsersOutput, error) {
server/internal/api/handler/dm_user_handler.go:		resp := &humaapi.UsersOutput{}
server/internal/api/handler/dm_user_handler.go:	}, func(ctx context.Context, input *humaapi.UpdateUserInput) (*humaapi.UserOutput, error) {
server/internal/api/handler/dm_user_handler.go:		resp := &humaapi.UserOutput{}

server/internal/api/handler/dm_user_handler.go:	}, func(ctx context.Context, input *humaapi.DeleteUserInput) (*struct{}, error) {
server/internal/api/handler/dm_post_handler.go:	}, func(ctx context.Context, input *humaapi.GetUserPostsInput) (*humaapi.UserPostsOutput, error) {
server/internal/api/handler/dm_post_handler.go:		resp := &humaapi.UserPostsOutput{}

server/internal/api/huma/inputs.go:type CreateUserInput struct {
server/internal/api/huma/inputs.go:type GetUserInput struct {
server/internal/api/huma/inputs.go:type ListUsersInput struct {
server/internal/api/huma/inputs.go:type UpdateUserInput struct {
server/internal/api/huma/inputs.go:type DeleteUserInput struct {

server/internal/api/huma/inputs.go:type GetUserPostsInput struct {
server/internal/api/huma/outputs.go:type UserOutput struct {
server/internal/api/huma/outputs.go:type UsersOutput struct {
server/internal/api/huma/outputs.go:type DeleteUserOutput struct {
server/internal/api/huma/outputs.go:type UserPostsOutput struct {
server/internal/api/huma/huma_test.go:func TestCreateUserInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := CreateUserInput{}

server/internal/api/huma/huma_test.go:func TestGetUserInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := GetUserInput{}
server/internal/api/huma/huma_test.go:func TestListUsersInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := ListUsersInput{}
server/internal/api/huma/huma_test.go:func TestUpdateUserInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := UpdateUserInput{}
server/internal/api/huma/huma_test.go:func TestDeleteUserInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := DeleteUserInput{}
server/internal/api/huma/huma_test.go:func TestUserOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	output := UserOutput{}
server/internal/api/huma/huma_test.go:func TestUsersOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	output := UsersOutput{}

server/internal/api/huma/huma_test.go:func TestDeleteUserOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	_ = DeleteUserOutput{}

server/internal/api/huma/huma_test.go:func TestGetUserPostsInput(t *testing.T) {
server/internal/api/huma/huma_test.go:	input := GetUserPostsInput{}
server/internal/api/huma/huma_test.go:func TestUserPostsOutput(t *testing.T) {
server/internal/api/huma/huma_test.go:	output := UserPostsOutput{}
server/internal/api/router/router_test.go:func TestRegisterUserEndpointsIntegration(t *testing.T) {
server/internal/api/router/router_test.go:		// handler.RegisterUserEndpoints(api, h) の形式で呼び出し可能
server/internal/api/router/router.go:func NewRouter(userHandler *handler.DmUserHandler, postHandler *handler.DmPostHandler, todayHandler *handler.TodayHandler, cfg *config.Config) *echo.Echo {
server/internal/api/router/router.go:	handler.RegisterUserEndpoints(humaAPI, userHandler)
server/internal/service/dm_post_service.go:func (s *DmPostService) ListDmPostsByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.DmPost, error) {

think.


その修正方針であっています。
修正が多いから、タスクの一覧をprogress.mdあたりに出力してから、
作業を進めてくれるかな？

APIサーバー、クライアントサーバー、GoAdminサーバーを
再起動してください。

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/54 に対して
pull requestを作成してください。








