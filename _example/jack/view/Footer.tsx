import Link from "./Link"

type FooterProps = {
  success?: boolean
}

export default function Footer(props: FooterProps) {
  return (
    <footer>
      <div className="left-scene" />
      <div className="scene">
        <div className="grass" />
        <Link href="/" prefetch>
          <div className="jack" />
        </Link>
      </div>
      <div className="right-scene">
        <div className="right-scene-inner" />
      </div>
    </footer>
  )
}
