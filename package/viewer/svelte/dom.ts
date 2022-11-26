/**
 * dom.ts is the runtime for rendering browser components
 */

import type { SvelteComponentTyped } from "svelte/internal"
import { noop, set_current_component } from "svelte/internal"

type View = {
  props?: Record<string, any>
  context?: Record<string, any>
}

type Page = View & {
  frames: View[]
}

const emptyPage: Page = {
  props: {},
  context: {},
  frames: []
}

type Component = typeof SvelteComponentTyped

type Input = {
  Page: Component
  frames: Component[]
}

export class Hydrator {
  constructor(private readonly input: Input) { }

  hydrate(targetElement?: HTMLElement, dataElement?: HTMLElement) {
    const { Page, frames } = this.input
    const data = this.data(dataElement)

    // Workaround to prevent the inline components from throwing.
    // Issue: https://github.com/sveltejs/svelte/issues/6584
    set_current_component({ $$: {} });

    // Get the first frame
    const TopFrame = frames.shift()
    if (!TopFrame) {
      // Render and mount the page
      return new Page({
        props: data.props,
        context: new Map(Object.entries(data.context || {})),
        target: targetElement || document.body
      })
    }

    // Initialize but don't mount the page
    let component = new Page({
      $$inline: true,
      props: data.props,
      context: new Map(Object.entries(data.context || {})),
      // @ts-ignore Svelte allows null targets with $$inline true, but the types
      // don't permit it.
      target: null
    })


    // Compose the middle frames
    for (let i = frames.length - 1; i > 0; i--) {
      const Frame = frames[i]
      const frameData = data.frames[i]
      component = new Frame({
        $$inline: true,
        context: new Map(Object.entries(frameData.context || {})),
        props: {
          ...frameData.props,
          $$scope: {},
          $$slots: {
            default: [this.slot(component)],
          },
        },
        // @ts-ignore Svelte allows null targets with $$inline true, but the types
        // don't permit it.
        target: null
      })
    }

    // Render the top frame
    const frameData = data.frames.shift() || {}
    return new TopFrame({
      hydrate: true,
      target: targetElement || document.body,
      context: new Map(Object.entries(frameData.context || {})),
      props: {
        ...frameData.props || {},
        $$scope: {},
        $$slots: {
          default: [this.slot(component)]
        },
      },
    })
  }

  data(dataElement?: HTMLElement): Page {
    if (!dataElement) {
      return emptyPage
    }
    try {
      return JSON.parse(dataElement.textContent || '')
    } catch (e) {
      return emptyPage
    }
  }

  // Internal implementation to support hydrating with slots
  // Based on: https://github.com/sveltejs/svelte/pull/4296
  private slot(component: SvelteComponentTyped) {
    // Load the fragment
    let frag = (component && component.$$ && component.$$.fragment) || {
      c: noop, // Create
      l: noop, // Listen(?)
      m: noop, // Mount
      d: noop, // Detach
    }
    return function () {
      return frag
    }
  }
}

