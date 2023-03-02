import type { create_ssr_component } from 'svelte/internal'
type Component = ReturnType<typeof create_ssr_component>

type View = {
  Component: Component
  key: string
}


type Pages = {
  [key: string]: Page
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
    const { Component, key, frames, layout } = this.state

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
      const { html: layoutHTML } = layout.Component.render(props[layout.key] || {}, {
        // context: new Map(Object.entries(page.layout?.context || {})),
        '$$slots': {
          default: () => `<div id="bud_target">${html}</div>`,
          head: () => heads.join("\n"),
          style: () => `<style>\n\t${styles.reverse().join("\n\t")}\n</style>`,
        }
      })
      html = layoutHTML
    }

    return html
  }
}


export class Viewer {
  constructor(private readonly pages: Pages) { }

  render(key: string, props: Props | null) {
    const page = this.pages[key]
    if (!page) {
      throw new Error(`svelte: unknown page ${key}`)
    }
    return page.render(props)
  }
}

