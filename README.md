# Bud

The Fullstack Go Framework for Prolific Web Developers.

Bud writes the boring code for you, helping you launch your website faster.

## Video Demo

Watch a video demonstrating how to build a minimal HN clone in 15 minutes with Bud.

[![](https://user-images.githubusercontent.com/170299/168361927-9165c2f9-55d4-4fa0-a53e-966028a79b39.png)](https://www.youtube.com/watch?v=LoypcRqn-xA)

## Documentation

Read [the documentation](https://denim-cub-301.notion.site/Hey-Bud-4d81622cc49942f9917c5033e5205c69) to learn how to get started with Bud.

# Installing Bud

Bud ships as a single binary that runs on Linux and Mac. You can follow along for Windows support in [this issue](https://github.com/livebud/bud/issues/7).

The easiest way to get started is by copying and pasting the command below in your terminal:

```diff
$ curl -sf curl https://raw.githubusercontent.com/livebud/bud/main/install.sh | sh
```

This script will download the right binary for your operating system and move the binary to the right location in your `$PATH`.

Confirm that you've installed Bud by typing `bud` in your terminal.

```bash
bud -h
```

You should see the following:

```bash
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

# Requirements

The following software is required to use Bud.

- Node v14+
   
   This is a temporary requirement that we plan to remove in [v0.3](https://github.com/livebud/bud/discussions/21)   
   
- Go v1.16+
   
   Bud relies heavily on `io/fs` and will take advantage of generics in the future, so while Go v1.16 will work, we suggest running Go v1.18+ if you can.

# Your First Project

With bud installed, you can now scaffold a new project:

```bash
$ bud create hello
$ cd hello
```

The create command will scaffold everything you need to get started with bud.

```bash
$ ls
go.mod  node_modules/  package-lock.json  package.json
```

... which is not very much by the way! Unlike most other fullstack frameworks, Bud starts out very minimal. As you add dependencies, Bud will generate all the boring code to glue your app together. Let's see this in action.

Start the development server with `bud run`:

```bash
$ bud run
| Listening on http://127.0.0.1:3000
```

Click on the link to open the browser. You'll be greeted with bud's welcome page.

Congrats! You're running your first web server with Bud. The welcome server is your jumping off point to learn more about the framework.

![CleanShot 2022-05-12 at 22.00.19@2x.png](https://denim-cub-301.notion.site/image/https%3A%2F%2Fs3-us-west-2.amazonaws.com%2Fsecure.notion-static.com%2Fdb7f750b-a699-4117-ac07-303124e5d2f4%2FCleanShot_2022-05-12_at_22.00.192x.png?table=block&id=9488d91f-b72d-4c6d-9ce0-358c31f7f964&spaceId=faf0f409-6e25-40a4-871e-3b311037350f&width=2000&userId=&cache=v2)

Check out the [Hacker News demo](https://www.youtube.com/watch?v=LoypcRqn-xA), [read the documentation](https://denim-cub-301.notion.site/Hey-Bud-4d81622cc49942f9917c5033e5205c69#156ea69b8d044bacb65fc2897f3e52b8), [schedule a quick call](https://cal.com/d/5iDN6qWZCLGMipwhFqhY3i/30min) or go on your own adventure. The only limit is your imagination.

## Contributing

Please refer to the [Contributing Guide](./Contributing.md) to learn how to run and test Bud locally.
