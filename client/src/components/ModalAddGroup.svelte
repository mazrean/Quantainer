<script type="ts">
  import type { GroupInfo } from "$lib/apis/generated/api";
  import { createEventDispatcher } from "svelte";
  import Button from "./Button.svelte";
import SubTitle from "./SubTitle.svelte";

  export let groups: GroupInfo[];

  let groupID: string;

  const dispatch = createEventDispatcher();
</script>

<div class="container">
  <SubTitle title="New Group File" />
  <select class="uk-select select" style="height: 56px;" bind:value={groupID}>
    {#each groups as group}
      <option value={group.id}>{group.name}</option>
    {/each}
  </select>
  <Button label="Add" on:click={()=>dispatch("add", {groupID})} />
</div>

<style>
  .select {
    height: 56px;
    width: 100%;
    background-color: rgba(17, 31, 77, 0.12);
    border-radius: 10px;
    color: #757575;
    margin-bottom: 10px;
    font-size: 20px;
  }
  .select:focus {
    background-color: rgba(17, 31, 77, 0.12);
    border-color: #111f4d;
  }
  .container {
    background-color: #f2f4f7;
    height: fit-content;
    padding: 15px;
    margin: auto;
    border-radius: 10px;
    min-width: 50%;
    max-width: 500px;
  }
</style>
