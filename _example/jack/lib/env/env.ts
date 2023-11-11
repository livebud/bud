// This module is provided by ESBuild
// TODO: make this a valid empty file
// @ts-expect-error
import ENV from 'bud/env.json'

interface Input {
  [key: string]: string | typeof String | number | typeof Number | boolean | typeof Boolean
}

type Output<T> = {
  [K in keyof T]:
  T[K] extends string | typeof String
  ? string
  : T[K] extends number | typeof Number
  ? number
  : T[K] extends boolean | typeof Boolean
  ? boolean
  : never
}


export default function env<Schema extends Input>(schema: Schema): Output<Schema> {
  // TODO: type output
  const out: any = {}
  for (let key in schema) {
    const value = ENV[key]
    // ENV[key] is undefined or an empty string
    if (typeof value === 'undefined' || value === '') {
      if (baseType(value)) {
        throw new Error(`environment variable "${key}" is not defined`)
      }
      const defaultValue = schema[key]
      if (typeof defaultValue != 'string' && typeof defaultValue != 'number' && typeof defaultValue != 'boolean') {
        throw new Error(`environment variable "${key}" must be either a string, number or boolean`)
      }
      out[key] = defaultValue
      continue
    }
    // we got something in ENV[key]
    switch (baseType(schema[key]) || defaultType(schema[key])) {
      case 'string':
        out[key] = value
        break
      case 'boolean':
        out[key] = value === 'false' ? false : Boolean(value)
        break
      case 'number':
        const n = parseInt(value, 10)
        if (isNaN(n)) {
          throw new Error(`environment variable "${key}" is not a number`)
        }
        out[key] = n
        break
      default:
        throw new Error(`environment variable "${key}" must be either a string, boolean or number`)
    }
  }
  return out
}


/**
 * Get the type
 */

function baseType(val: unknown): string | undefined {
  switch (val) {
    case String:
      return 'string'
    case Boolean:
      return 'boolean'
    case Number:
      return 'number'
    default:
      return undefined
  }
}

/**
 * Check if it's a default value
 */

function defaultType(val: unknown) {
  switch (typeof val) {
    case 'string':
      return 'string'
    case 'number':
      return 'number'
    case 'boolean':
      return 'boolean'
  }
}
