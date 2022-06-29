# Changelog

Get the latest release of Bud by running the following in your terminal:

```sh
curl -sf https://raw.githubusercontent.com/livebud/bud/main/install.sh | sh
```

## v0.1.11

- Fixed regression when using `bud new controller` (#173)

## v0.1.10

- Fix regression when running `bud create` outside of a `$GOPATH` (#167)
- Fix cache clear before code generation to ensure there's no stale generated code. Prior to v0.1.10, we were clearing part of the cache not all of it (#168)
- Allow explicit versions to be installed with `curl` by setting the `VERSION` environment variable (#168)

## v0.1.9

- Better hot reload DX with `bud run` (#131) thanks to @012e

  Prior to v0.1.9, whenever a file change occurs it would rebuild and print a ready message if successful or an error message if unsuccessful. The ready messages would often clutter the terminal over time.

  <a href="https://share.cleanshot.com/XCdkEnOr8BlCNgV035pu"><video src='https://share.cleanshot.com/XCdkEnOr8BlCNgV035pu/download' width='100%'/></a>

- Large internal refactor (#133). The goals of this refactor:

  1. Make it easier to understand. The double generated binary building was something that often confused me.
  2. Make it easier to contribute. I'm so impressed with the contributions so far, with this refactor it should be even easier.
  3. Make it faster during development. The slowest step in the build process is running `go build`. We now only run `go build` once on boot, not twice.

  Learn more details [in this comment](https://github.com/livebud/bud/pull/133#issuecomment-1166371510). This PR concludes the work necessary to release [v0.2](https://github.com/livebud/bud/discussions/18).

- Support glob embeds (#150) thanks to @vito

  Build caching now understands embedded globs like `// go:embed *.sql`. You'll no longer get stale builds when changing a file within an embedded glob.

- Improved Dockerfile in contributing (#140) thanks to @wheinze

  The Dockerfile now supports passing Node and Go versions to build a custom container. It also uses a smaller base image.

## v0.1.8

- Support `func(w, r)` controller actions (#147) (thanks @vito!)

  This release adds a highly-requested feature (atleast by me!) where you can now drop down to using the vanilla `http.HandlerFunc` signature.

  This is a useful escape hatch for webhooks, streaming and other complex use cases.

  ```go
  package webhooks

  type Controller struct {
    GH *github.Client
  }

  func (c *Controller) Index() string {
    return "GitHub webhook service!"
  }

  // Create handles incoming webhooks
  func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
    // Respond to the webhook using a vanilla HTTP handler function!
  }
  ```

- Replace redirect with case-insensitive routing (#142) (thanks @vito!)

  Prior to v0.1.8, if you defined the route `/bud`, but a user visited `/BUD`,
  they would be redirected to `/bud`. This was originally done for SEO purposes to prevent different casing from appearing as separate pages.

  However, this had an unfortunate side-effect in that you couldn't use parameters with mixed casing (e.g. base64 encoding).

  In v0.1.8, we changed this so that URL routing is now case insensitive, so `/BUD` will run the `/bud` action. This doesn't address the SEO issue, but that will be a [follow-up task](https://github.com/livebud/bud/pull/142#issuecomment-1159824008) for a later time.

## v0.1.7

- Ensure alignment between CLI and runtime round 2 (#128)

  In #126, we missed the case where the module wasn't downloaded yet. This version fixes that.

## v0.1.6

- Ensure alignment between CLI and runtime (#126)

  In v0.1.5, we had a breaking change in the runtime. If you had an existing project and upgraded the CLI to v0.1.5, but you were still using the v0.1.4 runtime in go.mod, you'd encounter an error. This change automatically updates your go.mod to align with the CLI that's building it. Fixes: #125.

## v0.1.5

This release focuses on paying down some technical debt that was accumulated prior to the release. It's part of the [v0.2](https://github.com/livebud/bud/discussions/18) plan.

- Rename `bud run [--port=<address>]` to `bud run [--listen=<address>]`

  This **breaking change** addresses the confusion discussed in https://github.com/livebud/bud/discussions/42.

- Rename `bud tool v8 client` to `bud tool v8 serve`

  This **breaking change** gives a better name for what the command does, listen for eval requests, evaluate the javascript and return the response.

- 204 No Content improvements (thanks @theEyeD!)

  Now when controller don't have a return value or return a nil error, they return `204 No Content`. For example:

  ```go
  package users
  type Controller struct {}
  func (c *Controller) Create() error {
    return nil
  }
  ```

  ```sh
  $ curl -X POST /users
  HTTP/1.1 204 No Content
  Content-Length: 0
  ```

- Dedupe and group watch events (#123)

  This reduces the number of events the watcher triggers on Linux causing less builds to trigger at once. It also hopefully fixed an issue where "remove" events sometimes weren't triggering rebuilds on all platforms.

- Improve CI (thanks @wheinze!) (#118)

  Waldemar took a much needed pass over the CI. He cleaned up the file, fixed caching and extended the test matrix, so we can more confidently say what's required to use Bud.

- Refactor, test and simplify the compiler (#93)

  The compiler wasn't really tested prior to v0.1.4. Issue that would crop up would be discovered due to E2E tests. In v0.1.5, the compiler was refactored to be testable, then tested. It was also simplified to reduce the number of ways you could use it programmatically. This makes it less likely a test will pass but the end-to-end usage will fail.

- Extend the error chain (#100)

  Errors were being lost while merging the filesystems for plugins. Now errors that occur in generators will be extended through the merged filesystem, so it's easier to see when there are issues in the generators

## v0.1.4

- Add support for custom actions (thanks @theEyeD!)

  This release adds support for defining custom actions on controllers that get mapped to GET requests. This feature closes: https://github.com/livebud/bud/pull/67.

  For example, given the following users controller in `controller/users/users.go`:

  ```go
  package users
  type Controller struct {}
  func (c *Controller) Deactivate() error { return nil }
  ```

  The `Deactivate` method would get mapped to `GET /users/deactivate`. Learn more in [the documentation](https://www.notion.so/Hey-Bud-4d81622cc49942f9917c5033e5205c69#cc536e505a7349e08c68a65aef266726).

- Speed up tests by using existing `GOMODCACHE`

  Some of the test suite used an empty `GOMODCACHE` to test plugin support. This turned added about a 1-minute overhead to those tests while dependencies were downloaded and the cache was populated.

  We now rely on real module fixtures: https://github.com/livebud/bud-test-plugin and https://github.com/livebud/bud-test-plugin.

- Add a version number to released assets

  This will make it easier to add Bud to other package managers like the Arch User Repository (AUR) for Arch Linux users. This feature fixes: https://github.com/livebud/bud/issues/52.

- Added a background section to the Readme (thanks @thepudds!)

  @thepudds re-wrote the Reddit post, simplifying and condensing it for the Readme. I always like when projects include a "why", so I really appreciate @thepudds adding this for Bud!

## v0.1.3

- Move build caching from $TMPDIR/bud/cache to $YOUR_APP/bud/cache

  This should fix https://github.com/livebud/bud/issues/27 issue on M1 Macs.

- Fallback to copying if renaming a directory fails

  This change fixes a `invalid cross-device link` error WSL users were encountering.

## v0.1.2

- Fix live reload for Windows Subsystem for Linux (WSL) users (thanks @theEyeD!)

## v0.1.1

- Improve the installation script (thanks @barelyhuman!)

  For cases where `/usr/local/bin` is not writable, the install script will now prompt you to escalate your privileges.

- Run Go's formatter on generated template files (thanks @codenoid!)

  Now the Go files generated into `bud/` will be formatted with the same formatter as `go fmt`.

- Fix `bud create` when working outside of `$GOPATH` (thanks @barelyhuman!)

  If you tried created a Bud project outside of `$GOPATH`, Bud would be unable to infer the Go module path and will prompt you to provide a module path. This prompt was being clobbered by NPM's progress bar. Now the prompt happens before running `npm install`.

- Add a [Contributor's Guide](./Contributing.md)

  Wrote an initial guide on how to contribute to Bud. Please have a look and let me know if you'd like to see anything else in that guide.

- Switch the default from 0.0.0.0 to 127.0.0.1 (#35)

  On the latest OSX you'll get a security prompt when you try binding to 0.0.0.0. On WSL, bud will just stall. This release switches the default to localhost (aka 127.0.0.1). Many thanks to @alecthomas, @theEyeD and @kevwan for helping me understand this issue better.

## v0.1.0

This release wraps up the [v0.1](https://github.com/livebud/bud/discussions/17) milestone 🎉

- Add initial darwin/arm64 (aka Apple M1) Support (please test!)
- Update the welcome page to include all the links.

I'm completely overwhelmed by all the love. As I write this, we're still on the front page of [Hacker News](https://news.ycombinator.com/item?id=31371340) (13 hours in!):

![HN](https://user-images.githubusercontent.com/170299/168419120-83faf09c-f5c8-406c-a98c-4c694fab5c2e.png)

And we're at the top of the [Go subreddit](https://www.reddit.com/r/golang).

![Reddit](https://user-images.githubusercontent.com/170299/168419117-a2b1a033-0f87-498d-8a8b-71d608c6dd57.png)

## v0.0.9

- Add support for `bud create .`
- Fix `bud create` for the development binary when outside the source tree.
- Reverse `bud new controller` from route:resource (e.g. `/:stories`) to resource:route (e.g. `stories:/`) to avoid confusion with URL params.
- Correct the "Ready" message

## v0.0.8

- Give the welcome page a fresh coat of paint.
- Support `bud version [key]` to make it easier to script getting the version of bud.

## v0.0.7

- Download go modules after scaffolding

  Hopefully fully fixing `bud create <dir>` this time. Bud now downloads the runtime modules upon create.

- Add some getting started instructions

  I've written up a bit of documentation on how to get started. There's a lot more to come here!

## v0.0.6

- Fixed `bud create`

  Now that bud is public and has versions, I did another pass over `bud create <dir>` to get it working as expected

- Add live reload to the welcome page

  In an effort to vanguish all page refreshes, needing to refresh the welcome page stuck out like a thorn. This release clips that thorn.

## v0.0.5

- Next attempt at fixing `bud run` after installing from Github

  It turns out .go files are also required by the runtime. I've unignored them, but will think of better ways to handle this in the future.

## v0.0.4

- Fix `bud run` after installing from Github

  When you build the binaries with --trimpath, it removes the import paths stored within the binary. This is preferable because you don't want to see my filepaths within the binary, but it also means that the downloaded binary was no longer able to find where the standard library packages were located. I fixed this by calling `go env` on boot. This adds about 7ms overhead on boot, so I'd like to find a way to do this without spawning, but we'll leave that as an exercise for the reader :)

  Another problem we encountered was that the runtime was missing some of the necessary embedded assets. I've removed them from .gitignore to fix the problem. Longer-term I think it makes sense to use bindata for this case to turn them into Go files that can be ignored in development but are built into the binary for production.

- Add back missing node_modules

  I overpruned the dependencies and that was causing failures in CI. I added them back in.

## v0.0.3

- Bud is finally open source on Github!

  I'm thrilled to finally share my side project with everyone! I first started working on this while in lockdown in Berlin on April 20th, 2020. A co-worker suggested I have a look at the Laracast videos about Laravel. I was just blown away by how productive you can be in Laravel.

  As a former solo developer, I place a lot of weight on having all the tools you need to build, launch and iterate on ideas as quickly as possible. Laravel provides a comprehensive and cohesive toolset for getting your ideas out there quickly.

  With Go being my preferred language of choice these days and a natural fit for building web backends, I started prototyping what a Laravel-like MVC framework would look like in Go.

  At this point, I just had the following goal for Bud:

  - Be just as productive as Laravel in a typed language like Go.

  I got the first version working about 6 months in and I tried building a blog from it. It fell flat. You needed to scaffold all these files just to get started. If you're coming from Rails or Laravel you may shrug, this is pretty normal. Unfortunately, I've been spoiled by the renaissance in frontend frameworks with Next.js. What I love Next is that it starts out barebones and every file you add incrementally enhances your web application. This keeps the complexity under control.

  With these newly discovered constraints, I started working on the next iteration of Bud.

  - Generate files only as you need them. Keep these generated files away from your application code.
  - Should feel like using a modern Javascript framework. This means it should work with modern frontend frameworks like Svelte and React, support live reload and have server-side rendering.

  With these new goals, the Bud you see today started to take shape. But along the way, I discovered a few more project goals:

  - The framework should be extensible from Day 1. Bud is too ambitious for one person. We're going to need an ambitious community behind this framework.
  - Bud should have great performance for the developer using Bud and for the person using websites built with Bud. We have an exciting journey ahead for both of these goals.
  - Bud should compile to a single binary. These days with platforms services like Heroku and Vercel, it's easy to not care about this, but I still cherish the idea that I can build a single binary that contains my entire web app and scp it up to tiny server.

  And this is the Bud you see before you. I have big plans for the framework and I sincerely hope you'll join me on this journey.

## v0.0.2

- Improve test caching

  The tests are slower than they should be. What was curious was that they didn't seem to be doing much during a lot of the test run.

  I dug into why and it turns out the test caching itself can be extremely slow. See this issue for more details: https://github.com/golang/go/issues/26562. It's slated to be fixed in Go 1.19, so I'm going to hold off on a fix for that.

  Additionally, the tests weren't ever being cached. If you'd like to read about the whole debugging saga that took a day to figure out, head over to this issue: https://github.com/golang/go/issues/26562. I've remedied this issue, but it's something to keep an eye on.

- Prep the build script

  I'm in the process of setting up `curl -sf https://github.com/livebud/bud/install.sh | sh`. In doing that, I'm ironing out the publishing pipeline and install script.

## v0.0.1

- Initial publish release
