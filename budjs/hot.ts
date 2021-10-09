import { parse } from "./url"

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

  // // load scripts in order
  // private loadScripts(...scripts: string[]) {
  //   const self = this
  //   const src = scripts.shift()
  //   if (!src) return
  //   const script = document.createElement("script")
  //   script.type = "module"
  //   script.src = src + "?ts=" + Math.random()
  //   function next() {
  //     self.loadScripts(...scripts)
  //   }
  //   function error() {
  //     // TODO: better error handling
  //     throw new Error("unable to load script")
  //   }
  //   script.addEventListener("error", function () {
  //     script.removeEventListener("error", error)
  //     error()
  //   })
  //   script.addEventListener("load", function () {
  //     script.removeEventListener("load", next)
  //     next()
  //   })
  //   // Add or replace existing script
  //   const existing = document.querySelector(`script[src="${src}"]`)
  //   if (existing && existing.parentNode) {
  //     existing.parentNode.replaceChild(script, existing)
  //   } else {
  //     document.head.appendChild(script)
  //   }
  // }

  close() {
    this.sse.removeEventListener("message", this.onmessage)
    this.sse.close()
  }
}
