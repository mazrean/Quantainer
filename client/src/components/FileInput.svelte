<script type="ts">
  import { createEventDispatcher } from "svelte";

  let fileInput: any;
  const dispatch = createEventDispatcher();

  function fileSelected(e: any) {
    const file = e.target.files[0];
    if (file) {
      dispatch("fileSelected", {
        file,
      });
    }
  }
</script>

<div class="container" on:click={()=>fileInput.click()}>
  <span uk-icon="icon:cloud-upload;ratio:5" />
  <p>Upload File</p>
</div>
<input style="display:none" type="file" on:change={fileSelected} bind:this={fileInput} >

<style>
  .container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 400px;
    background-color: rgba(17, 31, 77, 0.12);
    border-radius: 10px;
    flex-direction: column;
  }
  p {
    margin: 0;
    font-size: 20px;
  }
  span {
    width: 56px;
    height: 56px;
  }
</style>
