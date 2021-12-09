<script type="ts">
  import { goto } from "$app/navigation";
  import Button from "./Button.svelte";
  import SubTitle from "./SubTitle.svelte";

  export let title:string;
  export let buttonLabel: string;
  export let link:string;
</script>

<div class="container">
  <SubTitle title={title} />
  <span />
  <div class="btn">
    <Button label={buttonLabel} on:click={()=>goto(link)} />
  </div>
</div>

<style>
  .container {
    display: inline-flex;
    flex-direction: row;
    align-items: flex-end;
    padding: 0 16px;
  }
  span {
    height: 100%;
    width: 20px;
  }
  .btn {
    margin: 5px 0;
  }
</style>
