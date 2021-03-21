package test

import (
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"net/http"
	"regexp"
	"testing"
)

func TestStartupSAML(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockSAMLIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())
}

func TestUnauthorizedSAML(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockSAMLIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())

	requestAndTest(t, &testDefinition{
		Method:     "GET",
		URL:        "http://localhost:3000/",
		StatusCode: http.StatusFound,
	})
}

func TestAuthorizeSAML(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockSAMLIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer ipdServer.Shutdown(context.TODO())
	defer testServer.Shutdown(context.TODO())

	_, client, body := startRequest(t, &testDefinition{
		Method:         "GET",
		URL:            "http://localhost:3000/",
		FollowRedirect: true,
	})

	r := regexp.MustCompile(`(action=")([a-z:\/0-9]*)`)
	groups := r.FindAllStringSubmatch(body, -1)
	url := groups[0][2]

	r = regexp.MustCompile(`(name="SAMLRequest" value=")([a-zA-Z0-9+=\/]*)`)
	groups = r.FindAllStringSubmatch(body, -1)
	samlRequest := groups[0][2]

	r = regexp.MustCompile(`(name="RelayState" value=")([a-zA-Z0-9+=_]*)`)
	groups = r.FindAllStringSubmatch(body, -1)
	relayState := groups[0][2]

	_, client, body = startRequest(t, &testDefinition{
		Client: client,
		Method: "POST",
		URL:    url,
		FormData: map[string][]string{
			"user":        {"noface"},
			"password":    {"SpiritedAway"},
			"SAMLRequest": {samlRequest},
			"RelayState":  {relayState},
		},
		FollowRedirect: true,
	})

	r = regexp.MustCompile(`(action=")([a-z:\/0-9]*)`)
	groups = r.FindAllStringSubmatch(body, -1)
	url = groups[0][2]

	r = regexp.MustCompile(`(name="SAMLResponse" value=")([a-zA-Z0-9+=_]*)`)
	groups = r.FindAllStringSubmatch(body, -1)
	samlResponse := groups[0][2]

	r = regexp.MustCompile(`(name="RelayState" value=")([a-zA-Z0-9+=_]*)`)
	groups = r.FindAllStringSubmatch(body, -1)
	relayState = groups[0][2]

	requestAndTest(t, &testDefinition{
		Client: client,
		Method: "POST",
		URL:    url,
		FormData: map[string][]string{
			"SAMLResponse": {samlResponse},
			"RelayState":   {relayState},
		},
		FollowRedirect: true,
		StatusCode:     http.StatusOK,
		ExpectedBody:   "Hello from test service",
	})
}
