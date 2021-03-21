package test

import (
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"net/http"
	"testing"
)

func TestStartupLDAP(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockLDAPIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer func() { mocked.LDAPListenerCloseFlag = true }()
	defer (*ipdServer).Close()
	defer testServer.Shutdown(context.TODO())
}

func TestUnauthorizedLDAP(t *testing.T) {
	defer seq()()

	setDefaultConfig()

	testServer := mocked.RunMockTestService()
	ipdServer := mocked.RunMockLDAPIDP()
	server := server.Run(t)

	defer server.Shutdown(context.TODO())
	defer func() { mocked.LDAPListenerCloseFlag = true }()
	defer (*ipdServer).Close()
	defer testServer.Shutdown(context.TODO())

	requestAndTest(t, &testDefinition{
		Method:     "GET",
		URL:        "http://localhost:3000/",
		StatusCode: http.StatusUnauthorized,
	})
}

func TestAuthorizeLDAP(t *testing.T) {

}
