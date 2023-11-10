import { Component } from "preact"

export default class Head extends Component {
  render() {
    this.context.heads.push(this.props.children)
    return undefined
  }
}
