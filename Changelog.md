# Changelog

## Unreleased

- Fix `bud run` after installing from Github

  When you build the binaries with --trimpath, it removes the import paths stored within the binary. This is preferable because you don't want to see my filepaths within the binary, but it also means that the downloaded binary was no longer able to find where the standard library packages were located. I fixed this by calling `go env` on boot. This adds about 7ms overhead on boot, so I'd like to find a way to do this without spawning, but we'll leave that as an exercise for the reader :)

  Another problem we encountered was that the runtime was missing some of the necessary embedded assets. I've removed them from .gitignore to fix the problem. Longer-term I think it makes sense to use bindata for this case to turn them into Go files that can be ignored in development but are built into the binary for production.

## v0.0.3

- Bud is finally open source on Github!

  I'm thrilled to finally share my side project with everyone! I first started working on this while in lockdown in Berlin on April 20th, 2019. A co-worker suggested I have a look at the Laracast videos about Laravel. I was just blown away by how productive you can be in Laravel.

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
