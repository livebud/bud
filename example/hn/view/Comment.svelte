<script>
  export let comment = {}
  let numChildren = (comment.children || []).length
  let show = true

  function toggle() {
    show = !show
  }
</script>

<ul>
  <li class="pl-5 pt-6">
    <div class="ml-4 flex items-center">
      <span class="text-gray-600" />
      <span class="text-sm ml-2 text-gray-600"
        >{comment.author}
        {comment.created_at}
        [<span on:click={toggle}>{show ? "–" : `${numChildren} more`}</span>]
      </span>
    </div>
    {#if show}
      <div class="ml-6 mt-1">
        <p class="max-w-prose leading-relaxed text-gray-900">
          {@html comment.text}
        </p>
        <!-- <a class="mt-2 text-sm underline">reply</a> -->
      </div>

      {#if comment.children}
        <ul class="ml-6">
          {#each comment.children as comment}
            <svelte:self {comment} />
          {/each}
        </ul>
      {/if}
    {/if}
  </li>
</ul>
