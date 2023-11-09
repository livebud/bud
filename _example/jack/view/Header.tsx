import SlackButton from "./SlackButton"
import Button from "./Button"
import Link from "./Link"

type HeaderProps = {
  success?: boolean
}

export default function Header(props: HeaderProps = {}, context = {}) {
  return (
    <header class="Header">
      <div className="left-scene" />
      <div className="scene">
        <nav className="buttons">
          {/* <Button href={`mailto:${props.email}`}>Contact</Button> */}
          <Button href="/faq">FAQ</Button>
          <SlackButton success={props.success} />
        </nav>
        <div className="top-tagline">
          <span>Meet</span>
          <span className="jack-text" />
          <span className="comma">,</span>
        </div>
        <div className="tagline">A Slack Bot for your Standups</div>
        <div className="grass" />
        <Link href="/" prefetch>
          <div className="jack" />
        </Link>
        <div className="cloud-1" />
        <div className="cloud-2" />
      </div>
      <div className="right-scene">
        <div className="right-scene-inner" />
      </div>
      <nav className="mobile-buttons">
        {/* <a className="Button" href={`mailto:${props.email}`}>
          Contact
        </a> */}
        <Button href="/faq">FAQ</Button>
      </nav>
    </header>
  )
}
