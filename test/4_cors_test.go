package test

import (
	"auth-guardian/config"
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"net/http"
	"testing"
)

func TestCORS(t *testing.T) {
	defer seq()()

	setDefaultConfig()
	config.CORSUpstream = true

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
