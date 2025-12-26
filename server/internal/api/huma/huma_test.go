package humaapi

import (
	"reflect"
	"testing"
)

// TestCreateUserInput はCreateUserInputの構造を確認
func TestCreateUserInput(t *testing.T) {
	input := CreateUserInput{}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(input.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}

	// Bodyのフィールドを確認
	nameField, ok := bodyType.FieldByName("Name")
	if !ok {
		t.Error("Body should have Name field")
	}
	if nameField.Tag.Get("json") != "name" {
		t.Error("Name should have json:\"name\" tag")
	}

	emailField, ok := bodyType.FieldByName("Email")
	if !ok {
		t.Error("Body should have Email field")
	}
	if emailField.Tag.Get("json") != "email" {
		t.Error("Email should have json:\"email\" tag")
	}
}

// TestGetUserInput はGetUserInputの構造を確認
func TestGetUserInput(t *testing.T) {
	input := GetUserInput{}

	// IDフィールドの存在を確認
	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("GetUserInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}
}

// TestListUsersInput はListUsersInputの構造を確認
func TestListUsersInput(t *testing.T) {
	input := ListUsersInput{}

	inputType := reflect.TypeOf(input)
	limitField, ok := inputType.FieldByName("Limit")
	if !ok {
		t.Error("ListUsersInput should have Limit field")
	}
	if limitField.Tag.Get("query") != "limit" {
		t.Error("Limit should have query:\"limit\" tag")
	}

	offsetField, ok := inputType.FieldByName("Offset")
	if !ok {
		t.Error("ListUsersInput should have Offset field")
	}
	if offsetField.Tag.Get("query") != "offset" {
		t.Error("Offset should have query:\"offset\" tag")
	}
}

// TestUpdateUserInput はUpdateUserInputの構造を確認
func TestUpdateUserInput(t *testing.T) {
	input := UpdateUserInput{}

	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("UpdateUserInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(input.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}
}

// TestDeleteUserInput はDeleteUserInputの構造を確認
func TestDeleteUserInput(t *testing.T) {
	input := DeleteUserInput{}

	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("DeleteUserInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}
}

// TestUserOutput はUserOutputの構造を確認
func TestUserOutput(t *testing.T) {
	output := UserOutput{}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}
}

// TestUsersOutput はUsersOutputの構造を確認
func TestUsersOutput(t *testing.T) {
	output := UsersOutput{}

	// Bodyがスライス型であることを確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Slice {
		t.Error("Body should be a slice")
	}
}

// TestDeleteUserOutput はDeleteUserOutputの構造を確認
func TestDeleteUserOutput(t *testing.T) {
	_ = DeleteUserOutput{}
	// 204 No Contentなのでフィールドは不要
}

// TestCreatePostInput はCreatePostInputの構造を確認
func TestCreatePostInput(t *testing.T) {
	input := CreatePostInput{}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(input.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}

	// Bodyのフィールドを確認
	userIDField, ok := bodyType.FieldByName("UserID")
	if !ok {
		t.Error("Body should have UserID field")
	}
	if userIDField.Tag.Get("json") != "user_id" {
		t.Error("UserID should have json:\"user_id\" tag")
	}

	titleField, ok := bodyType.FieldByName("Title")
	if !ok {
		t.Error("Body should have Title field")
	}
	if titleField.Tag.Get("json") != "title" {
		t.Error("Title should have json:\"title\" tag")
	}
}

// TestGetPostInput はGetPostInputの構造を確認
func TestGetPostInput(t *testing.T) {
	input := GetPostInput{}

	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("GetPostInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}

	userIDField, ok := inputType.FieldByName("UserID")
	if !ok {
		t.Error("GetPostInput should have UserID field")
	}
	if userIDField.Tag.Get("query") != "user_id" {
		t.Error("UserID should have query:\"user_id\" tag")
	}
}

// TestListPostsInput はListPostsInputの構造を確認
func TestListPostsInput(t *testing.T) {
	input := ListPostsInput{}

	inputType := reflect.TypeOf(input)
	limitField, ok := inputType.FieldByName("Limit")
	if !ok {
		t.Error("ListPostsInput should have Limit field")
	}
	if limitField.Tag.Get("query") != "limit" {
		t.Error("Limit should have query:\"limit\" tag")
	}

	userIDField, ok := inputType.FieldByName("UserID")
	if !ok {
		t.Error("ListPostsInput should have UserID field")
	}
	if userIDField.Tag.Get("query") != "user_id" {
		t.Error("UserID should have query:\"user_id\" tag")
	}
}

// TestGetUserPostsInput はGetUserPostsInputの構造を確認
func TestGetUserPostsInput(t *testing.T) {
	input := GetUserPostsInput{}

	inputType := reflect.TypeOf(input)
	limitField, ok := inputType.FieldByName("Limit")
	if !ok {
		t.Error("GetUserPostsInput should have Limit field")
	}
	if limitField.Tag.Get("query") != "limit" {
		t.Error("Limit should have query:\"limit\" tag")
	}
}

// TestPostOutput はPostOutputの構造を確認
func TestPostOutput(t *testing.T) {
	output := PostOutput{}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}
}

// TestPostsOutput はPostsOutputの構造を確認
func TestPostsOutput(t *testing.T) {
	output := PostsOutput{}

	// Bodyがスライス型であることを確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Slice {
		t.Error("Body should be a slice")
	}
}

// TestUserPostsOutput はUserPostsOutputの構造を確認
func TestUserPostsOutput(t *testing.T) {
	output := UserPostsOutput{}

	// Bodyがスライス型であることを確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Slice {
		t.Error("Body is not a slice")
	}
}

// TestDeletePostOutput はDeletePostOutputの構造を確認
func TestDeletePostOutput(t *testing.T) {
	_ = DeletePostOutput{}
	// 204 No Contentなのでフィールドは不要
}
