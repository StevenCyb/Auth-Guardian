# Auth-Guardian

<table><tr>
<td>
Auth-Guardian is an <b>auth</b>entication and rules based <b>auth</b>orization reverse proxy.
Authentication is provided via generic OAuth and generic SAML.
Authorization (not implemented yet) should be rule-based and as flexible as possible.
</td>
<td>
  <img src="doc/media/animation.gif">
</td>
</tr></table>

The documentation is located [here](doc/DOC.md), changelog [here](doc/CHANGELOG.md) and the contribution guidelines [here](doc/CONTRIBUTING.md).

![Overview](doc/media/overview.jpg)

## Roadmap
### 0.4.8
- Implement tests
  - Authentication
    - LDAP
      - Wrong credentials
      - Correct credentials + Check if correct path and query
    - SAML 
      - Wrong credentials
      - Correct credentials + Check if correct path and query
    - OAuth
      - Wrong credentials
      - Correct credentials + Check if correct path and query
  - Rules
    - Whitelist
    - Required
    - Disallow
- Add tests to makefile (local test) (create new test cert saml_mock.crt and saml_mock.key)
### Release Alpha-0.4.8 
### 0.4.9
- Create docker image building (harden)
- Create pipeline to build image
### 0.4.10
- Add docker compose example
- Add helm example
### Release Beta-0.5.0
### Some improvements...