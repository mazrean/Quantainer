<script type="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
	export async function load({ page, fetch, session, stuff }) {
    const resourceID: string = page.params.slug;
    const resource = await apis.getResource(resourceID).then(r => {
      return r.data;
    }).catch(err => {
      console.log(err);
      toast.push("ファイル情報の取得に失敗しました", {
        theme: {
          background: '#e43a19',
          color: '#212121',
        },
      });
    });

    return {
      props: {
        resource,
      }
    };
	}
</script>

<script type="ts">
  import type { Resource } from "$lib/apis/generated/api";
  import { toast } from "@zerodevx/svelte-toast";
  import ImageContainer from "../../../components/ImageContainer.svelte";
  import SubTitle from "../../../components/SubTitle.svelte";
  import apis from '../../../lib/apis/api';

  export let resource: Resource;
</script>

<div class="container">
  <SubTitle title={resource.name} />
  <div class="form-container">
    <div class="file-input">
      <ImageContainer fileID={resource.fileID} />
    </div>
    <div class="info">
      <div class="field">
        <p class="created-at">Created At: {new Date(resource.createdAt).toLocaleString()}</p>
      </div>
      <div class="field">
        <p>Type: {resource.resourceType}</p>
      </div>
      <div class="field">
        <p class="comment">{resource.comment}</p>
      </div>
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
  .info {
    width: 100%;
  }
  .field {
    display: flex;
    align-items: center;
    height: 56px;
  }
  p {
    margin: 0;
    font-size: 20px;
    overflow-wrap: break-word;
  }
</style>
