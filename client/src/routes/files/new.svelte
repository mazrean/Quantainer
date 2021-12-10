<script type="ts">
  import type { ModelFile } from "$lib/apis/generated/api";
  import FileInput from "../../components/FileInput.svelte";
  import ImageContainer from "../../components/ImageContainer.svelte";
  import SubTitle from "../../components/SubTitle.svelte";

  let file: ModelFile = null;
</script>

<div class="container">
  <SubTitle title="New File" />
  <div class="form-container">
    {#if file}
      <ImageContainer file={file} />
    {:else}
      <FileInput />
    {/if}
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .form-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
  }
</style>
