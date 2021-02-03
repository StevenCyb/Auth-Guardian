# Auth-Guardian
Auth-Guardian is an **auth**entication and rules based **auth**orization reverse proxy.
Authentication is provided via generic OAuth and generic SAML.
Authorization (not implemented yet) should be rule-based and as flexible as possible.

The documentation is located [here](doc/doc.md), changelog [here](doc/CHANGELOG.md) and the contribution guidelines [here](doc/contributing.md).

![Overview](doc/media/overview.jpg)

## Roadmap
### 0.3.0
Add LDAP support (nee to rename project in this case since SAML != oauth)
### Release Alpha-0.3.0
### 0.4.0
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
### Release Alpha-0.4.0
### 0.4.1
- Create docker image
- Create helm chart
### Release Alpha-0.4.1
### 0.5.1
- Implement tests