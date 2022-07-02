<script>
  import cookie from "../module/cookie"
  // TODO: fix SSR flickering. cookie.get(...) uses document.cookie under the
  // hood. `document.cookie` should be accessible by V8 on the server-side.
  // We'll need to take care exposing APIs though. We shouldn't try to polyfill
  // the DOM on the server.
  const theme = cookie.get("theme")
  function submit(e) {
    this.form.submit()
  }
</script>

<div class={theme || "light"}>
  <h1>Change Theme</h1>
  <form action="/" method="post">
    <select name="theme" on:change={submit}>
      <option value="light" selected={theme === "light"}>Light</option>
      <option value="dark" selected={theme === "dark"}>Dark</option>
    </select>
  </form>
</div>

<style>
  .light {
    --var-color: black;
    --var-bgcolor: whiteSmoke;
  }
  .dark {
    --var-color: whiteSmoke;
    --var-bgcolor: black;
  }
  h1 {
    color: var(--var-color);
    background: var(--var-bgcolor);
  }
</style>
