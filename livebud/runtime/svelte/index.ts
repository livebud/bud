import { HydrateInput } from ".."

// TODO:
// - Support frames
// - Handle errors
export default function createView(input: HydrateInput) {
  if (input.target != null) {
    // TODO: for some reason Svelte isn't able to re-hydrate over itself during
    // a live reload. I wonder if they've figured this out in SvelteKit, but you
    // end up with a runtime error: "Cannot read properties of null (reading
    // 'removeChild')". Encountering this error will depend on your Svelte code.
    // For now, we'll clear the DOM in our target before hydrating.
    input.target.innerHTML = ""
  }
  new input.page({
    target: input.target,
    props: input.props,
    hydrate: true,
  })
}
