# Documentation
## DOC -- Table Of Contents
1) [How to configure](#how-to-configure)
2) [Configuration Options](#configuration-options)
3) [Authentication](#authentication) \
3.1) [Keycloak-OIDC](#keycloak-oidc) \
3.2) [Keycloak-SAML](#keycloak-saml) \
3.3) [Github-OAuth](#github-oauth) \
3.4) [ForumSys-LDAP](#forumsys-ldap)
4) [Rules](#rules) \
4.1) [How to define rules](#how-to-define-rules) \
4.1.1) [Via environment](#via-environment) \
4.1.2) [Via configuration file](#via-configuration-file) \
4.1.3) [Via arguments](#via-arguments) \
4.2) [Whitelist rules](#whitelist-rules) \
4.3) [Required rules](#required-rules) \
4.4) [Disallow rules](#disallow-rules)

## How to configure
The configuration can be done using environment variables, arguments and a configuration file.
The configuration file cannot be selected, it is explicitly searched for a `config.yaml` in the project folder.

Two or all three variants can also be used together.
Please note that the configurations overwrite each other at a fixed order.
This means that the configurations from the environment variables are overwritten by the configuration file, which in turn is overwritten by the arguments. 

## Configuration Options
The following table shows and describes the options currently available. 
They also show which of these options are mandatory.

| Argument               | Type            | Required | Default  | Description                                                                                    |
|------------------------|-----------------|:--------:|----------|------------------------------------------------------------------------------------------------|
| --auth-url             | string          | No 2*    | -        | Specifies the URL to redirect an unauthenticated user.                                         |
| --client-id            | string          | No 2*    | -        | Specifies the application's ID                                                                 |
| --client-secret        | string          | No 2*    | -        | Specifies the application's secret.                                                            |
| --ds-base-dn           | string          | No 5*    | -        | Specifies the base DN that should be used for the search.                                      |
| --ds-bind-dn           | string          | No 5*    | -        | Specifies the DN to use to bind to the directory server when performing simple authentication. |
| --ds-port              | int             | No 5*    | 389      | Specifies the port of the directory server.                                                    |
| --ds-host              | string          | No 5*    | -        | Specifies the host of the directory server.                                                    |
| --ds-bind-password     | string          | No 5*    | -        | Specifies password to use to access the key store contents.                                    |
| --ds-filter            | string          | No 5*    | -        | Specifies filter for the user object.                                                          |
| --forward-access-token | bool            | No 4*    | false    | Specifies whether the access token should be forwarded.                                        |
| --forward-userinfo     | bool            | No       | false    | Specifies whether the userinfo should be forwarded.                                            |
| --listen               | string          | No       | :8080    | Specifies where to listen to incoming requests.                                                |
| --log-file             | string          | No       | disabled | Specifies the log file location (default = file logging disabled).                             |
| --log-json             | bool            | No       | false    | Specifies if logs should have JSON format or formatted text.                                   |
| --log-level            | int             | No       | 2        | Set n for {any Panic, n >= 1 Errors, n >= 2 Warnings, n >= 3 Infos, n >= 4 Debugs}.            |
| --mock-ldap            | bool            | No       | false    | Specifies if should run mocked LDAP IPD.                                                       |
| --mock-oauth           | bool            | No       | false    | Specifies if should run mocked OAuth IDP.                                                      |
| --mock-saml            | bool            | No       | false    | Specifies if should run mocked SAML IDP.                                                       |
| --mock-test-service    | bool            | No       | false    | Specifies if is running in test mode.                                          |
| --redirect-url         | string          | *1       | -        | Specifies which redirect should be used.                                                       |
| --saml-crt             | string          | No 3*    | -        | Specifies the path to the crt file for SAML.                                                   |
| --saml-key             | string          | No 3*    | -        | Specifies the path to the key file for SAML.                                                   |
| --saml-metadata-url    | string          | No 3*    | -        | Specifies the URL to the IDP metadata.                                                         |
| --saml-register-url    | string          | No       | -        | Specifies the URL to register this SP.                                                         |
| --scopes               | string array    | No       | -        | Specifies optional requested permissions.                                                      |
| --self-root-url        | string          | No 3*    | -        | Specifies the root URL to self.                                                                |
| --server-crt           | string          | No       | -        | Specifies the path to the crt file for SAML.                                                   |
| --server-key           | string          | No       | -        | Specifies the path to the key file for SAML.                                                   |
| --session-lifetime     | int             | No       | 5        | Specifies the lifetime of a session (minutes).                                                 |
| --state-lifetime       | int             | No       | 5        | Specifies how long a state is valid (minutes)                                                  |
| --token-url            | string          | No 2*    | -        | Specifies the URL from which to get an access token.                                           |
| --upstream             | string          | Yes      | -        | Specifies the upstream behind this proxy.                                                      |
| --upstream-cors        | bool            | No       | false    | Specifies that the upstream not accept CORS and is not on the same domain.                     |
| --userinfo-url         | string          | No       | -        | Specifies the URL from which to get userinfos.                                                 |
| --rules                | string array 6* | No       | -        | Specifies rules for resources.                                                                 |

* 1* : Required when identity and access management has multiple redirect URls
* 2* : Yes if you want to use OAuth
* 3* : Yes if you want to use SAML
* 4* : Works only with OAuth
* 5* : Yes if you want to use LDAP
* 6* : Each string must be formatted as described on [whitelist rules](#rules)

## Authentication
### Keycloak-OIDC
#### Using environment variables
```bash
export client-id={client-id}
export client-secret={client-secret}
export auth-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/auth
export token-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/token
export userinfo-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/userinfo
export scopes=email,roles
```
#### Using configuration file:
```yml
client-id: {client-id}
client-secret: {client-secret}
auth-url: "http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/auth"
token-url: "http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/token"
userinfo-url: "http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/userinfo"
scopes: "email,roles"
```
#### Using arguments
```bash
go run main.go \
--client-id={client-id} \
--client-secret={client-secret} \
--auth-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/auth \
--token-url=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/token \
--userinfo=http://{keycloak-url}/auth/realms/{realm-name}/protocol/openid-connect/userinfo \
--scopes=email,roles
```
### Keycloak-SAML
#### Using environment variables
```bash
export saml-crt={path-to-saml-crt-file}
export saml-key={path-to-saml-key-file}
export saml-metadata-url="http://{keycloak-url}/auth/realms/master/protocol/saml/descriptor"
export self-root-url={self-root-url}
```
#### Using configuration file:
```yml
saml-crt: {path-to-saml-crt-file}
saml-key: {path-to-saml-key-file}
saml-metadata-url: "http://{keycloak-url}/auth/realms/master/protocol/saml/descriptor"
self-root-url: {self-root-url}
```
#### Using arguments
```bash
go run main.go \
--saml-crt={path-to-saml-crt-file} \
--saml-key={path-to-saml-key-file} \
--saml-metadata-url="http://{keycloak-url}/auth/realms/master/protocol/saml/descriptor" \
--self-root-url={self-root-url}
```
### Github-OAuth
#### Using environment variables
```bash
export client-id={client-id}
export client-secret={client-secret}
export auth-url=https://github.com/login/oauth/authorize
export token-url=https://github.com/login/oauth/access_token
export userinfo-url=https://api.github.com/user
export scopes=read:user,user:email
```
#### Using configuration file:
```yml
client-id: {client-id}
client-secret: {client-secret}
auth-url: "https://github.com/login/oauth/authorize"
token-url: "https://github.com/login/oauth/access_token"
userinfo-url: "https://api.github.com/user"
scopes: "read:user,user:email"
```
#### Using arguments
```bash
go run main.go \
--client-id={client-id} \
--client-secret={client-secret} \
--auth-url=https://github.com/login/oauth/authorize \
--token-url=https://github.com/login/oauth/access_token \
--userinfo=https://api.github.com/user \
--scopes=read:user,user:email
```
### ForumSys-LDAP
#### Using environment variables
```bash
export ds-base-dn="dc=example,dc=com"
export ds-bind-dn="cn=read-only-admin,dc=example,dc=com"
export ds-port=389
export ds-host="ldap.forumsys.com"
export ds-bind-password="password"
export ds-filter="(uid=%s)"
```
#### Using configuration file:
```yml
ds-base-dn: "dc=example,dc=com"
ds-bind-dn: "cn=read-only-admin,dc=example,dc=com"
ds-port: 389
ds-host: "ldap.forumsys.com"
ds-bind-password: "password"
ds-filter: "(uid=%s)"
```
#### Using arguments
```bash
go run main.go \
--ds-base-dn="dc=example,dc=com" \
--ds-bind-dn="cn=read-only-admin,dc=example,dc=com" \
--ds-port=389 \
--ds-host="ldap.forumsys.com" \
--ds-bind-password="password" \
--ds-filter="(uid=%s)"
```

## Rules
Rules are defined as JSON (if using environment or arguments) or YAML (if using config file).
Currently, rules support the following attributes:

| Key                 | Type          | Required |
|---------------------|---------------|----------|
| type                | string 1*     | Yes      |
| method              | string array  | No       |
| path                | regex string  | No       |
| userinfo            | string map 2* | No       |
| query-parameter     | string map 2* | No       |
| json-body-parameter | string map 2* | No       |

* 1* : Valid types are `whitelist`, `required` and `disallow`
* 2* : "path.to.value.dot.separated": "regex string"

It's required to specifie a valid `taype`.
But there is no need to set a value for other keys that are not required.
Just skip unwanted keys and they will be automatically ignored.
### How to define rules
#### Via environment
Rules are defined as JSON **array** on the environment.
Simply write a JSON array with rules and validate it e.g. on [jsonformatter.curiousconcept](https://jsonformatter.curiousconcept.com/).
```bash
export rules='[{"type": "", "method": [], path: "", "userinfo": {}, "query-parameter": {}, "json-body-parameter": {}}, {"type": "", "method": [], path: "", "userinfo": {}, "query-parameter": {}, "json-body-parameter": {}}]'
```
#### Via configuration file
Rules are defined as YAML in the configuration file.
**You must escape yaml special characters like "\"!**
For example a regex escape of `\` you are using `\\`.
So in YAML you an additional escape (like a chain), so use `\\\`.
```yml
rules:
- type: ""
  method: []
  path: ""
  userinfo: {}
  query-parameter: {}
  query-parameter: {}
- type: ""
  method: []
  path: ""
  userinfo: {}
  query-parameter: {}
  query-parameter: {}
```
#### Via arguments
The rule argument command expect a list of strings that contains a JSON rule definition.
```bash
go run main.go \
--rules='{"type": "", "method":[], "path": "", "userinfo": {}, "query-parameter": {}, "json-body-parameter": {}}' \
--rules='{"type": "", "method":[], "path": "", "userinfo": {}, "query-parameter": {}, "json-body-parameter": {}}'
```
### Whitelist rules
This rules define resources which can be accessed without authentication (public resources).
Whitelist rules are defined as describe in [How to define rules](#How-to-define-rules) and support `method` and `path` keys.

Let's assume you want to whitelist all `scripts`, `styles`, the `favicon` and the root path `/`.
Just for this example we want to allow `GET`, `POST`, `PUT`, `DELETE` on `favicon`. 
In this case you can define it as follows.
#### Using environment variables
```bash
export rules='[{"type": "whitelist", "path": "^(/)$"},{"type": "whitelist", "path": "^.*(.js)$"},{"type": "whitelist", "path": "^.*(.css)$"},{"type": "whitelist", "method": ["GET","POST","PUT","DELETE"],"path": "^(/favicon.ico)$"}]'
```
#### Using configuration file
```yml
rules:
- type: "whitelist"
  path: "^(/)$"
- type: "whitelist"
  path: "^.*(.js)$"
- type: "whitelist"
  path: "^.*(.css)$"
- type: "whitelist"
  method: 
  - GET
  - POST
  - PUT
  - DELETE
  path: "^(/favicon.ico)$"
```
#### Using arguments
```bash
go run main.go \
--rules='{"type": "whitelist", "path": "^(/)$}' \
--rules='{"type": "whitelist", "path":"(.js)$"}' \
--rules='{"type": "whitelist", "path": "(.css)$"}' \
--rules='{"type": "whitelist", "method": ["GET","POST","PUT","DELETE"], "path": "^(/favicon.ico)$"}'
```
### Required rules
This rules define resources that can be only accessed on defined conditions (e.g. if the user belongs to group or has a specific role).
Required rules support `userinfo`, `query-parameter` and `json-body-parameter` keys and the resource are defined with `method` and `path`.

For example if we want to allow root `/` access for users which belongs to the group `inetOrgPerson` and have an email address from `@test-company.com`.
We first can check how the userinfo is structured by seting the `--mock-test-service` argument to `true`.
Now we can navigate to the path `localhost/mirror` and decode the base64 encoded userinfo from the cookie to get the JSON.
This JSON could look like:
```json
{
  "username": "tesla",
  "email": "tesla@test-company.com",
  "company": {
    "group": [
      "developer",
      "inetOrgPerson",
      "internal"
    ]
  }
}
``` 
Then we could define the required role as follows.
#### Using environment variables
```bash
export rules='[{"type": "required", "path": "^(/)$", "userinfo": {"company.group": "^inetOrgPerson$", "email": "(@test-company.com)$"}}]'
```
#### Using configuration file
```yml
rules:
- type: "required"
  path: "^(/)$"
  userinfo:
    company.group: "^inetOrgPerson$"
    email: "(@test-company.com)$"
```
#### Using arguments
```bash
go run main.go \
--rules='{"type": "required", "path": "^(/)$", "userinfo": {"company.group": "^inetOrgPerson$", "email": "(@test-company.com)$"}}'
```
### Disallow rules
This rules define resources that can't be accessed on defined conditions (e.g. if all users are allowed excluding externals).
Required rules support `userinfo`, `query-parameter` and `json-body-parameter` keys and the resource are defined with `method` and `path`.

For example if we want to allow root `/` access for all users excluding externals which belogs to the group `externals`.
We first can check how the userinfo is structured by seting the `--mock-test-service` argument to `true`.
Now we can navigate to the path `localhost/mirror` and decode the base64 encoded userinfo from the cookie to get the JSON.
This JSON could look like:
```json
{
  "username": "someone",
  "email": "someone@external.test-company.com",
  "company": {
    "group": [
      "developer",
      "inetOrgPerson",
      "externals"
    ]
  }
}
``` 
Then we could define the required role as follows.
#### Using environment variables
```bash
export rules='[{"type": "disallow", "path": "^(/)$", "userinfo": {"company.group": "^externals$"}]'
```
#### Using configuration file
```yml
rules:
- type: "disallow"
  path: "^(/)$"
  userinfo:
    company.group: "^externals$"
```
#### Using arguments
```bash
go run main.go \
--rules='{"type": "disallow", "path": "^(/)$", "userinfo": {"company.group": "^externals$"}'
```