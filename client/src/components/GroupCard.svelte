<script type="ts">
  import type { GroupInfo } from "../lib/apis/generated/api";
  import { ResourceType } from "../lib/apis/generated/api";

  export let group: GroupInfo;
</script>

<div class="container uk-cover-container uk-card">
  <div class="sizer"></div>
  {#if group.mainResource.resourceType === ResourceType.Image }
    <img class="uk-cover thumbnail" src={`/api/v1/files/${group.mainResource.fileID}`} alt={group.name}>
  {:else}
    <div class="uk-cover file-icon">
      <span uk-icon="icon:file-text;ratio:5" />
    </div>
  {/if}

  <div class="description">
    <p>{group.name}</p>
  </div>
</div>

<style>
  .container {
    display: flex!important;
    align-items: flex-end!important;
    border-radius: 10px;
  }
  .sizer {
    width: 100%;
    padding-bottom: 100%;
    flex: 1 0 0px;
  }
  .thumbnail {
    height: 100%;
    width: 100%;
    object-fit: cover;
  }
  .file-icon {
    background-color: rgba(17, 31, 77, 0.12);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
  }
  .description {
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 5px;
    background: rgb(242, 244, 247, 0.5);
    width: 100%;
  }
  p {
    text-overflow: ellipsis;
    margin: 0;
    overflow: hidden;
    padding: 0 5px;
    height: 36px;
    font-size: 17px;
  }
</style>
