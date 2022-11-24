<script>
  export let title = "My Wonderful Blog"
  export let theme = "light"
</script>

<html lang="en">
  <head>
    <title>{title}</title>
  </head>
  <body data-theme={theme}>
    <h1>{title}</h1>
    <main>
      <slot />
    </main>
  </body>
</html>
