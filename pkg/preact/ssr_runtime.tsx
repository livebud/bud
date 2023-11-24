import renderToString from "preact-render-to-string"
import { h, FunctionComponent, Component, JSX, VNode } from "preact"

const heads: JSX.Element[] = []

class HeadProvider extends Component<any> {
  getChildContext() {
    return { heads: heads }
  }

  render() {
    return this.props.children
  }
}

type Options = {
  liveUrl: string
}

export function renderView(
  path: string,
  View: FunctionComponent,
  props: any,
  options: Options
): string {
  const html = renderToString(
    <HeadProvider>
      <View {...props} />
    </HeadProvider>
  )
  const out = {
    html: html,
    heads: heads
      .map(renderToJson)
      .concat(propScript(props), clientScript(path)),
  }
  if (options.liveUrl) {
    out.heads.push(liveScript(options.liveUrl))
  }
  return JSON.stringify(out)
}

function renderToJson(el: JSX.Element): VNode<any> {
  if (typeof el.type === "function") {
    throw new Error("rendering components inside head is not supported yet")
  }
  const props = el.props
  if (el.props.children) {
    // children can be undefined, an array, a component, or just a string
    props.children = Array.isArray(el.props.children)
      ? el.props.children.map(renderToJson)
      : el.props.children.type !== undefined
      ? renderToJson(el.props.children)
      : el.props.children
  }
  return {
    type: el.type,
    props: props,
    key: el.key,
  }
}

function propScript(props: any): VNode<any> {
  return {
    type: "script",
    props: {
      id: "bud#props",
      type: "text/template",
      defer: true,
      dangerouslySetInnerHTML: { __html: JSON.stringify(props) },
    },
    key: undefined,
  }
}

function clientScript(path: any): VNode<any> {
  return {
    type: "script",
    props: {
      src: `/view/${path}.js`,
      type: "text/javascript",
      defer: true,
    },
    key: undefined,
  }
}

function liveScript(liveUrl: string): VNode<any> {
  return {
    type: "script",
    props: {
      src: liveUrl,
      type: "text/javascript",
      defer: true,
    },
    key: undefined,
  }
}
