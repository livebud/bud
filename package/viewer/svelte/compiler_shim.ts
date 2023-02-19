/**
 * Shim for getting the svelte compiler to run in a V8 isolate.
 */

// URL shim for the browser
// TODO: properly shim URL
export class URL {
  constructor(url: string) {
    console.log(url)
  }
}

// TODO: properly shim performance.now()
export const self = {
  performance: {
    now(): number {
      return 0
    },
  },
}

// In development mode when compiling for the browser we hit this codepath:
// https://github.com/Rich-Harris/magic-string/blob/8f666889136ac2580356e48610b3ac95c276191e/src/SourceMap.js#L3-L10
// Since we're running in a V8 isolate, we don't have a window or a Buffer.
// TODO: shim btoa properly
export const window = {
  btoa: (data: string): string => {
    return ""
  },
}
