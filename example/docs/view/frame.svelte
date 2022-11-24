<script>
  export let categories = []
</script>

<main>
  <aside>
    {#each categories as category}
      <p>{category}</p>
    {/each}
  </aside>
  <article><slot /></article>
</main>
