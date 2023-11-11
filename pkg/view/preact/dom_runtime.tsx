import { render, FunctionComponent, Component } from "preact"

const heads: HTMLCollection = document.head.children || []

// TODO: support dynamically changing the <head> tag
class HeadProvider extends Component<any> {
  getChildContext() {
    return { heads: heads }
  }

  render() {
    return this.props.children
  }
}

const target = document.getElementById("bud") || document.body
const props = getProps(document.getElementById("bud#props"))
export function renderView(View: FunctionComponent): void {
  render(
    <HeadProvider>
      <View {...props} />
    </HeadProvider>,
    target
  )
}

function getProps(el: HTMLElement | null): Record<string, unknown> {
  if (!el) {
    return {}
  }
  const text = el.textContent
  if (!text) {
    return {}
  }
  try {
    return JSON.parse(text)
  } catch (e) {
    return {}
  }
}
