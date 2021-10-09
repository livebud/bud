/**
 * Imports
 */

import { ParsedUrlQuery } from "querystring"
import { parse as parseQS } from "../qs"

/**
 * URL type definition
 */

export type URL = {
  auth: string
  hash: string
  host: string
  hostname: string
  href: string
  pathname: string
  protocol: string
  search: string
  port: number
  query: ParsedUrlQuery
}

/**
 * Create a cached element
 */

/**
 * Parse a url into it's components
 *
 * TODO: use the native browser library if it exists
 * TODO: memoize the URLs
 */

export function parse(url: string): URL {
  // TODO: create a url/parse/server.ts & move these out
  const a = document.createElement("a")
  const location = window.location
  a.href = url

  // parse the query
  const q = a.search.slice(1)
  const query = q ? parseQS(q) : {}

  // handle auth
  const auth: string[] = []
  if (a.username) {
    auth.push(a.username)
    // cant have a password without a username
    if (a.password) {
      auth.push(a.password)
    }
  }
  // return the URL
  return {
    auth: auth.join(":"),
    hash: trimSlash(a.hash),
    host: a.host || location.host,
    hostname: a.hostname || location.hostname,
    href: trimSlash(a.href),
    protocol:
      !a.protocol || a.protocol === ":" ? location.protocol : a.protocol,
    pathname: trimSlash(
      a.pathname.charAt(0) !== "/" ? "/" + a.pathname : a.pathname
    ),
    port:
      a.port === "0" || a.port === ""
        ? port(a.protocol) || parseInt(location.port, 10)
        : parseInt(a.port, 10),
    query: query,
    search: trimSlash(a.search),
  }
}

/**
 * Return default port for `protocol`.
 */

function port(protocol: string): number {
  switch (protocol) {
    case "http:":
      return 80
    case "https:":
      return 443
    default:
      return 0
  }
}

/**
 * Clean path by stripping subsequent "//"'s. Without this
 * the user must be careful when to use "/" or not, which leads
 * to bad UX.
 */

function trimSlash(path: string): string {
  return path === "/" ? path : path.replace(/\/+$/g, "")
}
