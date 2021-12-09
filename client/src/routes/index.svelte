<script type="ts">
  import apis from '../lib/apis/api';
  import { Resource, ResourceType } from '../lib/apis/generated/api';
  import { toast } from '@zerodevx/svelte-toast';
  import ImageCard from '../components/ImageCard.svelte';

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
</script>

<div class="container">
  {#each imageResources as imageResource}
    <div class="item">
      <ImageCard resource={imageResource} />
    </div>
  {/each}
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
</style>

