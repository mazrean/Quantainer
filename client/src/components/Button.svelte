<script type="ts">
  export let label: string;
</script>

<button on:click>{label}</button>

<style>
  button {
    padding: 5px;
    background-color: #111f4d;
    color: #f2f4f7;
    border: none;
    border-radius: 10px;
    font-size: 17px;
    cursor: pointer;
  }
</style>
