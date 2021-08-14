<script lang="ts">
  import LinearProgress from "@smui/linear-progress/styled";
  import WowIcon from "../wow/WowIcon.svelte";
  import { TeamStore } from "./team_store";

  export let memberId: number;

  $: row = $TeamStore.row(memberId);
  $: member = row.member;
</script>

<div class="member">
  {#if !$member}
    <LinearProgress indeterminate />
  {:else}
    <WowIcon
      playerClass={$member.classId}
      spec={$member.config.primarySpec}
      class="icon"
      height={24}
    />
    <span>{$member.name}</span>
  {/if}
</div>

<style lang="scss">
  .member {
    max-width: 10em;

    :global(.icon) {
      padding-right: 4px;
    }
  }
</style>
