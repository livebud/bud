import { parse } from "../url"

/**
 * Hot reload
 */

export default class Hot {
  private subs: Array<() => void> = []
  private sse: EventSource

  constructor(path: string, private readonly components: Record<string, any>) {
    this.sse = new EventSource(path)
    this.sse.addEventListener("message", this.onmessage)
  }

  listen(fn: () => void) {
    this.subs.push(fn)
  }

  private onmessage = (e: MessageEvent) => {
    const payload: { scripts: string[] } = JSON.parse(e.data)
    this.loadScripts(payload.scripts).catch((err) => console.error(err))
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
