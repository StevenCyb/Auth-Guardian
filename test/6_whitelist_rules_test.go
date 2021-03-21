package test

import (
	"auth-guardian/config"
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"net/http"
	"testing"
)

func TestWhitelistRule(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	config.Rules = append(config.Rules,
		config.RuleConfig{
			Type:   "whitelist",
			Method: []string{},
			Path:   "^(/)$",
		})

	config.Rules = append(config.Rules,
		config.RuleConfig{
			Type: "whitelist",
			Path: "(.js)$",
		})

	config.Rules = append(config.Rules,
		config.RuleConfig{
			Type:   "whitelist",
			Method: []string{"GET"},
			Path:   "(.css)$",
		})

	config.Rules = append(config.Rules,
		config.RuleConfig{
			Type:   "whitelist",
			Method: []string{"GET", "POST", "PUT", "DELETE"},
			Path:   "^(/favicon.ico)$",
		})

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockOAuthIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())

	requestAndTest(t, &testDefinition{
		Method:     "GET",
		URL:        "http://localhost:3000/mirror",
		StatusCode: http.StatusTemporaryRedirect,
	})

	requestAndTest(t, &testDefinition{
		Method:       "GET",
		URL:          "http://localhost:3000/",
		StatusCode:   http.StatusOK,
		ExpectedBody: "Hello from test service",
	})

	requestAndTest(t, &testDefinition{
		Method:       "GET",
		URL:          "http://localhost:3000/script.js",
		StatusCode:   http.StatusOK,
		ExpectedBody: "console.log('I'm a script.');",
	})

	requestAndTest(t, &testDefinition{
		Method:       "GET",
		URL:          "http://localhost:3000/style.css",
		StatusCode:   http.StatusOK,
		ExpectedBody: "body { background-color: #333333; }",
	})
}
