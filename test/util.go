package test

import (
	"auth-guardian/config"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

type testDefinition struct {
	Client                   *http.Client
	URL                      string
	Method                   string
	FormData                 map[string][]string
	FollowRedirect           bool
	FollowClientSideRedirect bool
	StatusCode               int
	ExpectedBody             string
}

var seqMutex sync.Mutex

func seq() func() {
	seqMutex.Lock()
	return func() {
		seqMutex.Unlock()
	}
}

func setDefaultConfig() {
	config.Listen = ":3000"
	config.LogLevel = 0
	config.LogJSON = false
	config.LogFile = ""
	config.ServerCrt = ""
	config.ServerKey = ""
	config.Upstream = ""
	config.CORSUpstream = false
	config.ForwardUserinfo = true
	config.ForwardAccessToken = true
	config.SessionLifetime = 5
	config.ClientID = ""
	config.ClientSecret = ""
	config.Scopes = []string{}
	config.RedirectURL = ""
	config.AuthURL = ""
	config.TokenURL = ""
	config.UserinfoURL = ""
	config.StateLifetime = 5
	config.IdpMetadataURL = ""
	config.IdpRegisterURL = ""
	config.SelfRootURL = ""
	config.SAMLCrt = ""
	config.SAMLKey = ""
	config.DirectoryServerBaseDN = ""
	config.DirectoryServerBindDN = ""
	config.DirectoryServerPort = 389
	config.DirectoryServerHost = ""
	config.DirectoryServerBindPassword = ""
	config.DirectoryServerFilter = ""
	config.Rules = []config.RuleConfig{}
	config.IsHTTPS = false

	time.Sleep(1 * time.Second)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("got %v, want %v", a, b)
	}
}

func notError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func startRequest(t *testing.T, td *testDefinition) (*http.Response, *http.Client, string) {
	return request(t, td, 0)
}

func request(t *testing.T, td *testDefinition, counter int) (*http.Response, *http.Client, string) {
	var client *http.Client

	if td.Client != nil {
		client = td.Client
	} else {
		jar, err := cookiejar.New(nil)
		notError(t, err)
		client = &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 15 * time.Second,
				}).Dial,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				MaxIdleConnsPerHost: 14,
				TLSHandshakeTimeout: 15 * time.Second,
			},
			Timeout: 15 * time.Second,
			Jar:     jar,
		}
		td.Client = client
	}

	if !td.FollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var req *http.Request
	var err error
	if td.FormData != nil && len(td.FormData) > 0 {
		formData := url.Values{}
		for key, value := range td.FormData {
			formData[key] = value
		}

		req, err = http.NewRequest(td.Method, td.URL, strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(td.Method, td.URL, nil)
	}
	notError(t, err)

	res, err := client.Do(req)
	notError(t, err)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	notError(t, err)
	defer res.Body.Close()
	body := string(bodyBytes)

	if td.FollowClientSideRedirect && counter <= 5 {
		r := regexp.MustCompile(`(href=|window.location)([= '"]*)([a-zA-Z:/0-9?=&;+%_-]*)`)
		groups := r.FindAllStringSubmatch(body, -1)

		if len(groups) >= 1 && len(groups[0]) >= 3 {
			td.Method = "GET"
			td.URL = groups[0][3]
			counter++
			return request(t, td, counter)
		}
	}

	return res, client, body
}

func requestAndTest(t *testing.T, td *testDefinition) *http.Client {
	res, client, body := startRequest(t, td)

	assertEqual(t, res.StatusCode, td.StatusCode)

	if td.ExpectedBody != "" {
		assertEqual(t, body, td.ExpectedBody)
	}

	return client
}
