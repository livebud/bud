import Link from "./Link"
import { ComponentChildren } from "preact"

type ButtonProps = {
  href: string
  children: ComponentChildren
}

export default function Button(props: ButtonProps) {
  return (
    <Link href={props.href}>
      <button>{props.children}</button>
    </Link>
  )
}
