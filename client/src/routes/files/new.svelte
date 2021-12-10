<script type="ts">
  import { goto } from "$app/navigation";
  import type { ModelFile } from "$lib/apis/generated/api";
  import { ResourceType } from "$lib/apis/generated/api";
  import { toast } from "@zerodevx/svelte-toast";
  import Button from "../../components/Button.svelte";
  import FileInput from "../../components/FileInput.svelte";
  import ImageContainer from "../../components/ImageContainer.svelte";
  import SubTitle from "../../components/SubTitle.svelte";
  import apis from '../../lib/apis/api';

  let fileInput: any;

  let file: ModelFile = null;
  const resourceTypes: ResourceType[] = [ResourceType.Image, ResourceType.Other];

  let name: string = "";
  let resourceType: ResourceType = ResourceType.Image;
  let comment: string = "";

  async function selectFileEvent(e: any) {
    const fileData = e.detail.file;

    await uploadFile(fileData);
  }

  async function changeFileEvent(e: any) {
    const fileData = e.target.files[0];

    await uploadFile(fileData);
  }

  async function uploadFile(fileData: any) {
    const res = (await apis.postFile(fileData).catch(err => {
      console.log(err);
      toast.push("ファイルのアップロードに失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    }))

    if (res) {
      file = res.data;
    }
  }

  async function createResource(e: any) {
    const res = (await apis.postResource(
      file.id,
      {
        name,
        resourceType,
        comment,
      },
    ).catch(err => {
      console.log(err);
      toast.push("リソースの作成に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    }))

    if (res) {
      toast.push("リソースの作成に成功しました", {
        theme: {
          background: '#2e7d32',
          color: '#212121',
        },
      });
      goto(`/files/${res.data.id}/edit`);
    }
  }
</script>

<div class="container">
  <SubTitle title="New File" />
  <div class="form-container">
    <div class="file-input">
      {#if file}
        <ImageContainer file={file} />
        <Button label="Change" on:click={()=>fileInput.click()} />
        <input style="display:none" type="file" on:change={changeFileEvent} bind:this={fileInput} >
      {:else}
        <FileInput on:fileSelected={selectFileEvent} />
      {/if}
    </div>
    <div class="input">
      <input class="uk-input" placeholder="タイトル" type="text" bind:value={name}>
      <select class="uk-select select" style="height: 56px;" bind:value={resourceType}>
        {#each resourceTypes as resourceType}
          <option value={resourceType}>{resourceType}</option>
        {/each}
      </select>
      <textarea class="uk-textarea" placeholder="コメント" cols="30" rows="10" bind:value={comment} />
      <Button label="Create" on:click={createResource} />
    </div>
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
  }
  .form-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    column-gap: 10px;
    width: 100%;
    height: 100%;
    margin: 10px 0;
  }
  .input {
    width: 100%;
  }
  input {
    height: 56px;
    width: 100%;
    background-color: rgba(17, 31, 77, 0.12);
    border-radius: 10px;
    color: #757575;
    margin-bottom: 10px;
    font-size: 20px;
  }
  input:focus {
    background-color: rgba(17, 31, 77, 0.12);
    border-color: #111f4d;
  }
  textarea {
    width: 100%;
    background-color: rgba(17, 31, 77, 0.12);
    border-radius: 10px;
    color: #757575;
    margin-bottom: 10px;
    resize: none;
    font-size: 20px;
  }
  textarea:focus {
    background-color: rgba(17, 31, 77, 0.12);
    border-color: #111f4d;
  }
  .select {
    height: 56px;
    width: 100%;
    background-color: rgba(17, 31, 77, 0.12);
    border-radius: 10px;
    color: #757575;
    margin-bottom: 10px;
    font-size: 20px;
  }
  .select:focus {
    background-color: rgba(17, 31, 77, 0.12);
    border-color: #111f4d;
  }
</style>
