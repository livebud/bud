# Changelog

Get the latest release of Bud by running the following in your terminal:

```sh
curl -sf https://raw.githubusercontent.com/livebud/bud/main/install.sh | sh
```

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
