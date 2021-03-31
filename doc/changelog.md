# Changelog
All notable changes to this project will be documented in this file.

## [0.5.0-beta] - 2021-03-31
### Added
- Add examples "how to run it"
- Github-Action (CI) to build and publish latest-image to DockerHub
### Changed
- Now support config on configmap path

## [0.4.9] - 2021-03-27
### Added
- Github-Action (CI) to build and publish image to DockerHub
- Dockerfile to build image
- Github-Action (CI) to run test

## [0.4.8] - 2021-03-21
### Added
- Integration tests
- Check if argument definitions already set
### Changed
- Update readme: links, third-party package section + add LDAP hint
- Mocked servers now return server/listener so we can close them on test routines
- Move config override of mocking services to run functions
- Move server stuff to separated package
- Update config example
- Improve some logging
- If configuration file not exists - skip instead of throwing an error
- Improve rule parse types
- Improve flatdata parse types
- Mock services now override their configuration on initialization
- Some small improvements

## [0.4.7] - 2021-02-18
- Mock SAML IDP
- SAML authentication now support registration 
### Changed
- Fix copied word misspelling

## [0.4.6] - 2021-02-13
- Mock LDAP IDP
### Changed 
- Fixe and add some logging outputs
- Change order in roadmap

## [0.4.5] - 2021-02-10
### Added
- Mock OAuth IDP
- Logo and animation for README
### Changed
- Change order in roadmap
- Remove version argument
- Moved templates to github
- Rename postfix `yaml` to `yml`
- move test-service to mocked package
- rename `test-mode` argument to `mock-test-service`

## [0.4.4] - 2021-02-09
### Added
- Implementation of disallow rules for rule middleware
### Changed
- Some small improvements

## [0.4.3] - 2021-02-09
### Changed
- Rules are now configured by single argument `rule` and differentiated by `type` key 
- Update doc regarding rule changes
### Bugfix
- If no rules specified, return empty array and not nil to prevent parsing error

## [0.4.2] - 2021-02-09
### Added
- Required rule configuration
- Implement of required rules for rule middleware
- If forwarding token is enabled also send the claims
### Bugfix
- On custom StringArrayFlag type

## [0.4.1] - 2021-02-07
### Changed
- Config now support formatted rules 
- Update whitelist rule definition to use rules (more cumbersome but uniform)
- Some small improvements

## [0.4.0] - 2021-02-06
### Added
- New routs for test service (/style.css, /script.js, /favicon.ico, /(HTML))
- Whitelist rule middleware
- Authorization rule middleware (dummy)
### Changed
- Fix error logging for file config loading
- Config now support string arrays
- Test service route / moved to /mirror
- Change code examples type in doc
- Change overview graphic

## [0.3.0] - 2021-02-04
### Added
- Generic LDAP authentication
### Changeds
- Add debug log

## [0.2.0] - 2021-02-03 
### Added
- Configuration option for SAML
- Middleware for SAML authentication
### Changeds
- Use filename `.yaml` instead of `.yml`
- Remove old todo
- Upstream selection is now provided by provider
- Auth middleware usage is now generic 
- Add simple config validation for OAuth and SAML to detect targeted authentication mechanism
- `changelog.md` is now written as `CHANGELOG.md`

## [0.1.1] - 2021-02-01
### Added
- Options to forward userinfo and/or access-token
- Test option to run application in test mode (run a service behind which mirrors request)
### Changeds
- Roadmap re-organized 
- Error printing on oauth.go (did not use logging previously)
- Add changelog hint to contribution document 
- Move cookie.go, session.go and utils.go to package util

## [0.1.0] - 2021-01-31
### Added
- Basic structure of the software
- Logging handler to log in JSON or formatted string, and to log to a file
- Configuration handler with env, file and argumentation configuration options
- Generic OAuth middleware
- Upstream handler via reverse proxy lib and for cors without reverse proxy lib 
- README with description, links to doc, contribution and changelog and roadmap 
- This changelog file