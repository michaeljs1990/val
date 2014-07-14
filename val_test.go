package val

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

// Make string into io.ReadCloser
// this is just for convinience
func jsonFactory(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

// Ensure all required fields are matching
func TestRequired(t *testing.T) {

	// Test if STRING required is valid
	var testString struct {
		Test string `json:"something" validate:"required" `
	}

	testJSON := jsonFactory(`{"something": "hello"}`)

	if err := Guaranty(&testString, testJSON); err != nil {
		t.Error(err)
	}

	var testString2 struct {
		Test string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	if err := Guaranty(&testString2, testJSON); err == nil {
		t.Error("Required string, empty JSON object should return error but did not.")
	}

	// Test if INT require is valid
	var testInt struct {
		Test int `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": 2}`)

	if err := Guaranty(&testInt, testJSON); err != nil {
		t.Error(err)
	}

	// Need to work on getting this to pass
	var testInt2 struct {
		Test int `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	if err := Guaranty(&testInt2, testJSON); err == nil {
		t.Error("Required int, empty JSON object should return error but did not.")
	}

	// Test if BOOL required is valid
	var testBool struct {
		Test bool `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": true}`)

	if err := Guaranty(&testBool, testJSON); err != nil {
		t.Error(err)
	}

	var testBool2 struct {
		Test string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	if err := Guaranty(&testBool2, testJSON); err == nil {
		t.Error("Required bool, empty JSON object should return error but did not.")
	}

	// Test if ARRAY required is valid
	var testArray struct {
		Test []string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": ["test", "data"]}`)

	if err := Guaranty(&testArray, testJSON); err != nil {
		t.Error(err)
	}

	var testArray2 struct {
		Test []string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	if err := Guaranty(&testArray2, testJSON); err == nil {
		t.Error("Required array, empty JSON object should return error but did not.")
	}

	// Test is OBJECT required is valid
	type testObjectTP struct {
		Name string `json:"name" validate:"required" `
	}

	var testObject struct {
		Test testObjectTP `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": {"name": "test"}}`)

	if err := Guaranty(&testObject, testJSON); err != nil {
		t.Error(err)
	}

	type testObjectTP2 struct {
		Name string `json:"name" validate:"required" `
	}

	var testObject2 struct {
		Test testObjectTP2 `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	if err := Guaranty(&testObject2, testJSON); err == nil {
		t.Error("Required object, empty JSON object should return error but did not.")
	}
}

func TestEmail(t *testing.T) {

	var testValEmail struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON := jsonFactory(`{"email": "michaeljs@gmail.com"}`)

	if err := Guaranty(&testValEmail, testJSON); err != nil {
		t.Error(err)
	}

	var testValEmail2 struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON = jsonFactory(`{"email": "michaeljs@gail.edu"}`)

	if err := Guaranty(&testValEmail2, testJSON); err != nil {
		t.Error(err)
	}

	var testValEmail3 struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON = jsonFactory(`{"email": "michaeljs.edu"}`)

	if err := Guaranty(&testValEmail3, testJSON); err == nil {
		t.Error("Email test failed, michaeljs.edu passed as a valid email.")
	}

}

func TestIn(t *testing.T) {

	var testValIn struct {
		Test string `json:"special" validate:"in:admin,user,other" `
	}

	testJSON := jsonFactory(`{"special": "admin"}`)

	if err := Guaranty(&testValIn, testJSON); err != nil {
		t.Error(err)
	}

	var testValIn2 struct {
		Test int `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"special": 3}`)

	if err := Guaranty(&testValIn2, testJSON); err != nil {
		t.Error(err)
	}

	var testValIn3 struct {
		Test int `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"special": 6}`)

	if err := Guaranty(&testValIn3, testJSON); err == nil {
		t.Error("6 is not in validate in call, err should not have been nil.")
	}

	var testValIn4 struct {
		Test int `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"what": "fake"}`)

	if err := Guaranty(&testValIn4, testJSON); err == nil {
		t.Error("JSON not related to validation passed in error should not be null.")
	}

}
