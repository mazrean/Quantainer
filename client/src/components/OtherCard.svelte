<script type="ts">
import type { Resource } from "../lib/apis/generated/api";

  export let resource: Resource;
</script>

<div class="container uk-cover-container uk-card">
  <div class="sizer"></div>
  <div class="uk-cover file-icon">
    <span uk-icon="icon:file-text;ratio:5" />
  </div>
  <div class="description">
    <img class="icon" src={`https://q.trap.jp/api/v3/public/icon/${resource.creator}`} alt={resource.creator} />
    <p>{resource.name}</p>
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
  }
  .file-icon {
    background-color: rgba(17, 31, 77, 0.12);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
  }
</style>
