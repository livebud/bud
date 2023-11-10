import { Component, JSX, Fragment, h } from "preact"

type VNode = {
  name: string
  attrs: Record<string, string>
  children: (VNode | string)[]
}

export type DocumentProps = {
  script?: string
  head?: VNode[]
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
        head: this.props.head,
        page: this.props.page,
      },
    }
  }
  render(): JSX.Element {
    return <></>
  }
}

export class Head extends Component {
  getName(node: VNode): typeof Fragment | string {
    if (node.name === "Fragment") {
      return Fragment
    }
    return node.name
  }

  hydrate(node: VNode): JSX.Element {
    return h(
      node.name,
      node.attrs || {},
      (node.children || []).map((child) =>
        typeof child === "string" ? child : this.hydrate(child)
      )
    )
  }

  render() {
    const docProps = this.context._docProps || {}
    return (
      <head>
        {(this.props.children || ([] as any)).concat(
          docProps.head.map((head: VNode) => this.hydrate(head))
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
