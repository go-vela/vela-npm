# Contributing

## Getting Started

We'd love to accept your contributions to this project! If you are a first time contributor, please review our [Contributing Guidelines](https://go-vela.github.io/docs/community/contributing_guidelines/) before proceeding.

### Prerequisites

* [Review the commit guide we follow](https://chris.beams.io/posts/git-commit/#seven-rules) - ensure your commits follow our standards
* Review our [style guide](https://go-vela.github.io/docs/community/contributing_guidelines/#style-guide) to ensure your code is clean and consistent.
* [Docker](https://docs.docker.com/install/) - building block for local development
* [Docker Compose](https://docs.docker.com/compose/install/) - start up local development
* [Verdaccio](https://verdaccio.org/) - start up npm registry
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

* Write your code and tests to implement the changes you desire.

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

* Make sure to follow our [PR process](https://go-vela.github.io/docs/community/contributing_guidelines/#development-workflow) when opening a pull request

Thank you for your contribution!

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
