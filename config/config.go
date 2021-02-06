package config

// TestMode specifies if is running in test mode
var TestMode bool

// LogLevel set n for
// any Panic
// n >= 1 Errors
// n >= 2 Warnings
// n >= 3 Infos
// n >= 4 Debugs
var LogLevel int

// LogJSON specifies if logs should have JSON format or formatted text
var LogJSON bool

// LogFile specifies the log file location (default = file logging disabled)
var LogFile string

// Listen specifies where to listen to incoming requests
var Listen string

// ServerCrt specifies the path to the crt file
var ServerCrt string

// ServerKey specifies the path to the key file
var ServerKey string

// Upstream specifies the upstream behind this proxy
var Upstream string

// IsHTTPS is a flag to mark that HTTPS is used
var IsHTTPS bool

// CORSUpstream is a flag to mark that the upstream not accept CORS
var CORSUpstream bool

// ForwardAccessToken specifies whether the access token should be forwarded.
var ForwardAccessToken bool

// ForwardUserinfo specifies whether the userinfo should be forwarded
var ForwardUserinfo bool

// SessionLifetime specifies the lifetime of a session (minutes)
var SessionLifetime int

// ClientID specifies the application's ID
var ClientID string

// ClientSecret specifies the application's secret
var ClientSecret string

// Scopes specifies optional requested permissions
var Scopes []string

// RedirectURL specifies which redirect should be used
var RedirectURL string

// AuthURL specifies the URL to redirect an unauthenticated user
var AuthURL string

// TokenURL specifies the URL from which to get an access token
var TokenURL string

// UserinfoURL specifies the URL from which to get userinfos
var UserinfoURL string

// StateLifetime specifies how long a state is valid (minutes)
var StateLifetime int

// IdpMetadataURL specifies the URL to the IDP metadata
var IdpMetadataURL string

// SelfRootURL specifies the root URL to self
var SelfRootURL string

// SAMLCrt specifies the path to the crt file for SAML
var SAMLCrt string

// SAMLKey specifies the path to the key file for SAML
var SAMLKey string

// DirectoryServerBaseDN specifies the base DN that should be used for the search
var DirectoryServerBaseDN string

// DirectoryServerBindDN specifies the DN to use to bind to the directory server when performing simple authentication
var DirectoryServerBindDN string

// DirectoryServerPort specifies the port of the directory server
var DirectoryServerPort int

// DirectoryServerHost specifies the host of the directory server
var DirectoryServerHost string

// DirectoryServerBindPassword specifies password to use to access the key store contents
var DirectoryServerBindPassword string

// DirectoryServerFilter specifies filter for the user object
var DirectoryServerFilter string

// WhitelistRules specifies rules to whitelist resources
var WhitelistRules []string

