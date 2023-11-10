// _document is only rendered on the server side and not on the client side
import Document, { Head, Page, Scripts } from "./Document"

const site = {
  lang: "en",
  title: "Standup Jack",
  description:
    "Standup Jack is a Slack Bot for your Standups. Each weekday at a time of your choosing, Jack will message you a few questions. These questions are designed to help you plan your day and let your teammates know what you're up to.",
  url: "https://standupjack.com",
  twitter: "@mattmueller",
  card: "https://standupjack.com/static/images/card-wide.png",
  favicon: "https://standupjack.com/static/favicon.ico",
}

export default class Base extends Document {
  render() {
    return (
      <html lang={site.lang}>
        <Head>
          <meta name="description" content={site.description} />

          <meta property="og:title" content={site.title} />
          <meta property="og:url" content={site.url} />
          <meta property="og:description" content={site.description} />
          <meta property="og:image:type" content="image/png" />
          <meta property="og:image:width" content="940" />
          <meta property="og:image:height" content="550" />
          <meta property="og:image" content={`${site.card}`} />

          <meta name="twitter:card" content="summary_large_image" />
          <meta name="twitter:site" content={site.twitter} />
          <meta name="twitter:creator" content={site.twitter} />
          <meta name="twitter:image" content={`${site.card}`} />

          <link rel="shortcut icon" href={site.favicon} />
          <link rel="icon" sizes="16x16 32x32" href={site.favicon} />

          <meta name="viewport" content="width=device-width, initial-scale=1" />

          <meta httpEquiv="X-UA-Compatible" content="IE=edge,chrome=1" />
          <meta charSet="utf-8" />

          <link rel="stylesheet" href="/view/layout.css" />
          <link rel="stylesheet" href="/view/index.css" />
          <link rel="stylesheet" href="/view/faq.css" />
        </Head>
        <body>
          <Page />
          <Scripts />
          <script
            async
            src="https://www.googletagmanager.com/gtag/js?id=UA-10351690-15"
          />
          <script src="/js/ga.js" async defer></script>
        </body>
      </html>
    )
  }
}
