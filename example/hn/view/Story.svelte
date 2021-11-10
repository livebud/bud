<script>
  import { format as timeago } from "timeago.js"
  export let i = 0
  export let story = {
    id: 0,
    title: "",
    url: "",
    points: 0,
    author: "",
    created_at: new Date(),
    num_comments: 0,
  }

  function formatURL(url) {
    if (!url) return ""
    const parsed = new URL(url)
    return parsed.host
  }

  function formatComments(num_comments) {
    switch (num_comments) {
      case 1:
        return "1 comment"
      default:
        return `${num_comments || 0} comments`
    }
  }
</script>

<li class="p-5 bg-white flex items-center">
  {#if i > 0}<div class="text-sm text-gray-500">{i}.</div>{/if}
  <div class="ml-5 text-orange-600">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-6 w-6"
      viewBox="0 0 20 20"
      fill="currentColor"
    >
      <path
        fill-rule="evenodd"
        d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z"
        clip-rule="evenodd"
      />
    </svg>
  </div>
  <div class="pl-5">
    <div class="text-lg">
      <a href={story.url}>{story.title}</a>
      {#if story.url}
        <span class="text-sm text-gray-500">
          ({formatURL(story.url)})
        </span>
      {/if}
    </div>
    <div class="text-sm text-gray-500">
      {story.points} points by {story.author} • {timeago(story.created_at)} •
      <a href={`/${story.id}`}>{formatComments(story.num_comments)}</a>
    </div>
  </div>
</li>
