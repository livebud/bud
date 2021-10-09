import { compile as compileSvelte } from "svelte/compiler"

type Input = {
  code: string
  path: string
  target: "ssr" | "dom"
  dev: boolean
}

// Capitalized for Go
type Output =
  | {
      JS: string
      CSS: string
    }
  | {
      Error: {
        Path: string
        Name: string
        Message: string
        Stack?: string
      }
    }

// Compile svelte code
export function compile(input: Input): string {
  const { code, path, target, dev } = input
  const svelte = compileSvelte(code, {
    filename: path,
    generate: target,
    hydratable: true,
    format: "esm",
    dev: dev,
    css: dev,
  })
  return JSON.stringify({
    CSS: svelte.css.code,
    JS: svelte.js.code,
  } as Output)
}