// Load loads the config
func Load() (bool, error) {
	definition := map[string]map[string]interface{}{
		"version": {
			"desc":    "Get the version.",
			"type":    "bool",
			"default": false,
		},
		"test-mode": {
			"desc":    "Specifies if is running in test mode.",
			"type":    "bool",
			"default": false,
		},
		"log-level": {
			"desc":    "Set n for {any Panic, n >= 1 Errors, n >= 2 Warnings, n >= 3 Infos, n >= 4 Debugs}.",
			"type":    "int",
			"default": 2,
		},
		"log-file": {
			"desc":    "Specifies the log file location (default = file logging disabled).",
			"type":    "string",
			"default": "",
		},
		"log-json": {
			"desc":    "Specifies if logs should have JSON format or formatted text.",
			"type":    "bool",
			"default": false,
		},
		"listen": {
			"desc":    "Specifies where to listen to incoming requests (default = :8080).",
			"type":    "string",
			"default": ":8080",
		},
		"server-crt": {
			"desc":    "Specifies the path to the crt file.",
			"type":    "string",
			"default": "",
		},
		"server-key": {
			"desc":    "Specifies the path to the key file.",
			"type":    "string",
			"default": "",
		},
		"upstream": {
			"desc":    "Specifies the upstream behind this proxy.",
			"type":    "string",
			"default": "",
		},
		"upstream-cors": {
			"desc":    "Specifies that the upstream not accept CORS and is not on the same domain.",
			"type":    "bool",
			"default": false,
		},
		"forward-userinfo": {
			"desc":    "Specifies whether the userinfo should be forwarded.",
			"type":    "bool",
			"default": false,
		},
		"forward-access-token": {
			"desc":    "Specifies whether the access token should be forwarded.",
			"type":    "bool",
			"default": false,
		},
		"session-lifetime": {
			"desc":    "Specifies the lifetime of a session (minutes).",
			"type":    "int",
			"default": 5,
		},
		"client-id": {
			"desc":    "Specifies the application's ID.",
			"type":    "string",
			"default": "",
		},
		"client-secret": {
			"desc":    "Specifies the application's secret.",
			"type":    "string",
			"default": "",
		},
		"scopes": {
			"desc":    "Specifies optional requested permissions.",
			"type":    "string_array",
			"default": &StringArrayFlag{},
		},
		"redirect-url": {
			"desc":    "Specifies which redirect should be used.",
			"type":    "string",
			"default": "",
		},
		"auth-url": {
			"desc":    "Specifies the URL to redirect an unauthenticated user.",
			"type":    "string",
			"default": "",
		},
		"token-url": {
			"desc":    "Specifies the URL from which to get an access token.",
			"type":    "string",
			"default": "",
		},
		"userinfo-url": {
			"desc":    "Specifies the URL from which to get userinfos.",
			"type":    "string",
			"default": "",
		},
		"state-lifetime": {
			"desc":    "Specifies how long a state is valid (minutes).",
			"type":    "int",
			"default": 5,
		},
		"saml-metadata-url": {
			"desc":    "Specifies the URL to the IDP metadata.",
			"type":    "string",
			"default": "",
		},
		"self-root-url": {
			"desc":    "Specifies the root URL to self.",
			"type":    "string",
			"default": "",
		},
		"saml-crt": {
			"desc":    "Specifies the path to the crt file for SAML.",
			"type":    "string",
			"default": "",
		},
		"saml-key": {
			"desc":    "Specifies the path to the key file for SAML.",
			"type":    "string",
			"default": "",
		},
		"ds-base-dn": {
			"desc":    "Specifies the base DN that should be used for the search.",
			"type":    "string",
			"default": "",
		},
		"ds-bind-dn": {
			"desc":    "Specifies the DN to use to bind to the directory server when performing simple authentication.",
			"type":    "string",
			"default": "",
		},
		"ds-port": {
			"desc":    "Specifies the port of the directory server.",
			"type":    "int",
			"default": 389,
		},
		"ds-host": {
			"desc":    "Specifies the host of the directory server.",
			"type":    "string",
			"default": "",
		},
		"ds-bind-password": {
			"desc":    "Specifies password to use to access the key store contents.",
			"type":    "string",
			"default": "",
		},
		"ds-filter": {
			"desc":    "Specifies filter for the user object.",
			"type":    "string",
			"default": "",
		},
		"whitelist-rule": {
			"desc":    "Specifies rules to whitelist resources.",
			"type":    "string_array",
			"default": &StringArrayFlag{},
		},
	}

	// Config from env if exists
	getConfigFromEnv(&definition)

	// Get config from file if defined
	err := getConfigFromFile(&definition)
	if err != nil {
		return false, err
	}

	// Get config from arguments
	getConfigFromArguments(&definition)

	// Set mostly priories config value
	TestMode = getMostlyPrioriesConfigKey(definition["test-mode"]).(bool)

	LogLevel = getMostlyPrioriesConfigKey(definition["log-level"]).(int)
	LogJSON = getMostlyPrioriesConfigKey(definition["log-json"]).(bool)
	LogFile = getMostlyPrioriesConfigKey(definition["log-file"]).(string)

	Listen = getMostlyPrioriesConfigKey(definition["listen"]).(string)

	ServerCrt = getMostlyPrioriesConfigKey(definition["server-crt"]).(string)
	ServerKey = getMostlyPrioriesConfigKey(definition["server-key"]).(string)

	Upstream = getMostlyPrioriesConfigKey(definition["upstream"]).(string)
	CORSUpstream = getMostlyPrioriesConfigKey(definition["upstream-cors"]).(bool)

	ForwardAccessToken = getMostlyPrioriesConfigKey(definition["forward-access-token"]).(bool)
	ForwardUserinfo = getMostlyPrioriesConfigKey(definition["forward-userinfo"]).(bool)
	SessionLifetime = getMostlyPrioriesConfigKey(definition["session-lifetime"]).(int)

	ClientID = getMostlyPrioriesConfigKey(definition["client-id"]).(string)
	ClientSecret = getMostlyPrioriesConfigKey(definition["client-secret"]).(string)
	Scopes = getMostlyPrioriesConfigKey(definition["scopes"]).([]string)
	RedirectURL = getMostlyPrioriesConfigKey(definition["redirect-url"]).(string)
	AuthURL = getMostlyPrioriesConfigKey(definition["auth-url"]).(string)
	TokenURL = getMostlyPrioriesConfigKey(definition["token-url"]).(string)
	UserinfoURL = getMostlyPrioriesConfigKey(definition["userinfo-url"]).(string)
	StateLifetime = getMostlyPrioriesConfigKey(definition["state-lifetime"]).(int)

	IdpMetadataURL = getMostlyPrioriesConfigKey(definition["saml-metadata-url"]).(string)
	SelfRootURL = getMostlyPrioriesConfigKey(definition["self-root-url"]).(string)
	SAMLCrt = getMostlyPrioriesConfigKey(definition["saml-crt"]).(string)
	SAMLKey = getMostlyPrioriesConfigKey(definition["saml-key"]).(string)

	DirectoryServerBaseDN = getMostlyPrioriesConfigKey(definition["ds-base-dn"]).(string)
	DirectoryServerBindDN = getMostlyPrioriesConfigKey(definition["ds-bind-dn"]).(string)
	DirectoryServerPort = getMostlyPrioriesConfigKey(definition["ds-port"]).(int)
	DirectoryServerHost = getMostlyPrioriesConfigKey(definition["ds-host"]).(string)
	DirectoryServerBindPassword = getMostlyPrioriesConfigKey(definition["ds-bind-password"]).(string)
	DirectoryServerFilter = getMostlyPrioriesConfigKey(definition["ds-filter"]).(string)

	WhitelistRules = getMostlyPrioriesConfigKey(definition["whitelist-rule"]).([]string)

	// Set http flag
	IsHTTPS = (ServerKey != "" && ServerCrt != "")

	return getMostlyPrioriesConfigKey(definition["version"]).(bool), nil
}
