# Auth-Guardian
Auth-Guardian is an **auth**entication and rules based **auth**orization reverse proxy.
Authentication is provided via generic OAuth.
Authorization (not implemented yet) should be rule-based and as flexible as possible.

The documentation is located [here](doc/doc.md) and the contribution guidelines [here](doc/contributing.md).

![Overview](doc/media/overview.jpg)

## Roadmap
### 0.2.0
- Endpoint: Health
- Endpoint: Metrics
- Config-Option for which infos should be forwarded to upstream (userinfo and/or access_token)
### 0.3.0
Add rule middleware (strategy needs to be planed first) | oauthmiddleware -> rulemiddleware -> upstream
```yaml
rules:
  - allow: bool (def false)
    path: string (def "/*")
    roles: [string]
    method: string (def "*")
    query_parameter: {key: value}
    json_body_parameter: {key: value}
```
### 0.4.0
Add SAML support (nee to rename project in this case since SAML != oauth)
### 0.5.0
Add LDAP support (nee to rename project in this case since SAML != oauth)
### 0.5.1
Create docker image
### 0.5.2
Create helm chart
##### Go public
### 0.5.3
- Implement tests