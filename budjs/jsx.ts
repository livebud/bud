import { HydrateInput } from "./"
import ReactDOM from "react-dom"
import React from "react"

export default function createView(input: HydrateInput) {
  let component = React.createElement(input.page, input.props)
  for (let frame of input.frames) {
    component = React.createElement(frame, input.props, component)
  }
  ReactDOM.hydrate(component, input.target)
}
