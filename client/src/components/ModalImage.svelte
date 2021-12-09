<script type="ts">
  import type { Resource } from "../lib/apis/generated/api";

  export let resource: Resource;
</script>

<div class="container">
  <img src={`/api/v1/files/${resource.fileID}`} alt={resource.name}>
  <div class="description">
    <h1>{resource.name}</h1>
    <p>{resource.comment}</p>
  </div>
</div>

<style>
  .container {
    background-color: transparent;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
  }
  .description {
    background-color: #f2f4f7;
    height: 100%;
    padding: 15px;
    border-radius: 10px;
  }
  h1 {
    font-size: 28px;
    margin: 0;
  }
</style>
