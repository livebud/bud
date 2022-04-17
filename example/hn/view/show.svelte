<script>
  import Comment from "./Comment.svelte"
  import Header from "./Header.svelte"
  import Story from "./Story.svelte"
  export let story = {
    children: [],
  }
</script>

<div>
  <Header />
  <Story {story} />
  {#each story.children as comment}
    <Comment {comment} />
  {/each}
</div>
