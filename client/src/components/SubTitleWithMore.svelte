<script type="ts">
  import { goto } from "$app/navigation";
  import SubTitle from "./SubTitle.svelte";

  export let title:string;
  export let link:string;
</script>

<div class="container">
  <SubTitle title={title} />
  <div class="more" on:click={()=>goto(link)}>
    <p>more</p>
    <span uk-icon="icon:arrow-right" />
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-direction: row;
    align-items: flex-end;
    padding: 0;
  }
  .more {
    margin-left: 15px;
    display: inline-flex;
    flex-direction: row;
    align-items: center;
    color: #757575;
    font-size:30px;
    cursor: pointer;
  }
  p {
    margin: 0;
  }
  span {
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
