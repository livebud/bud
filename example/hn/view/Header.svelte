<div class="list">
  <img src="https://news.ycombinator.com/y18.gif" alt="Hacker News" />
  <a href="/"><h1>Hacker News</h1></a>
  <div class="links">
    <a href="/news">New</a>
    <a href="/ask">Ask</a>
    <a href="/show">Show</a>
  </div>
</div>

<style>
  .list {
    display: flex;
    align-items: center;
    padding: 10px;
  }

  a {
    text-decoration: none;
    color: inherit;
  }

  img {
    display: block;
    height: 20px;
    width: 20px;
  }

  h1 {
    margin-left: 10px;
    font-weight: 600;
    font-size: 1rem;
  }

  h1:hover {
    text-decoration: none;
  }

  .links {
    margin-left: 10px;
    display: flex;
    align-items: center;
  }

  .links > a {
    margin-left: 10px;
    font-size: 75%;
  }

  .links > a:hover {
    text-decoration: underline;
  }
</style>
