# Auth-Guardian

<table><tr>
<td>
Auth-Guardian is an <b>auth</b>entication and rules based <b>auth</b>orization reverse proxy.
Authentication can be provided via generic OAuth, generic SAML or LDAP.
The rule-based authorization is designed to be as flexible as possible.
</td>
<td>
  <img src="doc/media/animation.gif">
</td>
</tr></table>

![Overview](doc/media/overview.jpg)

## Links:
* [Documentation](/doc/doc.md)
* [Changelog](/doc/changelog.md) 
* [Contribution Guidelines](/doc/contributing.md).

## Used Third-Party Library
| Package                       | Version                            | Used for                                     |
| ----------------------------- |------------------------------------| -------------------------------------------- |
| github.com/crewjam/saml       | v0.4.5                             | the SAML authentication middleware           |
| github.com/shaj13/go-guardian | v1.5.11                            | the LDAP authentication middleware           |
| golang.org/x/crypto           | v0.0.0-20200622213623-75b288015ac9 | the certificate handling with SAML           |
| golang.org/x/oauth2           | v0.0.0-20210126194326-f9ce19ea3013 | the OAuth authentication middleware          |
| gopkg.in/asn1-ber.v1          | v1.0.0-20181015200546-f715ec2f112d | the mock of LDAP IDP                         |
| gopkg.in/square/go-jose.v2    | v2.5.1                             | the parsing of JWT token in OAuth middleware |
| gopkg.in/yaml.v2              | v2.4.0                             | the parsing of the configuration file        |

## Roadmap
### 0.4.9
- Create docker image building (harden)
- Create pipeline to build image
### 0.4.10
- Add docker compose example
- Add helm example
- Change logo
### Release Beta-0.5.0
### Some improvements...
