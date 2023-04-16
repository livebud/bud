<script>
  // import uid from "uid"
  export let planet = "Mars"
  let count = 0
</script>

<h1>Hello {planet}!</h1>

<button on:click={() => count++}>Total: {count}</button>
