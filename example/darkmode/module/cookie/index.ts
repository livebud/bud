/**
 * Export all by default
 * Based on: https://github.com/component/cookie
 */

export default { all, get, set }

/**
 * Cookie type
 */

type Cookies = {
  [cookie: string]: string
}

/**
 * Set options
 */

type Options = {
  maxage?: number
  expires?: Date
  domain?: string
  path?: string
  secure?: boolean
}

/**
 * Set cookie `name` to `value`.
 */

function set(name: string, value: string | null, options?: Options) {
  options = options || {}
  var str = encode(name) + '=' + encode(String(value))

  if (null == value) {
    options.maxage = -1
  }

  if (options.maxage) {
    options.expires = new Date(+new Date() + options.maxage)
  }

  if (options.path) str += '; path=' + options.path
  if (options.domain) str += '; domain=' + options.domain
  if (options.expires) str += '; expires=' + options.expires.toUTCString()
  if (options.secure) str += '; secure'

  document.cookie = str
}

/**
 * Return all cookies.
 *
 * This is isomorphic and may be called
 * from the server-side though it will
 * return nothing.
 */

function all(): Cookies {
  var str
  try {
    str = document.cookie
  } catch (err) {
    console.log(err)
    return {}
  }
  return parse(str)
}

/**
 * Get cookie `name`.
 */

function get(name: string): string | undefined {
  return all()[name]
}

/**
 * Parse cookie `str`.
 */

function parse(str: string): Cookies {
  var obj = <Cookies>{}
  var pairs = str.split(/ *; */)
  for (var i = 0; i < pairs.length; ++i) {
    var pair = pairs[i]
    var eqidx = pair.indexOf('=')
    if (eqidx === -1) {
      eqidx = pair.length
    }
    var name = decode(pair.substr(0, eqidx))
    // +1 because we don't want the =
    var value = decode(pair.substr(eqidx + 1))
    if (!name || !value) {
      continue
    }
    obj[name] = value
  }
  return obj
}

/**
 * Encode.
 */

function encode(value: string): string | undefined {
  try {
    return encodeURIComponent(value)
  } catch (e) {
    return
  }
}

/**
 * Decode.
 */

function decode(value: string): string | undefined {
  try {
    return decodeURIComponent(value)
  } catch (e) {
    return
  }
}