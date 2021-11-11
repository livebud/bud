/**
 * Imports
 */

import assert from "internal/assert"
import { parse, URL } from "."

/**
 * Test runner
 */

describe("parse", () => {
  tests().forEach((test) => {
    it(`${test.input}`, async () => {
      assert.deepEqual(parse(test.input), test.expect)
    })
  })
})

/**
 * Test type definition
 */

type Test = {
  input: string
  expect: URL
}

/**
 * Tests
 */

function tests(): Test[] {
  return [
    // {
    //   input: '127.0.0.1:57371',
    //   expect: {
    //     auth: '',
    //     protocol: 'http:',
    //     host: '127.0.0.1:57412',
    //     port: 57412,
    //     hostname: '127.0.0.1',
    //     hash: '',
    //     search: '',
    //     query: {},
    //     pathname: '/',
    //     href: '%3127.0.0.1:57371', // TODO
    //   },
    // },
    {
      input: "http://google.com/foo/bar",
      expect: {
        auth: "",
        hash: "",
        host: "google.com",
        hostname: "google.com",
        href: "http://google.com/foo/bar",
        pathname: "/foo/bar",
        port: 80,
        protocol: "http:",
        query: {},
        search: "",
      },
    },
    {
      input: "http://google.com:3000/foo/bar",
      expect: {
        auth: "",
        hash: "",
        host: "google.com:3000",
        hostname: "google.com",
        href: "http://google.com:3000/foo/bar",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "http:",
        query: {},
        search: "",
      },
    },
    {
      input: "https://google.com/foo/bar",
      expect: {
        auth: "",
        hash: "",
        host: "google.com",
        hostname: "google.com",
        href: "https://google.com/foo/bar",
        pathname: "/foo/bar",
        port: 443,
        protocol: "https:",
        query: {},
        search: "",
      },
    },
    {
      input: "http://google.com:80/foo/bar",
      expect: {
        auth: "",
        hash: "",
        host: "google.com",
        hostname: "google.com",
        href: "http://google.com/foo/bar",
        pathname: "/foo/bar",
        port: 80,
        protocol: "http:",
        query: {},
        search: "",
      },
    },
    {
      // this one is a bit weird
      input: "https://google.com:3000/foo/bar",
      expect: {
        auth: "",
        hash: "",
        host: "google.com:3000",
        hostname: "google.com",
        href: "https://google.com:3000/foo/bar",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "https:",
        query: {},
        search: "",
      },
    },
    {
      input: "http://google.com:3000/foo/bar?name=tobi",
      expect: {
        auth: "",
        hash: "",
        host: "google.com:3000",
        hostname: "google.com",
        href: "http://google.com:3000/foo/bar?name=tobi",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "http:",
        query: { name: "tobi" },
        search: "?name=tobi",
      },
    },
    {
      input: "http://google.com:3000/foo/bar#something",
      expect: {
        auth: "",
        hash: "#something",
        host: "google.com:3000",
        hostname: "google.com",
        href: "http://google.com:3000/foo/bar#something",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "http:",
        query: {},
        search: "",
      },
    },
    {
      input: "http://google.com:3000/foo/bar?name=tobi&foo=bar#something",
      expect: {
        auth: "",
        hash: "#something",
        host: "google.com:3000",
        hostname: "google.com",
        href: "http://google.com:3000/foo/bar?name=tobi&foo=bar#something",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "http:",
        query: { foo: "bar", name: "tobi" },
        search: "?name=tobi&foo=bar",
      },
    },
    {
      input: "http://a:b@google.com:3000/foo/bar?name=tobi&foo=bar#something",
      expect: {
        auth: "a:b",
        hash: "#something",
        host: "google.com:3000",
        hostname: "google.com",
        href: "http://a:b@google.com:3000/foo/bar?name=tobi&foo=bar#something",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "http:",
        query: { foo: "bar", name: "tobi" },
        search: "?name=tobi&foo=bar",
      },
    },
    {
      input: "http://a@google.com:3000/foo/bar?name=tobi&foo=bar#something",
      expect: {
        auth: "a",
        hash: "#something",
        host: "google.com:3000",
        hostname: "google.com",
        href: "http://a@google.com:3000/foo/bar?name=tobi&foo=bar#something",
        pathname: "/foo/bar",
        port: 3000,
        protocol: "http:",
        query: { foo: "bar", name: "tobi" },
        search: "?name=tobi&foo=bar",
      },
    },
  ]
}
