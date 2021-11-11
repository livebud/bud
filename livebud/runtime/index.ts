import Hot from "./hot"

export type HydrateInput<Props = Record<string, any>> = {
  page: any
  frames: any[]
  error?: any
  props: Props
  target: HTMLElement | null
}

type Hydrate<Props = Record<string, any>> = (input: HydrateInput<Props>) => void

/**
 * Mount function
 */

type MountInput = {
  components: Record<string, any>
  page: string
  frames: string[]
  target: HTMLElement
  error?: string
  createView: Hydrate
  hot?: Hot
}

export function mount(input: MountInput): void {
  const props = getProps(document.getElementById("bud_props"))
  input.createView({
    page: input.components[input.page],
    frames: input.frames.map((frame) => input.components[frame]),
    error: input.error ? input.components[input.error] : undefined,
    target: input.target,
    props: props,
  })
  if (input.hot) {
    input.hot.listen(() => {
      input.createView({
        page: input.components[input.page],
        frames: input.frames.map((frame) => input.components[frame]),
        error: input.error ? input.components[input.error] : undefined,
        target: input.target,
        props: props,
      })
    })
  }
}

function getProps(node: HTMLElement | null) {
  if (!node || !node.textContent) {
    return {}
  }
  try {
    return JSON.parse(node.textContent)
  } catch (err) {
    return {}
  }
}
