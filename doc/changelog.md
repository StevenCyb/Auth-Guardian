# Changelog
All notable changes to this project will be documented in this file.

## [0.1.1] - 2021-02-01
### Added
- Options to forward userinfo and/or access-token
- Test option to run application in test mode (run a service behind which mirrors request)
### Changes
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