<script>
  export let likes = 10
  let count = 0
  function onClick() {
    count += 1
  }
</script>

<div>
  <slot />
  <hr />
  <p>
    {likes} people liked this.
    <button on:click={onClick}>Count: {count}</button>
  </p>
</div>

<style>
  p {
    color: brown;
  }
</style>
