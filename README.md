go-sdk
======

[![Build Status](https://circleci.com/gh/blend/go-sdk.svg?style=shield)](https://circleci.com/gh/blend/go-sdk)
[![GoDoc](https://godoc.org/github.com/blend/go-sdk?status.svg)](https://godoc.org/github.com/blend/go-sdk)

`go-sdk` is our core library of packages. These packages can be composed to create anything from CLIs to fully featured web apps.

The general philosophy is to provide loosely coupled libraries that can be composed as a suite of tools, vs. a `do it all` framework.

# Requirements

This repository requires golang version 1.16+ to be installed.

To run tests, it is required that you have Docker installed, with docker-compose, or have postgres running locally.

# Addtional CLI Tools

We also provide the following CLI tools to help with development that leverage some of these packages:

- `cmd/ask` : securely input secrets and output to a file to be read by templates.
- `cmd/copyright` : injects and verifies copyright headers are present in files.
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

# Contributing

We currently don't accept PRs as this repository at this time, but feel free to log issues and we'll address when we can.

# Code Style Notes

## The "Options" Pattern

- The "Options" pattern is a variadic set of arguments as functions that mutate the returned object, typically in the constructor.
	- This lets callers add their own mutators to be used in constructors.
	- It also lets callers establish a set of "default" options that can be combined / overridden later.
	- It also reduces the amount of code required in the `go-sdk` repo itself.

## Other General Notes

- Where possible, follow the [golang proverbs](https://go-proverbs.github.io/).
- Make the zero value useful. Some situations require pointers, and are noted exceptions.
- Export all fields unless strictly internal state and would *never* be set by calling code.
- Where possible, packages should export configuration objects that can be used to create the core types of that package. Those configuration objects should be readable from both JSON and YAML.
- Anything that can return an error, should. Anything that needs to return a single value (but would return an error) should panic on that error and should be prefixed by `Must...`.
- Minimize dependencies between packages as much as possible; add external dependencies with *extreme* care.
	- Notable exceptions include:
		- airbrake
		- sentry
		- lib/pq
		- aws sdk
		- datadog

# Version Management

We follow calendar versioning. What that means in practice:

> [major].[year][month][day].[patch]

Major version changes are changes that break backwards-compatibility. A breaking change is defined as a change that would cause code written against the current major version having a build failure of any type. That's even for a trivial find and replace. Once you merge a new api or name for an object after CR, that's it. Major changes are rolled up into a release, at most, once-per-year.

Minor version changes are additions to exported types, values, or functions that don't cause break backwards compatibility. These types of changes are rolled into a new minor version at an approximate monthly cadence.

Patch versions is the used to distinguish different versions from the same calendar day.

The current version is stored in the `VERSION` file at the root of the package.

Currently we support 1 major version branches of the go-sdk.

- v1.xxxxxxxx.x

All other major versions are deprecated. The previous semver scheme is deprecated.


## Version Release Cycle

We will not be releasing more than 1 major version in any calendar year, and major versions will recieve at least 2 years of support before being deprecated.

# Bugs and Feature Requests

Please use the [issues page](https://github.com/blend/go-sdk/issues) to report any bugs or request new features be added to the SDK. We welcome contributions to any issue reported therein, including those you report. Please ping us on the issue itself for things like access requests.

# Maintainers

This repo is maintained by a core group of [Blend](https://github.com/blend) employees, including:
- [Addie Bendory](https://github.com/addie)
- [Alex Kaplan](https://github.com/KaplanAlex)
- [Chinmaya Sharma](https://github.com/chinmayasharma)
- [Igor Morozov](https://github.com/igmor)
- [Justin Fan](https://github.com/Justinfan827)
- [Kevin Wang](https://github.com/xorkevin)
- [Michael Turner](https://github.com/mat285)


