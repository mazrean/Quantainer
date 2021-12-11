<script type="ts">
  import { user } from "../store/user";
  import type { Resource } from "../lib/apis/generated/api";
  import { createEventDispatcher } from "svelte";

  export let resource: Resource;

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
  <img class="uk-cover thumbnail" uk-toggle="target: #resource-modal" src={`/api/v1/files/${resource.fileID}`} alt={resource.name}>
  <div class="description">
    <img class="icon" src={`https://q.trap.jp/api/v3/public/icon/${resource.creator}`} alt={resource.creator} />
    <p>{resource.name}</p>
    {#if resource.creator === userName}
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
  .thumbnail {
    height: 100%;
    width: 100%;
    object-fit: cover;
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
  }
  button {
    margin-left: auto;
    margin-right: 0;
  }
</style>
