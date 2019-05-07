go-sdk
======

[![Build Status](https://circleci.com/gh/blend/go-sdk.svg?style=shield)](https://circleci.com/gh/blend/go-sdk)
[![GoDoc](https://godoc.org/github.com/blend/go-sdk?status.svg)](https://godoc.org/github.com/blend/go-sdk)

`go-sdk` is our core library of packages. These packages can be composed to create anything from CLIs to fully featured web apps.

The general philosophy is to provide loosely coupled libraries that can be composed as a suite of tools, vs. a `do it all` framework.

# Addtional CLI Tools

We also provide the following CLI tools to help with development that leverage some of these packages:

- `cmd/ask` : securely input secrets and output to a file to be read by templates.
- `cmd/cover` : allows for project level coverage reporting and enforcement.
- `cmd/job` : run a command on a cron schedule; useful for writing jobs as kubernetes pods.
- `cmd/profanity` : profanity rules checking (i.e. fail on grep match).
- `cmd/recover` : recover crashed processes (to be used when debugging panics).
- `cmd/semver` : semver manipulation and validation.
- `cmd/shamir` : securely partition secrets using shamir's sharing scheme.
- `cmd/template` : commandline template generation using golang `text/template`.

# Repository Organization Notes

- The repository is organized into composible packages. They are designed to be as loosely coupled as possible, but are all in a single repo to facilitate breaking change management.
- Any commandline programs should live under `cmd/**` with the core library code in a top level package. The CLI should just be the bare minimum to run the library from the cli.

# Code Style Notes

- Where possible, follow the [golang proverbs](https://go-proverbs.github.io/).
- Where possible, make the zero value useful. If you do need to internalize a field, make sure it's accessible with an accessor in the form of the uppercase name of the field, i.e. field `foo` would have an accessor `Foo()`.
- Where possible, packages should export configuration objects that can be used to create the core types of that package. Those configuration objects should be readable from both JSON and YAML.
- Anything that can return an error, should. Anything that needs to return a single value (but would return an error) should panic on that error and should be prefixed by `Must...`.
- Minimize dependencies between packages as much as possible; add external dependencies with *extreme* care.

# Version Management

Generally we follow semantic versioning. What that means in practice:

> [major].[minor].[patch]

Major version changes are changes that break backwards-compatibility. A breaking change is defined as a change that would cause code written against the current major version having a build failure of any type. That's even for a trivial find and replace. Once you merge a new api or name for an object after CR, that's it. Major changes are rolled up into a release, at most, once-per-year.

Minor version changes are additions to exported types, values, or functions that don't cause break backwards compatibility. These types of changes are rolled into a new minor version at an approximate monthly cadence.

Patch versions are bugfixes and improvements made without changing the current set of exports. They can be cut at any time.

The current version is stored in the `VERSION` file at the root of the package.

Currently we support 2 major version branches of the go-sdk.

- v1.0 is scheduled for deprecation/end-of-support in Q4 2019
- v2.0 is the current version and will be supported until at least May 2021

Another version v3.0 is in development as the master branch and will likely be released in late Q2/early Q3 2019.

## Version Release Cycle

Our version release cycle has changed. Following v3.0's release later this year, we will not be releasing more than 1 major version in any calendar year, and major versions will recieve at least 2 years of support before being deprecated.

## To increment the local version

Patch:
> make increment-patch

Minor:
> make increment-minor

Major:
> make increment-major

# Bugs and Feature Requests

Please use the [issues page](https://github.com/blend/go-sdk/issues) to report any bugs or request new features be added to the SDK. We welcome contributions to any issue reported therein, including those you report. Please ping us on the issue itself for things like access requests.