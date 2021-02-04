# Documentation
## DOC -- Table Of Contents
1) [How to configure](#how-to-configure)
2) [Configuration Options](#configuration-options)
3) [Examples](#examples) \
3.1) [Keycloak-OIDC](#keycloak-oidc) \
3.2) [Keycloak-SAML](#keycloak-saml) \
3.3) [Github-OAuth](#github-oauth) \
3.4) [ForumSys-LDAP](#forumsys-ldap)

### How to configure
The configuration can be done using environment variables, arguments and a configuration file.
The configuration file cannot be selected, it is explicitly searched for a `config.yaml` in the project folder.

Two or all three variants can also be used together.
Please note that the configurations overwrite each other at a fixed order.
This means that the configurations from the environment variables are overwritten by the configuration file, which in turn is overwritten by the arguments. 

### Configuration Options
The following table shows and describes the options currently available. 
They also show which of these options are mandatory.

| Argument               | Type      | Required | Default  | Description                                                                                    |
|------------------------|-----------|:--------:|----------|------------------------------------------------------------------------------------------------|
| --auth-url             | string    | No 3*    | -        | Specifies the URL to redirect an unauthenticated user.                                         |
| --client-id            | string    | No 3*    | -        | Specifies the application's ID                                                                 |
| --client-secret        | string    | No 3*    | -        | Specifies the application's secret.                                                            |
| --ds-base-dn           | string    | No 6*    | -        | Specifies the base DN that should be used for the search.                                      |
| --ds-bind-dn           | string    | No 6*    | -        | Specifies the DN to use to bind to the directory server when performing simple authentication. |
| --ds-port              | int       | No 6*    | 389      | Specifies the port of the directory server.                                                    |
| --ds-host              | string    | No 6*    | -        | Specifies the host of the directory server.                                                    |
| --ds-bind-password     | string    | No 6*    | -        | Specifies password to use to access the key store contents.                                    |
| --ds-filter            | string    | No 6*    | -        | Specifies filter for the user object.                                                          |
| --forward-access-token | bool      | No 5*    | false    | Specifies whether the access token should be forwarded.                                        |
| --forward-userinfo     | bool      | No       | false    | Specifies whether the userinfo should be forwarded.                                            |
| --listen               | string    | No       | :8080    | Specifies where to listen to incoming requests.                                                |
| --log-file             | string    | No       | disabled | Specifies the log file location (default = file logging disabled).                             |
| --log-json             | bool      | No       | false    | Specifies if logs should have JSON format or formatted text.                                   |
| --log-level            | int       | No       | 2        | Set n for {any Panic, n >= 1 Errors, n >= 2 Warnings, n >= 3 Infos, n >= 4 Debugs}.            |
| --redirect-url         | string    | *1       | -        | Specifies which redirect should be used.                                                       |
| --saml-crt             | string    | No 4*    | -        | Specifies the path to the crt file for SAML.                                                   |
| --saml-key             | string    | No 4*    | -        | Specifies the path to the key file for SAML.                                                   |
| --saml-metadata-url    | string    | No 4*    | -        | Specifies the URL to the IDP metadata.                                                         |
| --scopes               | string 2* | No       | -        | Specifies optional requested permissions.                                                      |
| --self-root-url        | string    | No 4*    | -        | Specifies the root URL to self.                                                                |
| --server-crt           | string    | No       | -        | Specifies the path to the crt file for SAML.                                                   |
| --server-key           | string    | No       | -        | Specifies the path to the key file for SAML.                                                   |
| --session-lifetime     | int       | No       | 5        | Specifies the lifetime of a session (minutes).                                                 |
| --state-lifetime       | int       | No       | 5        | Specifies how long a state is valid (minutes)                                                  |
| --test-mode            | bool      | No       | false    | TestMode specifies if is running in test mode.                                                 |
| --token-url            | string    | No 3*    | -        | Specifies the URL from which to get an access token.                                           |
| --upstream             | string    | Yes      | -        | Specifies the upstream behind this proxy.                                                      |
| --upstream-cors        | bool      | No       | false    | Specifies that the upstream not accept CORS and is not on the same domain.                     |
| --userinfo-url         | string    | No       | -        | Specifies the URL from which to get userinfos.                                                 |
| --version              | bool      | No       | false    | Get the version.                                                                               |

* 1* : Required when identity and access management has multiple redirect URls
* 2* : String in array notation e.g. 1,2,3,4,5 or a,b,c,d,e,f
* 3* : Yes if you want to use OAuth
* 4* : Yes if you want to use SAML
* 5* : Works only with OAuth
* 6* : Yes if you want to use LDAP

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
#### Keycloak-SAML
##### Using environment variables
```
export saml-crt={path-to-saml-crt-file}
export saml-key={path-to-saml-key-file}
export saml-metadata-url="http://{keycloak-url}/auth/realms/master/protocol/saml/descriptor"
export self-root-url={self-root-url}
```
##### Using configuration file:
```yml
saml-crt: {path-to-saml-crt-file}
saml-key: {path-to-saml-key-file}
saml-metadata-url: "http://{keycloak-url}/auth/realms/master/protocol/saml/descriptor"
self-root-url: {self-root-url}
```
##### Using arguments
```shell
go run main.go \
--saml-crt={path-to-saml-crt-file} \
--saml-key={path-to-saml-key-file} \
--saml-metadata-url="http://{keycloak-url}/auth/realms/master/protocol/saml/descriptor" \
--self-root-url={self-root-url}
```
#### Github-OAuth
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
#### ForumSys-LDAP
##### Using environment variables
```
export ds-base-dn="dc=example,dc=com"
export ds-bind-dn="cn=read-only-admin,dc=example,dc=com"
export ds-port=389
export ds-host="ldap.forumsys.com"
export ds-bind-password="password"
export ds-filter="(uid=%s)"
```
##### Using configuration file:
```yml
ds-base-dn: "dc=example,dc=com"
ds-bind-dn: "cn=read-only-admin,dc=example,dc=com"
ds-port: 389
ds-host: "ldap.forumsys.com"
ds-bind-password: "password"
ds-filter: "(uid=%s)"
```
##### Using arguments
```shell
go run main.go \
--ds-base-dn="dc=example,dc=com" \
--ds-bind-dn="cn=read-only-admin,dc=example,dc=com" \
--ds-port=389 \
--ds-host="ldap.forumsys.com" \
--ds-bind-password="password" \
--ds-filter="(uid=%s)"
```