<script type="ts">
  import apis from '../lib/apis/api';
  import { GroupInfo, Resource, ResourceType } from '../lib/apis/generated/api';
  import { toast } from '@zerodevx/svelte-toast';
  import ImageCard from '../components/ImageCard.svelte';
  import GroupCard from '../components/GroupCard.svelte';
  import ModalImage from '../components/ModalImage.svelte';
  import SubTitleWithMore from '../components/SubTitleWithMore.svelte';
  import { goto } from '$app/navigation';

  let imageResources: Resource[] = [];
  apis.getResources([ResourceType.Image], undefined, undefined, 4, 0).then(r => {
    imageResources = r.data;
  }).catch(err => {
    console.log(err);
    toast.push("画像ファイル一覧の取得に失敗しました", {
      theme: {
        background: '#e43a19',
        color: '#212121',
      },
    });
  });

  let groups: GroupInfo[] = [];
  apis.getGroups(undefined, undefined, 4, 0).then(r => {
    groups = r.data;
  }).catch(err => {
    console.log(err);
    toast.push("グループ一覧の取得に失敗しました", {
      theme: {
        background: '#e43a19',
        color: '#212121',
      },
    });
  });

  let modalResourceID: number;
</script>

<div class="container">
  <div class="group">
    <SubTitleWithMore title="Latest Images" link="/files?type=image" />
    <div class="resources">
      {#each imageResources as imageResource, i}
        <div class="item">
          <button class="btn" uk-toggle="target: #resource-modal" type="button" on:click={()=>modalResourceID = i}>
            <ImageCard resource={imageResource} />
          </button>
        </div>
      {/each}
    </div>

    <div id="resource-modal" class="uk-flex-top" uk-modal>
      {#if modalResourceID === 0 || modalResourceID}
        <div class="uk-modal-dialog uk-margin-auto-vertical dialog">
          <button class="uk-modal-close-outside" type="button" uk-close></button>
          <ModalImage resource={imageResources[modalResourceID]} />
        </div>
      {/if}
    </div>
  </div>

  <div class="group">
    <SubTitleWithMore title="Latest Groups" link="/groups" />
    <div class="groups">
      {#each groups as group}
        <div class="item">
          <button class="btn" type="button" on:click={()=>goto(`/groups/${group.id}`)}>
            <GroupCard group={group} />
          </button>
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    width: 100%;
  }
  .group {
    display: flex;
    flex-direction: column;
    width: 100%;
  }
  .resources, .groups {
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
    width: fit-content;
    height: 100%;
    justify-content: center;
    align-items: center;
  }
</style>

