<script>
  import { format as timeago } from "timeago.js"
  export let story = {}

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
  <div>
    <a class="title" href={story.url || `/${story.id}`}>{story.title}</a>
    {#if story.url}
      <a class="url" href={story.url}>({formatURL(story.url)})</a>
    {/if}
  </div>
  <div class="meta">
    {story.points} points by {story.author} • {timeago(story.created_at)} •
    <a href={`/${story.id}`}>{formatComments(story.num_comments)}</a>
  </div>
</div>

<style>
  .story {
    padding: 10px;
    font-size: 14px;
  }
  a {
    text-decoration: none;
    color: inherit;
  }
  a[href]:hover {
    text-decoration: underline;
  }
  .title {
    font-weight: 500;
  }
  .url,
  .meta {
    color: gray;
  }
</style>
