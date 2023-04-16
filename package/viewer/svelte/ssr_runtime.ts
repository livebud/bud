import type { create_ssr_component } from 'svelte/internal'

type Component = ReturnType<typeof create_ssr_component>

type View = {
  Component: Component
  key: string
  client: string
}

type Props = {
  [key: string]: unknown
}

type State = View & {
  layout?: View
  frames: View[]
  error?: View
}

export class Page {
  constructor(private readonly state: State) { }

  render(props: Props | null) {
    props = props === null ? {} : props
    const { Component, key, client, frames, layout } = this.state

    // Load the page component
    const styles: string[] = []
    let heads: string[] = []
    const { head, css, html: pageHTML } = Component.render(props[key] || {}, {
      // context: new Map(Object.entries(page.context || {}))
    })
    if (head.length > 0) {
      heads.push(head)
    }
    if (css.code.length > 0) {
      styles.push(css.code)
    }
    let html = pageHTML

    // Render the frames
    for (let frame of frames) {
      const { head, css, html: frameHTML } = frame.Component.render(props[frame.key] || {}, {
        // context: new Map(Object.entries(frame.context || {})),
        '$$slots': { default: () => html }
      })
      if (head.length > 0) {
        heads.push(head)
      }
      if (css.code.length > 0) {
        styles.push(css.code)
      }
      html = frameHTML
    }

    // Render the layout
    if (layout) {
      const layoutProps = props[layout.key] || {}
      // Don't pass layout props down to the client
      delete props[layout.key]
      const clientScript = `<script src="${client}" type="module" async defer></script>`
      const { head, css, html: layoutHTML } = layout.Component.render(layoutProps, {
        // context: new Map(Object.entries(page.layout?.context || {})),
        '$$slots': {
          default: () => `<div id="bud_target">${html}</div><script id="bud_state" type="text/template">${escape({ props })}</script>`,
          head: () => clientScript,
        }
      })
      if (head.length > 0) {
        heads.push(head)
      }
      if (css.code.length > 0) {
        styles.push(css.code)
      }
      // Add the styles to the head
      if (styles.length) {
        heads.push(`<style id="bud_style">\n\t${styles.reverse().join("\n\t")}\n</style>`)
      }
      // Replace static client script with all the heads, including the client
      // script and styles.
      heads.push(clientScript)
      html = layoutHTML.replace(clientScript, heads.reverse().join('\n'))
    }

    return html
  }
}

// Based on: https://github.com/mathiasbynens/jsesc
// `jsesc(props, { isScriptContext: true, json: true })`
function escape(props: any): any {
  return JSON.stringify(props)
    .replace(/<\/(script|style)/gi, '<\\/$1')
    .replace(/<!--/g, '\\u003C!--');
}
