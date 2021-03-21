package test

import (
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"net/http"
	"testing"
)

func TestStartupOAuth(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockOAuthIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())
}

func TestUnauthorizedOAuth(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockOAuthIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())

	requestAndTest(t, &testDefinition{
		Method:     "GET",
		URL:        "http://localhost:3000/",
		StatusCode: http.StatusTemporaryRedirect,
	})
}

func TestAuthorizeOAuth(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockOAuthIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())

	requestAndTest(t, &testDefinition{
		Method:                   "GET",
		URL:                      "http://localhost:3000/",
		FollowRedirect:           true,
		FollowClientSideRedirect: true,
		StatusCode:               http.StatusOK,
		ExpectedBody:             "Hello from test service",
	})
}
