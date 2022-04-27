# Changelog

## v0.0.2

- Improve test caching

  The tests are slower than they should be. What was curious was that they didn't seem to be doing much during a lot of the test run.

  I dug into why and it turns out the test caching itself can be extremely slow. See this issue for more details: https://github.com/golang/go/issues/26562. It's slated to be fixed in Go 1.19, so I'm going to hold off on a fix for that.

  Additionally, the tests weren't ever being cached. If you'd like to read about the whole debugging saga that took a day to figure out, head over to this issue: https://github.com/golang/go/issues/26562. I've remedied this issue, but it's something to keep an eye on.

- Prep the build script

  I'm in the process of setting up `curl -sf https://github.com/livebud/bud/install.sh | sh`. In doing that, I'm ironing out the publishing pipeline and install script.

## v0.0.1

- Initial publish release
