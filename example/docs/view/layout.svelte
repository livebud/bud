<script lang="ts">
  export let title = "My Wonderful Blog"
  export let theme = "light"
  // import Header from "./Header.svelte"
</script>

<html lang="en">
  <head>
    <title>{title}</title>
    <slot name="head" />
  </head>
  <body data-theme={theme}>
    <h1>{title}</h1>
    <main>
      <slot />
    </main>
  </body>
</html>

<style>
  h1 {
    color: green;
  }
</style>
