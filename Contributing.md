# Bud's Contributor Guide

Thank you for your interest in helping to make Bud better! Bud is and forever will be a driven by the community. We depend on volunteers lending their time and expertise to make Bud better for everyone.

## Requirements

- Node 14
- Go 1.18

## Setting up Bud for Development

Run the following commands to download and run Bud locally:

```sh
git clone https://github.com/livebud/bud
cd bud
make install # fresh installs take a few minutes because of V8
go run main.go
```

After running `go run main.go`, you should see the following:

```
  Usage:
    bud [flags] [command]

  Flags:
    -C, --chdir  Change the working directory

  Commands:
    build    build the production server
    create   create a new project
    run      run the development server
    tool     extra tools
    version  Show package versions
```

If you run into any problems, please [open an issue](https://github.com/livebud/bud/issues/new).

## Developing Bud with a Project

Now that you have Bud running locally, you can use the `-C, --chdir` functionality to test Bud against different projects.

You can use one of the projects in `example/`

```sh
# Run the development server on the example/hn application
go run main.go -C example/hn run

# Build a binary for the example/hn application
go run main.go -C example/hn build
```

Or you can create your own project:

```sh
# Scaffold a new hello application
go run main.go create hello

# Run the development server for the hello application
go run main.go -C hello run

# Build a binary for the hello application
go run main.go -C hello build
```

## Running Tests

You can use the following commands to run all the tests:

```sh
make test
```

## Publishing a new Version

To publish a new version of Bud, you will need your a clean local main branch that matches the remote main branch. You'll also need to have commit access on both Github and publish rights on NPM.

You can run then run:

1. Write the `Changelog.md`. You can use `git log v0.1.0..HEAD` to see the commits since the last release.
2. Bump the version in `version.txt`.
3. Run `make publish`.
