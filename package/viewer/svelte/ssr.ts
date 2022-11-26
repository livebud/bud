/**
 * ssr.ts is the runtime for rendering server-side pages
 */

import type { create_ssr_component } from 'svelte/internal'
type Component = ReturnType<typeof create_ssr_component>

type Input = {
  page: Component
  frames: Component[]
  layout?: Component
  error?: Component
}

type RenderParams = Parameters<Component["render"]>
type RenderProps = RenderParams[0]
type RenderOptions = Exclude<RenderParams[1], undefined>
type RenderResult = ReturnType<Component["render"]>

type View = {
  key: string
  path?: string
  props?: Record<string, any>
  context?: Record<string, any>
}

type Page = View & {
  frames: View[]
  layout?: View
  error?: View
  client?: string
}

export class Viewer {
  constructor(private readonly input: Input) { }

  render(page: Page) {
    const input = this.input
    const styles: string[] = []
    let heads: string[] = []
    const { head, css, html: pageHTML } = input.page.render(page.props, {
      context: new Map(Object.entries(page.context || {}))
    })
    if (head.length > 0) {
      heads.push(head)
    }
    if (css.code.length > 0) {
      styles.push(css.code)
    }
    let html = pageHTML
    for (let i = 0; i < input.frames.length; i++) {
      const { head, css, html: frameHTML } = input.frames[i].render(page.frames[i].props, {
        context: new Map(Object.entries(page.frames[i].context || {})),
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
    // Add the data and client-side script. Note: we add to the head with defer
    // with the assumption that we won't block rendering but the browser will
    // start downloading dependencies earlier. TODO: double-check this.
    heads = heads.reverse()
    if (page.client) {
      heads.push(`<script id="bud_data" type="text/data" defer>${JSON.stringify(page)}</script>`)
      heads.push(`<script src="${page.client}" defer></script>`)
    }

    const layout = input.layout || new defaultLayout()
    const { html: layoutHTML } = layout.render(page.layout?.props, {
      context: new Map(Object.entries(page.layout?.context || {})),
      '$$slots': {
        default: () => `<div id="bud_target">${html}</div>`,
        head: () => heads.join("\n"),
        style: () => `<style>\n\t${styles.reverse().join("\n\t")}\n</style>`,
      }
    })
    return layoutHTML
  }
}

class defaultLayout {
  render(_: RenderProps, options: RenderOptions): RenderResult {
    const slots = options.$$slots || {}
    const head = 'head' in slots ? slots['head']() : ''
    const style = 'style' in slots ? slots['style']() : ''
    const body = 'default' in slots ? slots['default']() : ''
    return {
      head: "",
      css: { code: "", map: undefined },
      html: `<!doctype html>
  <html>
    <head>
      <meta charset="utf-8" />
      ${head}
      ${style}
    </head>
    <body>${body}</body>
  </html>
`
    }
  }
}
