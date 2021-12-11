<script type="ts">
  import { user } from "../store/user";
  import type { Resource } from "../lib/apis/generated/api";
  import { createEventDispatcher } from "svelte";

  export let resource: Resource;
  export let addable = true;
  export let modal = true;

  let userName: string = "";
  user.subscribe(user => {
    if (user === null) {
      return;
    }

    userName = user.name;
  });

  const dispatch = createEventDispatcher();

  function addGroup() {
    dispatch("group", {
      resource,
    });
  }
</script>

<div class="container uk-cover-container uk-card">
  <div class="sizer"></div>
  {#if modal}
  <div class="uk-cover file-icon" uk-toggle="target: #resource-modal">
    <span uk-icon="icon:file-text;ratio:5" />
  </div>
  {:else}
  <div class="uk-cover file-icon">
    <span uk-icon="icon:file-text;ratio:5" />
  </div>
  {/if}
  <div class="description">
    <img class="icon" src={`https://q.trap.jp/api/v3/public/icon/${resource.creator}`} alt={resource.creator} />
    <p>{resource.name}</p>
    {#if addable && resource.creator === userName}
    <button uk-icon="icon:plus;ratio:1.3" uk-toggle="target: #group-modal" on:click={addGroup} />
    {/if}
  </div>
</div>

<style>
  .container {
    display: flex!important;
    align-items: flex-end!important;
    border-radius: 10px;
    width: 100%;
    height: 100%;
  }
  .sizer {
    width: 100%;
    padding-bottom: 100%;
    flex: 1 0 0px;
  }
  .icon {
    width: 36px;
    height: 36px;
    border-radius: 50%;
  }
  .description {
    z-index: 1;
    display: inline-flex;
    align-items: center;
    padding: 5px;
    background: rgb(242, 244, 247, 0.5);
    width: 100%;
  }
  p {
    text-overflow: ellipsis;
    white-space: nowrap;
    margin: 0;
    overflow: hidden;
    padding-left: 5px;
    font-size: 17px;
    overflow-wrap: anywhere;
  }
  .file-icon {
    background-color: rgba(17, 31, 77, 0.12);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
  }
  button {
    margin-left: auto;
    margin-right: 0;
  }
</style>
