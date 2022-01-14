# Contributing

We'd love to accept your contributions to this project! There are just a few guidelines you need to follow.

## Bugs

Bug reports should be opened up as [issues](https://help.github.com/en/github/managing-your-work-on-github/about-issues) on the [go-vela/community](https://github.com/go-vela/community) repository!

## Feature Requests

Feature Requests should be opened up as [issues](https://help.github.com/en/github/managing-your-work-on-github/about-issues) on the [go-vela/community](https://github.com/go-vela/community) repository!

## Pull Requests

**NOTE: We recommend you start by opening a new issue describing the bug or feature you're intending to fix. Even if you think it's relatively minor, it's helpful to know what people are working on.**

We are always open to new PRs! You can follow the below guide for learning how you can contribute to the project!

## Getting Started

### Prerequisites

* [Review the GoLang package structure we use](https://github.com/golang-standards/project-layout) - ensure code is organized to our standards
* [Review the commit guide we follow](https://chris.beams.io/posts/git-commit/#seven-rules) - ensure your commits follow our standards
* [Docker](https://docs.docker.com/install/) - building block for local development
* [Docker Compose](https://docs.docker.com/compose/install/) - start up local development
* [Verdaccio](https://verdaccio.org/) - start up npm registry
* [Golang](https://golang.org/dl/) - for source code and [dependency management](https://github.com/golang/go/wiki/Modules)
* [Make](https://www.gnu.org/software/make/) - start up local development

### Setup

* [Fork](/fork) this repository

* Clone this repository to your workstation:

```bash
# Clone the project
git clone git@github.com:go-vela/vela-npm.git $HOME/go-vela/vela-npm
```

* Navigate to the repository code:

```bash
# Change into the project directory
cd $HOME/go-vela/vela-npm
```

* Point the original code at your fork:

```bash
# Add a remote branch pointing to your fork
git remote add fork https://github.com/your_fork/vela-npm
```

### Running Locally

* Navigate to the repository code:

```bash
# Change into the project directory
cd $HOME/go-vela/vela-npm
```

* Build the repository code:

```bash
# Build the code with `make`
make build
```

* Run the repository code:

```bash
# Run the code with `make`
make run
```

### Development

* Navigate to the repository code:

```bash
# Change into the project directory
cd $HOME/go-vela/vela-npm
```

* Make a test user for your local Verdaccio instance, this will create a .env file with the access token in it

```bash
# username: testuser
# password: testpass
# email: test@test.com
make user
```

* Write your code
  - Please be sure to [follow our commit rules](https://chris.beams.io/posts/git-commit/#seven-rules)

* Test your code against the provided example:

```bash
# Build code and test plugin against local registry
make test-e2e
```

* Ensure your code meets the project standards:

```bash
# Clean the code with `make`
make clean
```

* Push to your fork:

```bash
# Push your code up to your fork
git push fork main
```

* Open a pull request. Thank you for your contribution!

### Testing A Different Project

Sometimes it's helpful to test out publishing npm packages locally before trying in a CI environment. This repo is also helpful in facilitating that.

* Start Verdaccio
```bash
make registry
```

* Create a user
```bash
make user
```

* Build local plugin image
```bash
make build
```

* Clone your npm project repository

```bash
git clone my-project
```

* Get your project ready for publishing

```bash
cd my-project
npm install --production
# run build scripts if needed
npm run build
```

* From plugin repo, run local plugin against your project, where $PROJECT_PATH is the directory path to your npm project

```bash
PROJECT_PATH=/User/helpful-contributor/my-project make docker-run
```
