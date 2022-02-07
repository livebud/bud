<script>
  import Header from "./Header.svelte"
  import Story from "./Story.svelte"
  export let props = {
    stories: [],
  }
</script>

<!-- <svelte:head>
  <link rel="stylesheet" href="/tailwind/preflight.css" />
</svelte:head> -->

<main>
  <Header />

  <div>
    {#each props.stories as story, i}
      <Story i={i + 1} {story} />
    {/each}
  </div>
</main>
