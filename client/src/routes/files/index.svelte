<script type="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
	export async function load({ page, fetch, session, stuff }) {
    const types = page.query.getAll("type")
    for (const type of types) {
      if (type !== ResourceType.Image && type !== ResourceType.Other) {
        toast.push("ファイルの種類が誤っています", {
          theme: {
            background: '#e43a19',
            color: '#212121',
          },
        });
      }
    }

    const strPageNum = page.query.get("page");
    const pageNum = strPageNum ? Number(page.query.get("page")) : 1;

    const resources = await apis.getResources(types, undefined, undefined, 20, 20*(pageNum - 1)).then(r => {
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

    return {
      props: {
        resources,
        pageNum,
        types,
      }
    };
	}
</script>

<script type="ts">
  import apis from "$lib/apis/api";
  import { Resource, ResourceType } from "$lib/apis/generated/api";
  import { toast } from "@zerodevx/svelte-toast";
  import OtherCard from "../../components/OtherCard.svelte";
  import ImageCard from "../../components/ImageCard.svelte";
  import ModalImage from "../../components/ModalImage.svelte";
  import SubTitleWithButton from "../../components/SubTitleWithButton.svelte";
  import ModalDescription from "../../components/ModalDescription.svelte";
  import Pagenation from "../../components/Pagenation.svelte";
  import { goto } from "$app/navigation";

  export let resources: Resource[];
  export let pageNum: number;
  export let types: ResourceType[];

  let path = "/files?";
  if (types.length > 0) {
    path += "types=" + types.map(t => t.toString()).join("&type=");
  }

  let modalResourceID: number;
</script>

<div class="container">
  <SubTitleWithButton title="Files" buttonLabel="New File" link="/files/new" />
  <div class="resources" style="grid-template-rows: repeat({(resources.length+3)/4}, 1fr);">
    {#each resources as resource, i}
      <div class="item">
        {#if resource.resourceType === ResourceType.Image}
          <button class="btn" uk-toggle="target: #resource-modal" type="button" on:click={()=>modalResourceID = i}>
            <ImageCard resource={resource} />
          </button>
        {:else}
          <button class="btn" uk-toggle="target: #resource-modal" type="button" on:click={()=>modalResourceID = i}>
            <OtherCard resource={resource} />
          </button>
        {/if}
      </div>
    {/each}
  </div>

  <div class="pagenation">
    <Pagenation nowPage={pageNum} end={resources.length < 20} on:page={e=>goto(`${path}&page=${e.detail.page}`)} />
  </div>

  <div id="resource-modal" class="uk-flex-top" uk-modal>
    {#if modalResourceID === 0 || modalResourceID}
      <div class="uk-modal-dialog uk-margin-auto-vertical dialog">
        <button class="uk-modal-close-outside" type="button" uk-close></button>
        {#if resources[modalResourceID].resourceType === ResourceType.Image}
          <ModalImage resource={resources[modalResourceID]} />
        {:else}
          <ModalDescription resource={resources[modalResourceID]} />
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    width: 100%;
  }
  .resources {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
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
  #resource-modal {
    justify-content: center;
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
