<script lang="ts">
  import { Classes } from "./api";

  import WowIcon from "./WowIcon.svelte";

  export let playerClass: number;

  $: className = $Classes.then((cs) => cs.get(playerClass).name);
</script>

<div>
  <WowIcon {playerClass} class="class-icon" height={24} />
  <span>
    {#await className}
      ...
    {:then name}
      {name}
    {/await}
  </span>
</div>
