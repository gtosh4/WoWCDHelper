<script lang="ts">
  import { Classes } from "./api";

  import WowIcon from "./WowIcon.svelte";

  export let playerClass: number;

  $: name = $Classes.get(playerClass).name;
</script>

<div>
  <WowIcon {playerClass} class="class-icon" height={24} />
  <span>{name}</span>
</div>
