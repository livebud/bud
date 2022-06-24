type Input = {
  route: string
  view: any // TODO: type this
  props: Record<string, any>
  context: Record<string, any>
}

type Response = {
  status: number
  headers: Record<string, string>
  body: string
}

export function renderHTML(input: Input): Response {
  // Handle the missing view
  if (!input.view) {
    return {
      status: 404,
      headers: {},
      body: fallback(new Error('Missing page "' + input.route + '"')),
    }
  }
  return input.view({ props: input.props, context: input.context })
}

function fallback(err: Error) {
  return `fallback error: ${err.message}`
}
