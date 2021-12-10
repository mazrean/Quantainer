<script type="ts">
  import type { Resource } from "../lib/apis/generated/api";
  import ModalDescription from "./ModalDescription.svelte";

  export let resource: Resource;
</script>

<div class="container">
  <div class="img-container">
    <img src={`/api/v1/files/${resource.fileID}`} alt={resource.name}>
  </div>
  <div class="description-container">
    <ModalDescription resource={resource} />
  </div>
</div>

<style>
  .container {
    background-color: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
  }
  .img-container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-content: center;
    max-width: 700px;
  }
  .description-container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-content: center;
    margin: 0 5px;
  }
  img {
    object-fit: contain;
    max-width: 100%;
    max-height: 100%;
    inset: 0;
  }
</style>
