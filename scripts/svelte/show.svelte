<script>
  export let planet = "Mars"
  export let id = 0
  export let count = 0
</script>

<h1>Hello {planet} from {id}</h1>

<button on:click={() => count++}>Total: {count}</button>
