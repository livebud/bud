<script>
  export let categories = ["Soccer", "Football"]
</script>

<main>
  <aside>
    {#each categories as category}
      <p>{category}</p>
    {/each}
  </aside>
  <article><slot /></article>
</main>

<style>
  p {
    color: orange;
  }
</style>
