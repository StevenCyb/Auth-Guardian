package test

import (
	"auth-guardian/config"
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"net/http"
	"testing"
)

func TestRequiredRule(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	config.Rules = append(config.Rules,
		config.RuleConfig{
			Type:   "required",
			Method: []string{},
			Path:   "^(/notallow)$",
			Userinfo: map[string]string{
				"role": "^monster$",
			},
		},
	)
	config.Rules = append(config.Rules,
		config.RuleConfig{
			Type:   "required",
			Method: []string{},
			Path:   "^(/)$",
			Userinfo: map[string]string{
				"role": "^Superhero$",
			},
		},
	)

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockOAuthIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())

	client := requestAndTest(t, &testDefinition{
		Method:                   "GET",
		URL:                      "http://localhost:3000/",
		FollowRedirect:           true,
		FollowClientSideRedirect: true,
		StatusCode:               http.StatusOK,
		ExpectedBody:             "Hello from test service",
	})

	requestAndTest(t, &testDefinition{
		Client:                   client,
		Method:                   "GET",
		URL:                      "http://localhost:3000/notallow",
		FollowRedirect:           true,
		FollowClientSideRedirect: true,
		StatusCode:               http.StatusForbidden,
	})
}
