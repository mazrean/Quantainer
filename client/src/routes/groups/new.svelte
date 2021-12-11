<script type="ts">
  import { goto } from "$app/navigation";
  import type { Resource } from "$lib/apis/generated/api";
  import { GroupType, ResourceType, ReadPermission, WritePermission } from "$lib/apis/generated/api";
  import { toast } from "@zerodevx/svelte-toast";
  import Pagenation from "../../components/Pagenation.svelte";
  import Button from "../../components/Button.svelte";
  import SubTitle from "../../components/SubTitle.svelte";
  import apis from '../../lib/apis/api';
  import ImageCard from "../../components/ImageCard.svelte";
  import OtherCard from "../../components/OtherCard.svelte";

  let page = 1;

  let name: string = "";
  let type: GroupType = GroupType.ArtBook;
  let description: string = "";
  let readPermission: ReadPermission = ReadPermission.Public;
  let writePermission: WritePermission = WritePermission.Public;

  class ResourceItem {
    resource: Resource;
    status: string;
  }

  let mainResourceItem: ResourceItem = undefined;
  let selectedResourceItems: ResourceItem[] = [];

  let pageNum = 1;
  let resourceItems: ResourceItem[] = [];
  async function updateResource(pageNum: number) {
    resourceItems = [];
    for (let resource of await getResources(pageNum)) {
      let status = "none";
      if (mainResourceItem && resource.id === mainResourceItem.resource.id) {
        status = "main";
      } else if (selectedResourceItems.find((item) => item.resource.id === resource.id)) {
        status = "normal";
      }

      resourceItems = [{ resource, status }, ...resourceItems];
    }
  };
  updateResource(pageNum);

  function next(i: number) {
    switch (resourceItems[i].status) {
      case "none":
        resourceItems[i] = {
          resource: resourceItems[i].resource,
          status: "normal",
        };

        selectedResourceItems.push(resourceItems[i]);
        break;
      case "normal":
        resourceItems[i] = {
          resource: resourceItems[i].resource,
          status: "main",
        };

        if (mainResourceItem) {
          selectedResourceItems.push(mainResourceItem);
        }
        mainResourceItem = resourceItems[i];
        for (let j = 0; j < selectedResourceItems.length; j++) {
          if (selectedResourceItems[j].resource.id === mainResourceItem.resource.id) {
            selectedResourceItems.splice(j, 1);
          }
        }

        for (let j = 0; j < resourceItems.length; j++) {
          if (j !== i && resourceItems[j].status === "main") {
            resourceItems[j] = {
              resource: resourceItems[j].resource,
              status: "normal",
            };
            break;
          }
        }
        break;
      case "main":
        mainResourceItem = undefined;

        resourceItems[i] = {
          resource: resourceItems[i].resource,
          status: "none",
        };
        break;
    }
    resourceItems = [...resourceItems];
  }

  async function createGroup(e: any) {
    let resourceIDs: string[] = selectedResourceItems.map((item) => item.resource.id);

    const res = (await apis.postGroup({
      name,
      type,
      description,
      readPermission,
      writePermission,
      mainResourceID: mainResourceItem ? mainResourceItem.resource.id : undefined,
      resourceIDs,
    }).catch(err => {
      console.log(err);
      toast.push("グループの作成に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    }))

    if (res) {
      toast.push("グループの作成に成功しました", {
        theme: {
          background: '#2e7d32',
          color: '#212121',
        },
      });
      goto(`/groups/${res.data.id}`);
    }
  }

  async function getResources(pageNum: number) {
    const newResources = await apis.getResources(undefined, undefined, undefined, 20, 20*(pageNum - 1)).then(r => {
      return r.data;
    }).catch(err => {
      console.log(err);
      toast.push("ファイル一覧の取得に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    });

    if (newResources) {
      return newResources;
    }
    return [];
  }
</script>

<div class="container">
  <SubTitle title="New Group" />
  {#if page === 1}
    <div class="form-container">
      <input class="uk-input" placeholder="名前" type="text" bind:value={name}>
      <select class="uk-select select" style="height: 56px;" bind:value={type}>
        {#each [GroupType.ArtBook, GroupType.Other] as groupType}
          <option value={groupType}>{groupType}</option>
        {/each}
      </select>
      <textarea class="uk-textarea" placeholder="説明" cols="30" rows="10" bind:value={description} />
      <div class="permission-container">
        <div class="select-container">
          <p>ファイル閲覧権限:</p>
          <select class="uk-select select" style="height: 56px;" bind:value={readPermission}>
            {#each [ReadPermission.Public, ReadPermission.Private] as readPermission}
              <option value={readPermission}>{readPermission}</option>
            {/each}
          </select>
        </div>
        <div class="select-container">
          <p>ファイル追加権限:</p>
          <select class="uk-select select" style="height: 56px;" bind:value={writePermission}>
            {#each [WritePermission.Public, WritePermission.Private] as writePermission}
              <option value={writePermission}>{writePermission}</option>
            {/each}
          </select>
        </div>
      </div>
      <Button label="Next" on:click={()=>page=2} />
    </div>
  {:else}
    <div class="resources" style="grid-template-rows: repeat({(resourceItems.length+4)/5}, 1fr);">
      {#each resourceItems as resourceItem, i}
        <div class="item">
          <button class="btn" type="button" on:click={()=>next(i)}>
            <div class="card-wrapper {resourceItems[i].status}">
              {#if resourceItem.resource.resourceType === ResourceType.Image}
                <ImageCard resource={resourceItem.resource} />
              {:else}
                <OtherCard resource={resourceItem.resource} />
              {/if}
            </div>
          </button>
        </div>
      {/each}
    </div>

    <div class="pagenation">
      <Pagenation nowPage={pageNum} end={resourceItems.length < 20} on:page={e=>{pageNum=e.detail.page;updateResource(pageNum)}} />
    </div>

    <div>
      <Button label="Prev" on:click={()=>page =1} />
      <Button label="Create" on:click={createGroup} />
    </div>
  {/if}
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
  }
  .form-container {
    width: 100%;
  }
  .permission-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    column-gap: 10px;
    width: 100%;
    height: 100%;
    margin: 10px 0;
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
    width: 100%;
    height: 56px;
    background-color: rgba(17, 31, 77, 0.12);
    border-radius: 10px;
    color: #757575;
    font-size: 20px;
  }
  .select:focus {
    background-color: rgba(17, 31, 77, 0.12);
    border-color: #111f4d;
  }
  .card-wrapper {
    border-radius: 10px;
  }
  .normal {
    background-color: #111F4D;
    padding: 3px;
  }
  .main {
    background-color: #e43a19;
    padding: 3px;
  }
  .select-container {
    height: 56px;
    display: flex;
    flex-direction: row;
    align-items: center;
    margin-bottom: 10px;
  }
  p {
    white-space: nowrap;
    margin: 0;
    width: fit-content;
  }
  .resources {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    column-gap: 16px;
  }
  .item {
    width: 100%;
    height: 100%;
  }
  .btn {
    border: 0;
    padding: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
  }
  .dialog {
    background-color: transparent;
    margin: 0!important;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .pagenation {
    display: flex;
    width: 100%;
    justify-content: center;
    align-items: center;
  }
</style>
