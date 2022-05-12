# Bud

The Fullstack Go Framework for Prolific Web Developers.

Bud writes the boring code for you, helping you launch your website faster.

## Installation

Run the following command in your terminal to install `bud` into your `$PATH`:

```sh
curl -sf curl https://raw.githubusercontent.com/livebud/bud/main/install.sh | sh
```

You can confirm the binary is working as expected with `bud -h`:

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

## Create a new project

With bud installed, you can now scaffold a new project. Let's create a minimal hackernews clone:

```
$ bud create hackernews
$ cd hackernews
```

The create command will scaffold everything you need to get started with bud.

```sh
$ ls
go.mod  node_modules/  package-lock.json  package.json
```

... which is not very much by the way! Unlike most other fullstack frameworks, bud starts out very minimal. As you add dependencies, bud generates all the boring code to glue your app together. Let's see this in action.

Start the development server with `bud run`:

```sh
$ bud run
| Listening on http://0.0.0.0:3000
```

Click on the link to open the browser. You'll be greeted with bud's welcome page.

Congrats! You're running your first web server with Bud.

## Going Further

Use the quick links on the welcome page to learn more about how to go further with Bud.
