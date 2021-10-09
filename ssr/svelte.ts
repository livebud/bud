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
        <html>
          <head>${slots.head(props)}</head>
          <body>${slots.default(props)}</body>
        </html>
      `,
    }
  },
}
