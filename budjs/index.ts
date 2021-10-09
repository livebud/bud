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
  error?: string
  hydrate: Hydrate
  hot?: Hot
}

export function mount(input: MountInput): void {
  const props = getProps(document.getElementById("duo_props"))
  input.hydrate({
    page: input.components[input.page],
    frames: input.frames.map((frame) => input.components[frame]),
    error: input.components[input.error],
    target: document.getElementById("duo_target"),
    props: props,
  })
  if (input.hot) {
    input.hot.listen(() => {
      input.hydrate({
        page: input.components[input.page],
        frames: input.frames.map((frame) => input.components[frame]),
        error: input.components[input.error],
        target: document.getElementById("duo_target"),
        props: props,
      })
    })
  }
}

function getProps(node) {
  if (!node) {
    return {}
  }
  try {
    return JSON.parse(node.textContent)
  } catch (err) {
    return {}
  }
}
