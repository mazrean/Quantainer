<script type="ts">
  import { goto } from "$app/navigation";

  const items = [
    {
      label: "File",
      icon: "file-text",
      link: "/files",
      hover: false,
    },
    {
      label: "Group",
      icon: "folder",
      link: "/groups",
      hover: false,
    },
  ];
</script>

<div>
  <ul class="uk-list">
    {#each items as item}
      <!-- svelte-ignore a11y-mouse-events-have-key-events -->
      <li class:hover={item.hover} on:click={() => goto(item.link)} on:mouseover={() => item.hover=true} on:mouseout={() => item.hover=false}>
        <span uk-icon="icon:{item.icon};ratio:1.5" />
        <span class="span" />
        <p>{item.label}</p>
      </li>
    {/each}
  </ul>
</div>

<style>
  div {
    width: 100%;
  }
  ul {
    background-color: #f2f4f7;
    width: 100%;
    margin: 0;
  }
  li {
    padding: 10px;
    height: 56px;
    display: flex;
    align-items: center;
    font-size: 20px;
    color: #212121;
    border-radius: 10px;
  }
  .span {
    height: 100%;
    width: 10px;
  }
  p {
    margin: 0;
  }
  .hover {
    background-color: rgba(17, 31, 77, 0.12);
  }
</style>
