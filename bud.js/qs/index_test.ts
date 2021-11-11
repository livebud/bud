/**
 * Imports
 */

import { ParsedUrlQuery } from "querystring"
import assert from "internal/assert"
import { parse } from "../qs"

describe("qs/parse", () => {
  tests().forEach((test) => {
    it(test.title || test.input, () => {
      assert.deepEqual(parse(test.input), test.expected)
    })
  })
})

/**
 * Test type definition
 */

type Test = {
  title?: string
  input: string
  expected: ParsedUrlQuery
}

/**
 * Tests
 */

function tests(): Test[] {
  return [
    { input: "", expected: {}, title: "empty" },
    { input: "name&species", expected: { name: "", species: "" } },
    { input: "name=false", expected: { name: "false" } },
    {
      input: "name=tobi&species=ferret",
      expected: { name: "tobi", species: "ferret" },
    },
    {
      input: "?names=friends+and+family",
      expected: { names: "friends and family" },
    },
    {
      input: "?name=tobi&species=ferret",
      expected: { name: "tobi", species: "ferret" },
    },
    {
      input: "items=1&items=2&items=3&key=a",
      expected: { items: ["1", "2", "3"], key: "a" },
    },
  ]
}
