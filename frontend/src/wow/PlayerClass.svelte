<script lang="ts">
  import { Classes } from "./api";

  import WowIcon from "./WowIcon.svelte";

  export let playerClass: number;

  $: name = Classes.then((cs) => cs.get(playerClass).name);
</script>

<div {...$$restProps}>
  <WowIcon {playerClass} class="class-icon" height={24} />
  <span>
    {#await name}
      ...
    {:then n}
      {n}
    {/await}
  </span>
</div>
