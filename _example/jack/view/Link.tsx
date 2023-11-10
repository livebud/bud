import { Component } from "preact"

type LinkProps = {
  href: string
  prefetch?: boolean
}

export default class Link extends Component<LinkProps> {
  render() {
    return <a href={this.props.href}>{this.props.children}</a>
  }
}
