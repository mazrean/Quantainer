<script type="ts">
  import Header from '../components/Header.svelte';
</script>

<div>
  <Header userName="mazrean"></Header>
  <div class="main">
    <slot></slot>
  </div>
</div>
