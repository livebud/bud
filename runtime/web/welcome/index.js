import JSConfetti from "js-confetti"

const jsConfetti = new JSConfetti()

jsConfetti.addConfetti()

const sse = new EventSource("http://0.0.0.0:35729")
sse.addEventListener("message", () => {
  location.reload()
})
