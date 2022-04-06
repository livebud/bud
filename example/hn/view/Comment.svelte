<script>
  export let comment = {}
  let numChildren = (comment.children || []).length
  let show = true

  function toggle() {
    show = !show
  }
</script>

<div class="pl-5 pt-6">
  <div class="ml-4 flex items-center">
    <span class="text-gray-600" />
    <span class="text-sm ml-2 text-gray-600"
      >{comment.author}
      {comment.created_at}
      [<a href={"#"} on:click={toggle}>{show ? "–" : `${numChildren} more`}</a>]
    </span>
  </div>
  {#if show}
    <div class="mt-1 ml-6">
      <div class="max-w-prose leading-relaxed text-gray-900">
        {@html comment.text}
      </div>
    </div>
    {#if comment.children}
      {#each comment.children as comment}
        <svelte:self {comment} />
      {/each}
    {/if}
  {/if}
</div>
