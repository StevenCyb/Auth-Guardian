# Documentation
## DOC -- Table Of Contents
1) [How to configure](#how-to-configure)
2) [Configuration Options](#configuration-options)
3) [Examples](#examples) \
3.1) [Keycloak-OIDC](#keycloak-oidc) \
3.1) [Github](#github)

### How to configure
The configuration can be done using environment variables, arguments and a configuration file.
The configuration file cannot be selected, it is explicitly searched for a `config.yml` in the project folder.

Two or all three variants can also be used together.
Please note that the configurations overwrite each other at a fixed order.
This means that the configurations from the environment variables are overwritten by the configuration file, which in turn is overwritten by the arguments. 

### Configuration Options
The following table shows and describes the options currently available. 
They also show which of these options are mandatory.

| Argument           | Type      | Required | Default  | Description                                                                         |
|--------------------|-----------|:--------:|----------|-------------------------------------------------------------------------------------|
| --auth-url         | string    | Yes      | -        | Specifies the URL to redirect an unauthenticated user.                              |
| --client-id        | string    | Yes      | -        | Specifies the application's ID                                                      |
| --client-secret    | string    | Yes      | -        | Specifies the application's secret.                                                 |
| --listen           | string    | No       | :8080    | Specifies where to listen to incoming requests.                                     |
| --log-file         | string    | No       | disabled | Specifies the log file location (default = file logging disabled).                  |
| --log-json         | bool      | No       | false    | Specifies if logs should have JSON format or formatted text.                        |
| --log-level        | int       | No       | 2        | Set n for {any Panic, n >= 1 Errors, n >= 2 Warnings, n >= 3 Infos, n >= 4 Debugs}. |
| --redirect-url     | string    | *1       | -        | Specifies which redirect should be used.                                            |
| --scopes           | string 2* | No       | -        | Specifies optional requested permissions.                                           |
| --server-crt       | string    | No       | -        | Specifies the path to the crt file.                                                 |
| --server-key       | string    | No       | -        | Specifies the path to the key file.                                                 |
| --session-lifetime | int       | No       | 5        | Specifies the lifetime of a session (minutes).                                      |
| --state-lifetime   | int       | No       | 5        | Specifies how long a state is valid (minutes)                                       |
| --token-url        | string    | Yes      | -        | Specifies the URL from which to get an access token.                                |
| --upstream         | string    | Yes      | -        | Specifies the upstream behind this proxy.                                           |
| --upstream-cors    | bool      | No       | false    | Specifies that the upstream not accept CORS and is not on the same domain.          |
| --userinfo-url     | string    | No       | -        | Specifies the URL from which to get userinfos.                                      |
| --version          | bool      | No       | false    | Get the version.                                                                    |

* 1* : Required when identity and access management has multiple redirect URls
* 2* : String in array notation e.g. 1,2,3,4,5 or a,b,c,d,e,f

### Examples
#### Keycloak-OIDC
##### Using environment variables
```
export client-id={client-id}
export client-secret={client-secret}
export auth-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/auth
export token-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/token
export userinfo-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/userinfo
export scopes=email,roles
```
##### Using configuration file:
```yml
client-id: {client-id}
client-secret: {client-secret}
auth-url: "http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/auth"
token-url: "http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/token"
userinfo-url: "http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/userinfo"
scopes: "email,roles"
```
##### Using arguments
```shell
go run main.go \
--client-id={client-id} \
--client-secret={client-secret} \
--auth-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/auth \
--token-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/token \
--userinfo=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/userinfo \
--scopes=email,roles
```
#### Github
##### Using environment variables
```
export client-id={client-id}
export client-secret={client-secret}
export auth-url=https://github.com/login/oauth/authorize
export token-url=https://github.com/login/oauth/access_token
export userinfo-url=https://api.github.com/user
export scopes=read:user,user:email
```
##### Using configuration file:
```yml
client-id: {client-id}
client-secret: {client-secret}
auth-url: "https://github.com/login/oauth/authorize"
token-url: "https://github.com/login/oauth/access_token"
userinfo-url: "https://api.github.com/user"
scopes: "read:user,user:email"
```
##### Using arguments
```shell
go run main.go \
--client-id={client-id} \
--client-secret={client-secret} \
--auth-url=https://github.com/login/oauth/authorize \
--token-url=https://github.com/login/oauth/access_token \
--userinfo=https://api.github.com/user \
--scopes=read:user,user:email
```
