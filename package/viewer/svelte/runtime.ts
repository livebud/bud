import type { create_ssr_component } from 'svelte/internal'
type Component = ReturnType<typeof create_ssr_component>

type Entry = {
  main: Component
  frames: Component[]
  layout?: Component
  error?: Component
  scripts: string[]
}

type RenderParams = Parameters<Component["render"]>
type RenderProps = RenderParams[0]
type RenderOptions = Exclude<RenderParams[1], undefined>
type RenderResult = ReturnType<Component["render"]>

type Input = {
  props: RenderProps
  context?: RenderOptions["context"]
}

type PropMap = {
  main: Input
  frames: Input[]
  layout?: Input
  error?: Input
}

export default class View {
  constructor(private readonly entry: Entry) { }

  render(propMap: PropMap) {
    const entry = this.entry
    const heads: string[] = []
    const styles: string[] = []

    const { head, css, html: mainHTML } = entry.main.render(propMap.main.props, {
      context: propMap.main.context
    })
    if (head.length > 0) {
      heads.push(head)
    }
    if (css.code.length > 0) {
      styles.push(css.code)
    }
    let html = mainHTML
    for (let i = 0; i < entry.frames.length; i++) {
      const { head, css, html: frameHTML } = entry.frames[i].render(propMap.frames[i].props, {
        context: propMap.frames[i].context,
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
    const layout = entry.layout || new defaultLayout()
    const { html: layoutHTML } = layout.render(propMap.layout?.props, {
      context: propMap.layout?.context,
      '$$slots': {
        default: () => html,
        head: () => heads.reverse().join("\n"),
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
