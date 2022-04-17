<script>
  import { format as timeago } from "timeago.js"
  export let comment = {}
  let show = true
  function toggle() {
    show = !show
  }
</script>

<div class="comment">
  <div class="header">
    <a class="fold" href={"#"} on:click={toggle}>{show ? "↓" : `→`}</a>
    {comment.author}
    {timeago(comment.created_at)}
  </div>
  {#if show}
    <div class="body">
      {@html comment.text}
    </div>
    {#if comment.children}
      {#each comment.children as comment}
        <svelte:self {comment} />
      {/each}
    {/if}
  {/if}
</div>

<style>
  .comment {
    padding: 10px;
  }
  .header {
    color: gray;
    font-size: 75%;
  }
  .fold {
    text-decoration: none;
    color: inherit;
  }
  .body {
    padding-left: 13px;
    font-size: 14px;
  }
  .body :global(a) {
    text-decoration: none;
    color: inherit;
  }
  .body :global(a:hover) {
    text-decoration: underline;
  }
</style>
