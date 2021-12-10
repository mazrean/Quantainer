<script type="ts">
  import type { ModelFile } from "../lib/apis/generated/api";
  export let fileID: string;

</script>

<div class="img-container">
  <img src={`/api/v1/files/${fileID}`} alt={fileID}>
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
