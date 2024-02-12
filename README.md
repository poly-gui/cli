# Poly CLI

[What is Poly?](https://polygui.org)

This CLI provides commands to create and manage Poly projects. It is the official way to get started with creating Poly
applications.

## Installation

The CLI is available as pre-built binaries on the [GitHub Releases page](https://github.com/poly-gui/cli/releases).

### Building The CLI

The CLI can also be built from the source code, which requires:

- [Go >= 1.20](https://go.dev/dl/)

First, clone this repo:

```
git clone https://github.com/poly-gui/cli.git
```

This clones the repo into a directory called "cli". If the name is too generic, feel free to clone the repo into a
directory with a different name.

Change into the repo, then build and install the binary:

```
go install poly-cli/cmd/poly
```

The `poly` command should now be installed and ready to use. Make sure `GOBIN`
is in PATH, which defaults to `$(go env GOPATH)/bin`. If not, add:

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

or if `GOBIN` is set:

```
export PATH="$PATH:$(go env GOBIN)/bin"
```

to your path.

For more information on `GOPATH` and `GOBIN`, please consult
the [official documentation](https://go.dev/doc/install/source#gopath).

## Usage

### ```poly generate```

The `generate` command creates a new Poly project directory in the current working directory by default. The following
flags are available:

|    Flags    |                                                                      Description                                                                      |
|:-----------:|:-----------------------------------------------------------------------------------------------------------------------------------------------------:|
| `--output`  |                      The directory in which the project directory should be created. Defaults to the current working directory.                       |
|  `--name`   |                                                  The name of the application. Defaults to "PolyApp"                                                   |
| `--package` | The package name/bundle identifier/application ID of the application. This typically uses the reverse domain name notation. Defaults to "org.polygui" |

### Examples

```
poly generate --name=MyNewApp --package=org.my
```

This generates a new Poly app called "MyNewApp" that has "org.my" as the package identifier. A "MyNewApp" directory is created in the current working directory.
