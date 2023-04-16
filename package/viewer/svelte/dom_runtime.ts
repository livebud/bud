import { noop, set_current_component } from 'svelte/internal'
import { SvelteComponentTyped } from 'svelte'
import type Hot from 'livebud/runtime/hot'

type Key = string

type Page = {
  key: Key
  frames: Key[]
  error?: Key
  components: {
    [key: string]: typeof SvelteComponentTyped
  }
  hot?: Hot
}

export function mount(page: Page) {
  const target = document.getElementById('bud_target')
  if (!target) return
  const data = getState(document.getElementById('bud_state'))
  data.props = data.props || {}

  // Workaround to prevent the inline components from throwing.
  // Issue: https://github.com/sveltejs/svelte/issues/6584
  set_current_component({ $$: {} });

  // organize views from innermost page to outermost frame
  let keys: Key[] = [page.key, ...page.frames.reverse()]
  let outermost: Key | undefined = keys.pop()
  if (!outermost) return

  // Compose the components
  let component: SvelteComponentTyped | undefined
  for (let key of keys) {
    let props = data.props[key] || {}

    // If we already have an inner component, slot it in
    if (component) {
      props['$$scope'] = {}
      props['$$slots'] = {
        default: [
          slot(component)
        ],
      }
    }

    // @ts-expect-error ts(2345) target is required in the typedef, but not
    // required in practice.
    component = new page.components[key]({
      $$inline: true,
      props: props
    })
  }

  // Hydrate the outermost (most likely a frame)
  let props = data.props[outermost] || {}
  // If we already have an inner component, slot it in
  if (component) {
    props['$$scope'] = {}
    props['$$slots'] = {
      default: [
        slot(component)
      ],
    }
  }
  new page.components[outermost]({
    target: target,
    hydrate: true,
    props: props,
  })

  // Start listening for hot reloads
  if (page.hot) {
    page.hot.listen(() => {
      console.log("oh hai.")
    })
  }
}

// Try getting the state from an HTML element or return an empty object
function getState(node: HTMLElement | null) {
  if (!node || !node.textContent) {
    return {}
  }
  try {
    return JSON.parse(node.textContent)
  } catch (err) {
    return {}
  }
}


// Internal implementation to support hydrating with slots
// Based on: https://github.com/sveltejs/svelte/pull/4296
function slot(element: SvelteComponentTyped) {
  // Load the fragment
  let frag = (element && element.$$ && element.$$.fragment) || {
    c: noop, // Create
    l: noop, // Listen(?)
    m: noop, // Mount
    d: noop, // Detach
  }
  return function () {
    return frag
  }
}