<script>
  export let post = {}
</script>

<h1>Edit Post</h1>

<form method="post" action={`/${post.id || 0}`}>
	<input type="hidden" name="_method" value="patch">
	<!-- Add input fields here -->
	<input type="submit" value="Update Post" />
</form>

<br />

<a href={`/`}>Back</a>
<span> | </span>
<a href={`/${post.id || 0}`}>Show Post</a>
