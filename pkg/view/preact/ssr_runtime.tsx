import renderToString from "preact-render-to-string"
import { FunctionComponent, Component, JSX } from "preact"

const heads: JSX.Element[] = []

class HeadProvider extends Component<any> {
  getChildContext() {
    return { heads: heads }
  }

  render() {
    return this.props.children
  }
}

export function renderView(View: FunctionComponent, props: any): string {
  const html = renderToString(
    <HeadProvider>
      <View {...props} />
    </HeadProvider>
  )
  return JSON.stringify({
    html: html,
    // head: heads.map((head) => renderToJson(head)),
    head: heads.map((head) => renderToString(head)).join(""),
  })
}

type VNode = {}

function renderToJson(vnode: JSX.Element): VNode {
  console.log("vnode", JSON.stringify(vnode))
  return {}
}
