type View = {
  page: any
  frames: any[]
  layout: any
  error?: any
  client: string
}

// TODO:
// - Test custom layouts
// - Support frames
// - Support default errors
// - Support custom errors
export function createView(view: View) {
  view.layout = view.layout || defaultLayout
  return function ({ props, context }) {
    // Wrap what's passed from the server as props. This gives a consistent
    // exported variable that can be any value:
    // - For index.svelte, that would be a list of resources, e.g. props = [...]
    // - For show.svelte, that would be an single resources, e.g. props = {...}
    props = { props }

    const page = view.page.render(props)
    let css = page.css.code
    let html = page.html
    let head = page.head
    // Render the layout
    const hydrate = JSON.stringify(props)
    const layout = view.layout.render(props, {
      head: function () {
        return `
          ${head}
          <style>#bud{}${css}</style>
          <script id="bud_props" type="text/template" defer>${hydrate}</script>
          <script type="module" src="${view.client}" defer></script>
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

const defaultLayout = {
  render(props, slots) {
    return {
      css: {
        code: "",
      },
      head: "",
      html: `
        <!doctype html>
        <html>
          <head>
            <meta charset="utf-8"/>
            ${slots.head(props)}
          </head>
          <body>${slots.default(props)}</body>
        </html>
      `,
    }
  },
}
