# Auth-Guardian
Auth-Guardian is an **auth**entication and rules based **auth**orization reverse proxy.
Authentication is provided via generic OAuth and generic SAML.
Authorization (not implemented yet) should be rule-based and as flexible as possible.

The documentation is located [here](doc/doc.md), changelog [here](doc/CHANGELOG.md) and the contribution guidelines [here](doc/contributing.md).

![Overview](doc/media/overview.jpg)

## Roadmap
### 0.4.2
Add rule middleware with allow rules:
```yaml
allow-rule:
  # For ...
  # Condition
  roles: []string
  userinfo: []string
  query_parameter: {key: regex}
  json_body_parameter: {key: regex}
```
Example https://github.com/louketo/louketo-proxy/blob/master/docs/user-guide.md
### 0.4.3
Add rule middleware with not allow only rules:
```yaml
disallow-rule:
  # For ...
  # Condition
  roles: []string
  userinfo: []string
  query_parameter: {key: regex}
  json_body_parameter: {key: regex}
```
### Release Alpha-0.4.3
### 0.4.1
- Create docker image
- Create helm chart
### Release Alpha-0.4.1
### 0.5.1
- Implement tests