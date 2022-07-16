<script>
  export let _string = ""
</script>

<h1>{_string}</h1>

<form action="/10" method="post">
  <input type="hidden" name="_method" value="patch" />
  <input type="text" name="title" placeholder="title" />
  <input type="submit" value="Update Post" />
</form>

<style>
  h1 {
    margin: 0;
    background: whitesmoke;
    border-bottom: 2px dashed red;
    padding: 20px;
    color: red;
  }
</style>
