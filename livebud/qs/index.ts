/**
 * Imports
 */

import { ParsedUrlQuery } from "querystring"

/**
 * Simple query string parser.
 *
 * @param {String} query The query string that needs to be parsed.
 * @returns {Query}
 * @api public
 */

export function parse(query: string): ParsedUrlQuery {
  if (!query) return {}
  const parser = /([^=?&]+)=?([^&]*)/g
  const result: ParsedUrlQuery = {}
  let part

  // Little nifty parsing hack, leverage the fact that RegExp.exec
  // increments the lastIndex property so we can continue executing
  // this loop until we've parsed all results.
  while ((part = parser.exec(query))) {
    /** @type {string|boolean} */
    const val = decodeComponent(part[2])
    const key = decodeComponent(part[1])

    // support arrays (item=a&item=b)
    let existing = result[key]
    if (typeof existing !== "undefined") {
      result[key] = ([] as string[]).concat(existing, val)
      continue
    }
    // Add query to result
    result[key] = val
  }

  return result
}

/**
 * Decode the URI component
 */

function decodeComponent(input: string): string {
  return decodeURIComponent(input.replace(/\+/g, " "))
}
