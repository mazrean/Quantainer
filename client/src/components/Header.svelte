<script type="ts">
  export let userName: string;
</script>

<div class="container">
  <h1>Quantainer</h1>
  <div class="right">
    <span uk-icon="icon:plus;ratio:1.5" class="icon"></span>
    <span class="span"></span>
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
  .span {
    width: 10px;
    height: 36px;
  }
  img {
    width: 36px;
    height: 36px;
    border-radius: 50%;
  }
  .icon {
    color: #f2f4f7;
  }
</style>
