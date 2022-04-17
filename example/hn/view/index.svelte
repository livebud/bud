<script>
  import Header from "./Header.svelte"
  import Story from "./Story.svelte"
  export let news = {
    stories: [],
  }
</script>

<main>
  <Header />

  <div>
    {#each news.stories as story, i}
      <Story i={i + 1} {story} />
    {/each}
  </div>
</main>
