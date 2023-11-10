import { Component } from "preact"

type LinkProps = {
  href: string
  prefetch?: boolean
}

export default class Link extends Component<LinkProps> {
  render() {
    // TODO: fill in
    return this.props.children
  }
}
