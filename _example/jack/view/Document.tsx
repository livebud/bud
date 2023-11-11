import { Component, JSX, Fragment, h } from "preact"

type VNode = {
  type: string
  props: Record<string, string>
  // children: (VNode | string)[]
}

export type DocumentProps = {
  script?: string
  heads: VNode[]
  style?: string
  page?: string
}

type DocumentContext = {
  _docProps: DocumentProps
}

export default class Document extends Component<DocumentProps> {
  static Head: typeof Head
  static Page: typeof Page
  static Scripts: typeof Scripts
  getChildContext(): DocumentContext {
    return {
      _docProps: {
        heads: this.props.heads,
        page: this.props.page,
      },
    }
  }
  render(): JSX.Element {
    return <></>
  }
}

export class Head extends Component {
  render() {
    const docProps = this.context._docProps || {}
    const heads = docProps.heads || []
    return (
      <head>
        {(this.props.children || ([] as any)).concat(
          heads.map((node: VNode) => h(node.type, node.props))
        )}
      </head>
    )
  }
}

export class Page extends Component {
  render() {
    const docProps = this.context._docProps || {}
    return (
      <div
        id="bud#target"
        dangerouslySetInnerHTML={{ __html: docProps.page }}
      ></div>
    )
  }
}

export class Scripts extends Component {
  render() {
    return []
  }
}
