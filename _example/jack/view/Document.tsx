import { Component, JSX, Fragment } from "preact"

export type DocumentProps = {
  script?: string
  head?: string
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
  constructor() {
    super()
  }
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
  render() {
    const docProps = this.context._docProps || {}
    // TODO: merge with dangerouslySetInnerHTML={{ __html: docProps.head }}
    return <head>{this.props.children}</head>
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
