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

<div class="story">
  <div class="text-lg">
    <a class="title" href={story.url}>{i}. {story.title}</a>
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

<style>
  .story {
    padding: 10px;
  }
  .title {
    font-weight: 500;
  }
</style>
