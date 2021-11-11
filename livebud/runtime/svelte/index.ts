import { HydrateInput } from ".."

// TODO:
// - Support frames
// - Handle errors
export default function createView(input: HydrateInput) {
  const component = new input.page({
    target: input.target,
    props: input.props,
    hydrate: true,
  })
}
