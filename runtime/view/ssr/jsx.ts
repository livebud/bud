import ReactSSR from "react-dom/server"
import React from "react"

type View = {
  page: any
  frames: any[]
  layout: any
  error?: any
  client: string
}

export function createView(view: View) {
  return function ({ props, context }) {
    let component = React.createElement(view.page, props, [])
    for (let frame of view.frames) {
      component = React.createElement(frame, props, component)
    }
    let component2 = React.createElement("div", { id: "bud_target" }, component)
    const layout = view.layout || defaultLayout
    let component3 = React.createElement(layout, props, component2)
    let html = ReactSSR.renderToString(component3)
    let inject = ""
    const hydrate = JSON.stringify(props)
    inject += `<script id="bud_props" type="text/template" defer>${hydrate}</script>`
    inject += `<script type="module" src="${view.client}" defer></script>`
    html = html.replace("</head>", inject + `</head>`)
    return {
      status: 200,
      headers: {
        "Content-Type": "text/html",
      },
      body: html,
    }
  }
}

function defaultLayout(props) {
  return React.createElement(
    "html",
    null,
    React.createElement("head", null),
    React.createElement("body", null, props.children)
  )
}
