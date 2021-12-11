<script type="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
	export async function load({ page, fetch, session, stuff }) {
    const types = page.query.getAll("type")
    for (const type of types) {
      if (type !== GroupType.ArtBook && type !== GroupType.Other) {
        toast.push("グループの種類が誤っています", {
          theme: {
            background: '#e43a19',
            color: '#212121',
          },
        });
      }
    }

    const strPageNum = page.query.get("page");
    const pageNum = strPageNum ? Number(page.query.get("page")) : 1;

    const groups = await apis.getGroups(types, undefined, 20, 20*(pageNum - 1)).then(r => {
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
        groups,
        pageNum,
        types,
      }
    };
	}
</script>

<script type="ts">
  import apis from "$lib/apis/api";
  import { GroupInfo, GroupType } from "$lib/apis/generated/api";
  import { toast } from "@zerodevx/svelte-toast";
  import SubTitleWithButton from "../../components/SubTitleWithButton.svelte";
  import Pagenation from "../../components/Pagenation.svelte";
  import { goto } from "$app/navigation";
  import GroupCard from "../../components/GroupCard.svelte";

  export let groups: GroupInfo[];
  export let pageNum: number;
  export let types: GroupType[];

  let path = "/groups?";
  if (types.length > 0) {
    path += "types=" + types.map(t => t.toString()).join("&type=");
  }
</script>

<div class="container">
  <SubTitleWithButton title="Groups" buttonLabel="New Group" link="/groups/new" />
  <div class="resources" style="grid-template-rows: repeat({(groups.length+3)/4}, 1fr);">
    {#if groups.length > 0}
    {#each groups as group}
      <div class="item">
        <button class="btn" type="button" on:click={()=>goto(`/groups/${group.id}`)}>
          <GroupCard group={group} />
        </button>
      </div>
    {/each}
    {:else}
      No Groups
    {/if}
  </div>

  <div class="pagenation">
    <Pagenation nowPage={pageNum} end={groups.length < 20} on:page={e=>goto(`${path}&page=${e.detail.page}`)} />
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
  .pagenation {
    display: flex;
    width: 100%;
    justify-content: center;
    align-items: center;
  }
</style>
