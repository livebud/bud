<script>
  import Header from "./Header.svelte"
  import Story from "./Story.svelte"
  export let stories = []
</script>

<main>
  <Header />

  <div>
    {#each stories as story}
      <Story {story} />
    {/each}
  </div>
</main>
