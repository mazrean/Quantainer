<script type="ts">
  import type { Resource } from "../lib/apis/generated/api";
  import ModalDescription from "./ModalDescription.svelte";

  export let resource: Resource;
</script>

<div class="container">
  <img src={`/api/v1/files/${resource.fileID}`} alt={resource.name}>
  <ModalDescription resource={resource} />
</div>

<style>
  .container {
    background-color: transparent;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
  }
  img {
    width: inherit;
    height: inherit;
    object-fit: contain;
  }
</style>
