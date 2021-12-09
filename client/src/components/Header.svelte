<script type="ts">
  export let userName: string;
</script>

<div class="container">
  <a href="/"><h1>Quantainer</h1></a>
  <div class="right">
    <!--<span uk-icon="icon:plus;ratio:1.3" class="icon"></span>
    <span class="span"></span>-->
    <img src="https://q.trap.jp/api/v3/public/icon/{userName}" alt="{userName}"/>
  </div>
</div>

<style>
  .container {
    display: inline-flex;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    background-color: #111f4d;
  }
  a {
    text-decoration: none;
  }
  h1 {
    font-size: 28px;
    margin: 15px 16px;
    color: #f2f4f7;
  }
  .right {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-direction: row;
    margin-bottom: 1rem;
    margin: 15px 16px;
  }
  img {
    width: 36px;
    height: 36px;
    border-radius: 50%;
  }
  /*.span {
    width: 10px;
    height: 36px;
  }
  .icon {
    color: #f2f4f7;
  }*/
</style>
