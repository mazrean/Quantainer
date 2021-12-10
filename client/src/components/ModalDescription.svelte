<script type="ts">
  import type { Resource } from "../lib/apis/generated/api";

  export let resource: Resource;
</script>

<div class="description">
  <h1>{resource.name}</h1>
  <p>{resource.comment}</p>
</div>

<style>
  .description {
    background-color: #f2f4f7;
    height: fit-content;
    width: inherit;
    padding: 15px;
    margin: 10px;
    border-radius: 10px;
    min-width: 200px;
  }
  h1 {
    font-size: 28px;
    margin: 0;
  }
  p {
    overflow-wrap: break-word;
  }
</style>
