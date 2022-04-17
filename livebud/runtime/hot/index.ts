import { parse } from "../../url"

/**
 * Hot reload
 */

export default class Hot {
  private subs: Array<() => void> = []
  private sse: EventSource
  private queue = new Queue()

  constructor(path: string, private readonly components: Record<string, any>) {
    this.sse = new EventSource(path)
    this.sse.addEventListener("message", this.onmessage)
  }

  listen(fn: () => void) {
    this.subs.push(fn)
  }

  private onmessage = (e: MessageEvent) => {
    // TODO: define a protocol
    const payload: { scripts: string[]; reload: boolean } = JSON.parse(e.data)
    if (payload.reload) {
      location.reload()
      return
    }
    this.queue.enqueue(() => {
      this.loadScripts(payload.scripts).catch((err) => console.error(err))
    })
  }

  private async loadScripts(scripts: string[]) {
    for (let scriptPath of scripts) {
      const imported = await import(scriptPath)
      const url = parse(scriptPath)
      this.components[url.pathname] = imported.default
      for (let sub of this.subs) {
        sub()
      }
    }
  }

  close() {
    this.sse.removeEventListener("message", this.onmessage)
    this.sse.close()
  }
}

/**
 * Simple queue to ensure updates only happen one at a time, in order.
 */
class Queue {
  private queue: Array<() => void> = []
  private idle: boolean = true
  enqueue(fn: () => void) {
    this.queue.push(fn)
    this.process()
  }
  private process() {
    if (!this.idle) return
    this.idle = false
    let fn: (() => void) | undefined
    while ((fn = this.queue.shift())) {
      fn()
    }
    this.idle = true
  }
}
