package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"./externalservice"
)

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestPOSTCallsAndReturnsJSONfromExternalServicePOST(t *testing.T) {
	// Descirption
	//
	// Write a test that accepts a POST request on the server and sends it the
	// fake external service with the posted form body return the response.
	//
	// Use the externalservice.Client interface to create a mock and assert the
	// client was called <n> times.
	//
	// ---
	//
	// Server should receive a request on
	//
	//  [POST] /api/posts/:id
	//  application/json
	//
	// With the form body
	//
	//  application/x-www-form-urlencoded
	//	title=Hello World!
	//	description=Lorem Ipsum Dolor Sit Amen.
	//
	// The server should then relay this data to the external service by way of
	// the Client POST method and return the returned value out as JSON.
	//
	// ---
	//
	// Assert that the externalservice.Client#POST was called 1 times with the
	// provided `:id` and post body and that the returned Post (from
	// externalservice.Client#POST) is written out as `application/json`.

	req := httptest.NewRequest("POST", "/api/posts/87", nil)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	if err := req.ParseForm(); err != nil {
		panic(err.Error())
	}
	req.Form.Add("title", "Hello World!")
	req.Form.Add("description", "Lorem Ipsum Dolor Sit Amen.")

	rec := httptest.NewRecorder()
	s := Server{
		Client: &externalservice.ClientImpl{
			Posts: make(map[int]*externalservice.Post),
		},
	}
	s.InitializeRoutes()
	s.Router.ServeHTTP(rec, req)

	if http.StatusCreated != rec.Code {
		t.Fatalf("Expected response code %d. Got %d\n", http.StatusCreated, rec.Code)
	}

	// Response content-type should be 'application/json; charset=UTF-8'
	if rec.Header().Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Error("Invalid response Content-Type")
	}

	postPayload := &externalservice.Post{
		ID:          87,
		Title:       "Hello World!",
		Description: "Lorem Ipsum Dolor Sit Amen.",
	}

	// error message should be 'Post id already exists'
	if _, err := s.Client.POST(87, postPayload); err == nil || err.Error() != "Post id already exists" {
		t.Error("Should get error with message  'Post id already exists'")
	}
}

func TestPOSTCallsAndReturnsErrorAsJSONFromExternalServiceGET(t *testing.T) {
	// Description
	//
	// Write a test that accepts a GET request on the server and returns the
	// error returned from the external service.
	//
	// Use the externalservice.Client interface to create a mock and assert the
	// client was called <n> times.
	//
	// ---
	//
	// Server should receive a request on
	//
	//	[GET] /api/posts/:id
	//
	// The server should then return the error from the external service out as
	// JSON.
	//
	// The error response returned from the external service would look like
	//
	//	400 application/json
	//
	//	{
	//		"code": 400,
	//		"message": "Bad Request"
	//	}
	//
	// ---
	//
	// Assert that the externalservice.Client#GET was called 1 times with the
	// provided `:id` and the returned error (above) is output as the response
	// as
	//
	//	{
	//		"code": 400,
	//		"message": "Bad Request",
	//		"path": "/api/posts/:id
	//	}
	//
	// Note: *`:id` should be the actual `:id` in the original request.*

	req := httptest.NewRequest("GET", "/api/posts/87", nil)
	rec := httptest.NewRecorder()
	s := Server{
		Client: &externalservice.ClientImpl{
			Posts: make(map[int]*externalservice.Post),
		},
	}
	s.InitializeRoutes()
	s.Router.ServeHTTP(rec, req)

	// status code should be 400
	if rec.Code != 400 {
		t.Error("Status code should be 400")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(rec.Body.String()), &result); err != nil {
		t.Error("Invalid json response schema")
	}

	path := strings.Split(result["path"].(string), "/")

	id, _ := strconv.Atoi(path[len(path)-1])
	if id != 87 {
		t.Error("Path should be the actual `:id` in the original request.")
	}

	// error message should be 'Post id already exists'
	if _, err := s.Client.GET(87); err == nil || err.Error() != "Post not found" {
		t.Error("Should get error with message 'Post not found'")
	}
}
