# Contributing to Auth-Guardian
First of all - thanks for your contribution!

The following is a set of guidelines for contributing to GMG on GitHub.
These are guidelines, not rules. 
Try to provide the necessary information so that your pull request can be understood.

## Table Of Contents
* [Reporting Bugs](#reporting-bugs)
* [Suggesting Enhancements](#suggesting-enhancements) 
* [Code Styleguide](#code-styleguide)
* [Write Tests](#write-tests)
* [Update the Documentation and Changelog](#update-the-documentation-and-changelog)
* [Versioning](#versioning)
* [Git Commit Messages](#git-commit-message)

### Reporting Bugs
Please check existing bugs before reporting, sometimes bugs are already reported. 
If you found a closed report which seems to be similar, then link it in your report. Try to add as many details as possible and use [these template](templates/bug_report.md).

### Suggesting Enhancements
Please check the [roadmap](../README.md#roadmap) before suggesting an enhancement as you might find out that you don't need to create one.
When you are creating an enhancement suggestion, please include as many details as possible and use [this template](/.github/ISSUE_TEMPLATE/feature_request.md).

### Code Styleguide
- Comment new functions with the name and a one-line description `// func-/struct-/interface-Name and your description`
- Try to not use external packages if possible 
- Add additional comments in the line before and not at the end!

### Write Tests
Nobody likes to write test but writing tests is important to determine errors which could occur through changes.
Changes that effect the code and are testable should have tests.
Even if tests exist, check if some of them needs to be updated to.

### Update the Documentation and Changelog
Most changes require an adjustment of the documentation. 
Check the full documentation for possible points that would need to be adjusted together with your contribution.
In addition, the changes should be included in the changelog.

### Versioning
Let us keep the versioning simple. 
Currently the project is in the beta phase. 
Therefore, the versions are only incremented on `Minor` for feature implementation and `Path/Bugfix/Hotfix` for parts of a feature.
After the beta phase you should increase the versions depending on the degree of compatibility:

| Major |.| Minor |.| Path/Bugfix/Hotfix |
|-------|-|-------|-|--------------------|
| 1     |.| 2     |.| 3                  |

### Git Commit Messages
Writing a commit message is not that hard ;)
- Write your message in present tense ("Change, Add" not "Changed, Added")
- Use the imperative mood ("Move" not "Moves")
- Do not exceed 72 characters on a single line
