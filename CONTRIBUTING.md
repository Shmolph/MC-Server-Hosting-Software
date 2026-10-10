# Contributing to Shmolph Cloud

Welcome to the Shmolph Cloud project! We are excited to have you contribute to our open-source project. This guide will help you get started with setting up your development environment, understanding our coding standards, and making your first or next contribution.

Please also read our [Code of Conduct](https://github.com/Shmolph/MC-Server-Hosting-Software).

## Getting started

To start contributing to Shmolph Cloud, you need to have a basic understanding of the following:
- Backend: Go
- Frontend: Svelte & npm
- Git & GitHub

## Backend dev Environment Setup

### Pre-reqs

You will need the [latest go version](https://go.dev/dl/) (currently 1.27.2) and [postgresql](https://www.postgresql.org/download/) (You can have it locally installed or in a docker container) installed.

Optionally you can install a tool like Hoppscotch to interface with   

### How-to

1. Fork the repository and clone your fork.
2. Install backend dependencies: `cd backend/src && go mod download`
3. Setup your postgres install with a database and a user.
4. Edit the `database.go` file at line 12 with your correct credentials to connect to your database
5. Compile the project and run it `go build . -o shmolphhosting.exe && ./shmolphhosting.exe`

As an IDE, we recommend Visual Studio Code or Jetbrains GoLand for backend dev..

## Coding Standards

We currently don't have any code checkers but later we will use go [golangci-lint](https://github.com/golangci/golangci-lint) to enforce certain code styles and standards.

See `code-standards.md` for the full rundown of tooling, tests, and codebase conventions.

## Making Contributions

From your forked repository, make your own changes on your own branch. **Do not make changes directly to main!**

When you are ready, you can submit a pull request to the Shmolph Cloud repository. If you are still working on your pull request or need help with something, make sure to mark it as a **Draft**.

Also, please make sure that your pull requests are as targeted and simple as possible and don't do a hundred things at a time. If you want to add/change/fix 5 different things, you should make 5 different pull requests.

## Before opening a pull request

Make sure your code looks good :)

## Contributor License Agreement

All contributors must sign our Contributor License Agreement. The CLA Assistant bot will ask you to sign it on your first pull request; posting the comment it asks for completes the signature.

## Translations

Crowdin coming soon™

## Code Review Process

Your pull request will then be reviewed by the maintainers.

Once you have an approval from a maintainer, another will merge it once it’s confirmed. Depending on the pull request size, this process can take multiple days.

## Community and Support

- Help: Send a message through [matrix](https://element.io/en) to @shmolph:gekohosting167.online
- Bugs: [https://github.com/Shmolph/MC-Server-Hosting-Software/issues](https://github.com/Shmolph/MC-Server-Hosting-Software/issues)
- Features: Forums coming soon™
- Security vulnerabilities: See our [security policy](https://github.com/Shmolph/MC-Server-Hosting-Software/security).
