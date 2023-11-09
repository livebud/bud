import { Component, JSX } from "preact"

export type DocumentProps = {
  common?: any
  script?: any
  head?: any
  style?: any
  page?: any
  props?: any
}

type DocumentContext = {
  _documentProps: DocumentProps
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
      _documentProps: {},
    }
  }
  render(): JSX.Element {
    return <></>
  }
}

export class Head extends Component {
  render() {
    return undefined
  }
}

export class Page extends Component {
  render() {
    return undefined
  }
}

export class Scripts extends Component {
  render() {
    return undefined
  }
}
