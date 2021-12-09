<script type="ts">
  import UIkit from 'uikit'
  import apis from '../lib/apis/api';
  import { Resource, ResourceType } from '../lib/apis/generated/api';
  import { toast } from '@zerodevx/svelte-toast';
  import ImageCard from '../components/ImageCard.svelte';
  import ModalImage from '../components/ModalImage.svelte';

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

  let modalResourceID: number;
</script>

<div class="container">
  {#each imageResources as imageResource, i}
    <div class="item">
      <button class="image-btn" uk-toggle="target: #resource-modal" type="button" on:click={()=>modalResourceID = i}>
        <ImageCard resource={imageResource} />
      </button>
    </div>
  {/each}
  <div id="resource-modal" class="uk-flex-top" uk-modal>
    {#if modalResourceID === 0 || modalResourceID}
      <div class="uk-modal-dialog uk-margin-auto-vertical dialog">
        <button class="uk-modal-close-outside" type="button" uk-close></button>
        <ModalImage resource={imageResources[modalResourceID]} />
      </div>
    {/if}
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-wrap: wrap;
    width: 100%;
  }
  .item {
    flex: 0 0 25%;
    max-width: 25%;
    padding: 2%;
  }
  .image-btn {
    border: 0;
    padding: 0;
    width: 100%;
    height: 100%;
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

