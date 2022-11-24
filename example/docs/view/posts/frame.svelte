<script>
  export let likes = 0
</script>

<div>
  <slot />
  <hr />
  {likes} people liked this.
</div>
