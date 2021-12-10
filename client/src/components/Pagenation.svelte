<script type="ts">
  import { createEventDispatcher } from "svelte";

  export let nowPage: number;
  export let end = false;

  const dispatch = createEventDispatcher();

  export function setPage(active: boolean, page: number) {
    if (active) {
      dispatch("page", {
        page: page,
      });
    }
  }
</script>

<div class="container">
  <button class={nowPage===1?"disabled":"enabled"} uk-icon="icon:chevron-left" on:click={()=>setPage(nowPage!==1, nowPage-1)} />
  <button class={nowPage===1?"disabled":"enabled"} on:click={()=>setPage(nowPage!==1, nowPage-1)}>{nowPage===1?"":nowPage-1}</button>
  <button class="active">{nowPage}</button>
  <button class={end?"disabled":"enabled"} on:click={()=>setPage(!end, nowPage+1)}>{end?"":nowPage+1}</button>
  <button class={end?"disabled":"enabled"} uk-icon="icon:chevron-right" on:click={()=>setPage(!end, nowPage+1)} />
</div>

<style>
  .container {
    display: inline-flex;
    justify-content: center;
    align-items: center;
  }
  button {
    width: 35px;
    height: 45px;
    border-radius: 10px;
    color: #212121;
    border: none;
  }
  button.enabled:hover {
    background-color: rgba(17, 31, 77, 0.12);
  }
  button.active {
    background-color: #111F4D;
    color: #f2f4f7;
  }
  button.disabled {
    color: #757575;
    background-color: transparent;
    cursor: default;
  }
</style>
