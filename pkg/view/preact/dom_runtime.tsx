import { render, FunctionComponent, Component } from "preact"

const heads: any[] = []

class HeadProvider extends Component<any> {
  getChildContext() {
    return { heads: heads }
  }

  render() {
    return this.props.children
  }
}

const target = document.getElementById("bud") || document.body
const props = JSON.parse(
  document.getElementById("bud#props")?.textContent || "{}"
)

export function renderView(View: FunctionComponent): void {
  render(
    <HeadProvider>
      <View {...props} />
    </HeadProvider>,
    target
  )
}
