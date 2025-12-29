package humaapi

import (
	"reflect"
	"testing"
)

// TestCreateDmUserInput はCreateDmUserInputの構造を確認
func TestCreateDmUserInput(t *testing.T) {
	input := CreateDmUserInput{}

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

// TestGetDmUserInput はGetDmUserInputの構造を確認
func TestGetDmUserInput(t *testing.T) {
	input := GetDmUserInput{}

	// IDフィールドの存在を確認
	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("GetDmUserInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}
}

// TestListDmUsersInput はListDmUsersInputの構造を確認
func TestListDmUsersInput(t *testing.T) {
	input := ListDmUsersInput{}

	inputType := reflect.TypeOf(input)
	limitField, ok := inputType.FieldByName("Limit")
	if !ok {
		t.Error("ListDmUsersInput should have Limit field")
	}
	if limitField.Tag.Get("query") != "limit" {
		t.Error("Limit should have query:\"limit\" tag")
	}

	offsetField, ok := inputType.FieldByName("Offset")
	if !ok {
		t.Error("ListDmUsersInput should have Offset field")
	}
	if offsetField.Tag.Get("query") != "offset" {
		t.Error("Offset should have query:\"offset\" tag")
	}
}

// TestUpdateDmUserInput はUpdateDmUserInputの構造を確認
func TestUpdateDmUserInput(t *testing.T) {
	input := UpdateDmUserInput{}

	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("UpdateDmUserInput should have ID field")
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

// TestDeleteDmUserInput はDeleteDmUserInputの構造を確認
func TestDeleteDmUserInput(t *testing.T) {
	input := DeleteDmUserInput{}

	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("DeleteDmUserInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}
}

// TestDmUserOutput はDmUserOutputの構造を確認
func TestDmUserOutput(t *testing.T) {
	output := DmUserOutput{}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}
}

// TestDmUsersOutput はDmUsersOutputの構造を確認
func TestDmUsersOutput(t *testing.T) {
	output := DmUsersOutput{}

	// Bodyがスライス型であることを確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Slice {
		t.Error("Body should be a slice")
	}
}

// TestDeleteDmUserOutput はDeleteDmUserOutputの構造を確認
func TestDeleteDmUserOutput(t *testing.T) {
	_ = DeleteDmUserOutput{}
	// 204 No Contentなのでフィールドは不要
}

// TestCreateDmPostInput はCreateDmPostInputの構造を確認
func TestCreateDmPostInput(t *testing.T) {
	input := CreateDmPostInput{}

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

// TestGetDmPostInput はGetDmPostInputの構造を確認
func TestGetDmPostInput(t *testing.T) {
	input := GetDmPostInput{}

	inputType := reflect.TypeOf(input)
	idField, ok := inputType.FieldByName("ID")
	if !ok {
		t.Error("GetDmPostInput should have ID field")
	}
	if idField.Tag.Get("path") != "id" {
		t.Error("ID should have path:\"id\" tag")
	}

	userIDField, ok := inputType.FieldByName("UserID")
	if !ok {
		t.Error("GetDmPostInput should have UserID field")
	}
	if userIDField.Tag.Get("query") != "user_id" {
		t.Error("UserID should have query:\"user_id\" tag")
	}
}

// TestListDmPostsInput はListDmPostsInputの構造を確認
func TestListDmPostsInput(t *testing.T) {
	input := ListDmPostsInput{}

	inputType := reflect.TypeOf(input)
	limitField, ok := inputType.FieldByName("Limit")
	if !ok {
		t.Error("ListDmPostsInput should have Limit field")
	}
	if limitField.Tag.Get("query") != "limit" {
		t.Error("Limit should have query:\"limit\" tag")
	}

	userIDField, ok := inputType.FieldByName("UserID")
	if !ok {
		t.Error("ListDmPostsInput should have UserID field")
	}
	if userIDField.Tag.Get("query") != "user_id" {
		t.Error("UserID should have query:\"user_id\" tag")
	}
}

// TestGetDmUserPostsInput はGetDmUserPostsInputの構造を確認
func TestGetDmUserPostsInput(t *testing.T) {
	input := GetDmUserPostsInput{}

	inputType := reflect.TypeOf(input)
	limitField, ok := inputType.FieldByName("Limit")
	if !ok {
		t.Error("GetDmUserPostsInput should have Limit field")
	}
	if limitField.Tag.Get("query") != "limit" {
		t.Error("Limit should have query:\"limit\" tag")
	}
}

// TestDmPostOutput はDmPostOutputの構造を確認
func TestDmPostOutput(t *testing.T) {
	output := DmPostOutput{}

	// Body構造体の存在を確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Struct {
		t.Error("Body should be a struct")
	}
}

// TestDmPostsOutput はDmPostsOutputの構造を確認
func TestDmPostsOutput(t *testing.T) {
	output := DmPostsOutput{}

	// Bodyがスライス型であることを確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Slice {
		t.Error("Body should be a slice")
	}
}

// TestDmUserPostsOutput はDmUserPostsOutputの構造を確認
func TestDmUserPostsOutput(t *testing.T) {
	output := DmUserPostsOutput{}

	// Bodyがスライス型であることを確認
	bodyType := reflect.TypeOf(output.Body)
	if bodyType.Kind() != reflect.Slice {
		t.Error("Body is not a slice")
	}
}

// TestDeleteDmPostOutput はDeleteDmPostOutputの構造を確認
func TestDeleteDmPostOutput(t *testing.T) {
	_ = DeleteDmPostOutput{}
	// 204 No Contentなのでフィールドは不要
}
