<script>
  export let id = 0
  export let title = "Untitled"
  export let body = "No body"
</script>

<article {id}>
  <h1>{title}</h1>
  <p>{body}</p>
</article>

<style>
  h1 {
    color: blue;
  }
  p {
    color: red;
  }
</style>
