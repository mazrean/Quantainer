<script type="ts">
  import type { ModelFile } from "../lib/apis/generated/api";
  export let file: ModelFile;

</script>

<div class="img-container">
  <img src={`/api/v1/files/${file.id}`} alt={file.id}>
</div>

<style>
  .img-container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-content: center;
    max-width: 700px;
  }
  img {
    object-fit: contain;
    max-width: 100%;
    max-height: 100%;
    inset: 0;
  }
</style>
