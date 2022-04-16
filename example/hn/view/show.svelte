<script>
  import Comment from "./Comment.svelte"
  import Header from "./Header.svelte"
  import Story from "./Story.svelte"
  export let story = {
    children: [],
  }
</script>

<svelte:head>
  <link rel="stylesheet" href="/tailwind/preflight.css" />
</svelte:head>

<div>
  <Header />

  <ul class="mb-12">
    <li class="pl-5 pt-6">
      <Story {story} />
      {#each story.children as comment}
        <Comment {comment} />
      {/each}
    </li>
  </ul>
</div>
