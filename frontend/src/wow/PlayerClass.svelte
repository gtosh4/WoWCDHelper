<script lang="ts">
  import { Classes } from "./api";

  import WowIcon from "./WowIcon.svelte";

  export let playerClass: number;
  let classes = "";
  export { classes as class };

  $: name = Classes.then((cs) => cs.get(playerClass).name);
</script>

<div class={`flex items-start ${classes}`} {...$$restProps}>
  <WowIcon {playerClass} class="pr-1" height={24} />
  <span>
    {#await name}
      ...
    {:then n}
      {n}
    {/await}
  </span>
</div>
