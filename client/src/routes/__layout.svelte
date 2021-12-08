<script type="ts">
  import Header from '../components/Header.svelte';
  import Sidebar from '../components/Sidebar.svelte';
</script>

<div class="container">
  <Header userName="mazrean"></Header>
  <div class="main">
    <div class="sidebar">
      <Sidebar></Sidebar>
    </div>
    <slot></slot>
  </div>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  .main {
    flex: 1;
    display: flex;
    flex-direction: row;
    background-color: #f2f4f7;
    height: 100%;
  }
  .sidebar {
    width: 240px!important;
    min-height: 100%;
    border-right: 1px #e5e5e5 solid;
  }
</style>
