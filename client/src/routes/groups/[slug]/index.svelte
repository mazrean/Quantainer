<script type="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
	export async function load({ page, fetch, session, stuff }) {
    const strPageNum = page.query.get("page");
    const pageNum = strPageNum ? Number(page.query.get("page")) : 1;

    const groupID = page.params.slug;

    const limit = pageNum === 1?19:20;
    const offset = pageNum === 1?0:20*(pageNum - 1)-1;
    const resources = await apis.getResources(undefined, undefined, groupID, limit, offset).then(r => {
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

    const group = await apis.getGroup(groupID).then(r => {
      return r.data;
    }).catch(err => {
      console.log(err);
      toast.push("グループ情報の取得に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    });

    const groups = await apis.getGroups(undefined, undefined, 20, 20*(pageNum - 1)).then(r => {
      return r.data;
    }).catch(err => {
      console.log(err);
      toast.push("グループ一覧の取得に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    });

    return {
      props: {
        group,
        resources,
        groups,
        pageNum,
      }
    };
	}
</script>

<script type="ts">
  import apis from "$lib/apis/api";
  import { GroupDetail, GroupInfo, Resource, ResourceType, WritePermission } from "$lib/apis/generated/api";
  import { toast } from "@zerodevx/svelte-toast";
  import OtherCard from "../../../components/OtherCard.svelte";
  import ImageCard from "../../../components/ImageCard.svelte";
  import ModalImage from "../../../components/ModalImage.svelte";
  import ModalDescription from "../../../components/ModalDescription.svelte";
  import Pagenation from "../../../components/Pagenation.svelte";
  import { goto } from "$app/navigation";
  import SubTitleWithButton from "../../../components/SubTitleWithButton.svelte";
  import ModalAddGroup from "../../../components/ModalAddGroup.svelte";

  export let group: GroupDetail;
  export let resources: Resource[];
  export let groups: GroupInfo[];
  export let pageNum: number;

  groups = groups.filter(g => g.writePermission === WritePermission.Public);

  if (pageNum == 1) {
    resources = [group.mainResource, ...resources]
  }

  let path = `/groups/${group.id}?`;

  let modalResourceID: number;
  let selectedResource: Resource;

  async function addResourceEvent(e: any) {
    await apis.postResourceToGroup(e.detail.groupID, selectedResource.id).then(r => {
      toast.push("グループに追加しました", {
        theme: {
          background: '#4caf50',
          color: '#212121',
        },
      });
      goto(`/groups/${e.detail.groupID}`);
    }).catch(err => {
      console.log(err);
      toast.push("ファイルのグループへの追加に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    });
  }
</script>

<div class="container">
  <SubTitleWithButton title={group.name} buttonLabel="Edit" link={`/groups/${group.id}/edit`} />
  <div class="group-info">
    <p>Admin: {group.administrators.join(', ')}</p>
    <p>{group.description}</p>
  </div>
  <div class="resources" style="grid-template-rows: repeat({(resources.length+3)/4}, 1fr);">
    {#if resources.length > 0}
    {#each resources as resource, i}
      <div class="item">
        {#if resource.resourceType === ResourceType.Image}
          <button class="btn" type="button" on:click={()=>modalResourceID = i}>
            <ImageCard resource={resource} on:group={e=>{selectedResource=e.detail.resource}} />
          </button>
        {:else}
          <button class="btn" uk-toggle="target: #resource-modal" type="button" on:click={()=>modalResourceID = i}>
            <OtherCard resource={resource} on:group={e=>{selectedResource=e.detail.resource}} />
          </button>
        {/if}
      </div>
    {/each}
    {:else}
      No Files
    {/if}
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

  <div id="group-modal" class="uk-flex-top" uk-modal>
    <div class="uk-modal-dialog uk-margin-auto-vertical dialog">
      <button class="uk-modal-close-outside" type="button" uk-close="target: #group-modal"></button>
      <ModalAddGroup groups={groups} on:add={addResourceEvent} />
    </div>
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
  p {
    margin: 0;
    font-size: 20px;
    overflow-wrap: break-word;
  }
  .group-info {
    margin-top: 10px;
    margin-bottom: 10px;
  }
</style>
