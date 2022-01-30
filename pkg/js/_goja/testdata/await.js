var bud = (() => {
  var __defProp = Object.defineProperty
  var __markAsModule = (target) =>
    __defProp(target, "__esModule", { value: true })
  var __export = (target, all) => {
    __markAsModule(target)
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true })
  }

  // ssr:./bud/view/_ssr.js
  var ssr_exports = {}
  __export(ssr_exports, {
    render: () => render,
  })

  // ssr_runtime:./bud/view/_ssr_runtime.ts
  function renderHTML(input) {
    if (!input.view) {
      return {
        status: 404,
        headers: {},
        body: fallback(new Error('Missing page "' + input.route + '"')),
      }
    }
    return input.view({ props: input.props, context: input.context })
  }
  function fallback(err) {
    return `fallback error: ${err.message}`
  }

  // svelte_runtime:./bud/view/_svelte.ts
  function createView(view) {
    view.layout = view.layout || defaultLayout
    return function ({ props, context }) {
      const page = view.page.render(props)
      let css = page.css.code
      let html = page.html
      let head = page.head
      const hydrate = JSON.stringify(props)
      const layout = view.layout.render(props, {
        head: function () {
          return `
          ${head}
          <style>#bud{}${css}</style>
          <script id="bud_props" type="text/template" defer>${hydrate}<\/script>
          <script type="module" src="${view.client}" defer><\/script>
        `
        },
        default: function () {
          return '<div id="bud_target">' + html + "</div>"
        },
      })
      html = layout.html.replace("#bud{}", layout.css.code)
      return {
        status: 200,
        headers: {
          "Content-Type": "text/html",
        },
        body: html,
      }
    }
  }
  var defaultLayout = {
    render(props, slots) {
      return {
        css: {
          code: "",
        },
        head: "",
        html: `
        <html>
          <head>${slots.head(props)}</head>
          <body>${slots.default(props)}</body>
        </html>
      `,
      }
    },
  }

  // node_modules/svelte/internal/index.mjs
  function noop() {}
  function is_promise(value) {
    return (
      value && typeof value === "object" && typeof value.then === "function"
    )
  }
  function run(fn) {
    return fn()
  }
  function blank_object() {
    return Object.create(null)
  }
  function run_all(fns) {
    fns.forEach(run)
  }
  function is_function(thing) {
    return typeof thing === "function"
  }
  function is_empty(obj) {
    return Object.keys(obj).length === 0
  }
  var tasks = new Set()
  var active_docs = new Set()
  var current_component
  function set_current_component(component) {
    current_component = component
  }
  var resolved_promise = Promise.resolve()
  var seen_callbacks = new Set()
  var outroing = new Set()
  var globals =
    typeof window !== "undefined"
      ? window
      : typeof globalThis !== "undefined"
      ? globalThis
      : global
  var boolean_attributes = new Set([
    "allowfullscreen",
    "allowpaymentrequest",
    "async",
    "autofocus",
    "autoplay",
    "checked",
    "controls",
    "default",
    "defer",
    "disabled",
    "formnovalidate",
    "hidden",
    "ismap",
    "loop",
    "multiple",
    "muted",
    "nomodule",
    "novalidate",
    "open",
    "playsinline",
    "readonly",
    "required",
    "reversed",
    "selected",
  ])
  var escaped = {
    '"': "&quot;",
    "'": "&#39;",
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
  }
  function escape(html) {
    return String(html).replace(/["'&<>]/g, (match) => escaped[match])
  }
  var on_destroy
  function create_ssr_component(fn) {
    function $$render(result, props, bindings, slots, context) {
      const parent_component = current_component
      const $$ = {
        on_destroy,
        context: new Map(
          parent_component ? parent_component.$$.context : context || []
        ),
        on_mount: [],
        before_update: [],
        after_update: [],
        callbacks: blank_object(),
      }
      set_current_component({ $$ })
      const html = fn(result, props, bindings, slots)
      set_current_component(parent_component)
      return html
    }
    return {
      render: (props = {}, { $$slots = {}, context = new Map() } = {}) => {
        on_destroy = []
        const result = { title: "", head: "", css: new Set() }
        const html = $$render(result, props, {}, $$slots, context)
        run_all(on_destroy)
        return {
          html,
          css: {
            code: Array.from(result.css)
              .map((css) => css.code)
              .join("\n"),
            map: null,
          },
          head: result.title + result.head,
        }
      },
      $$render,
    }
  }
  function destroy_component(component, detaching) {
    const $$ = component.$$
    if ($$.fragment !== null) {
      run_all($$.on_destroy)
      $$.fragment && $$.fragment.d(detaching)
      $$.on_destroy = $$.fragment = null
      $$.ctx = []
    }
  }
  var SvelteElement
  if (typeof HTMLElement === "function") {
    SvelteElement = class extends HTMLElement {
      constructor() {
        super()
        this.attachShadow({ mode: "open" })
      }
      connectedCallback() {
        const { on_mount } = this.$$
        this.$$.on_disconnect = on_mount.map(run).filter(is_function)
        for (const key in this.$$.slotted) {
          this.appendChild(this.$$.slotted[key])
        }
      }
      attributeChangedCallback(attr, _oldValue, newValue) {
        this[attr] = newValue
      }
      disconnectedCallback() {
        run_all(this.$$.on_disconnect)
      }
      $destroy() {
        destroy_component(this, 1)
        this.$destroy = noop
      }
      $on(type, callback) {
        const callbacks =
          this.$$.callbacks[type] || (this.$$.callbacks[type] = [])
        callbacks.push(callback)
        return () => {
          const index = callbacks.indexOf(callback)
          if (index !== -1) callbacks.splice(index, 1)
        }
      }
      $set($$props) {
        if (this.$$set && !is_empty($$props)) {
          this.$$.skip_bound = true
          this.$$set($$props)
          this.$$.skip_bound = false
        }
      }
    }
  }

  // view/index.svelte
  var View = create_ssr_component(($$result, $$props, $$bindings, slots) => {
    let promise = fetch("http://127.0.0.1:55748").then((res) => res.text())
    return `<div>${(function (__value) {
      if (is_promise(__value)) {
        __value.then(null, noop)
        return `
                                        Loading...
                                `
      }
      return (function (value) {
        return `
                                        response: ${escape(value)}
                                `
      })(__value)
    })(promise)}</div>`
  })
  var view_default = View

  // svelte:./bud/view/index.svelte
  var view_default2 = createView({
    page: view_default,
    frames: [],
    client: "/bud/view/_index.svelte",
  })

  // ssr:./bud/view/_ssr.js
  var views = {}
  views["/"] = view_default2
  function render(route, props, context) {
    const view = views[route]
    if (!view) {
      return JSON.stringify({
        status: 404,
      })
    }
    return JSON.stringify(
      renderHTML({
        context,
        props,
        route,
        view,
      })
    )
  }
  return ssr_exports
})()
