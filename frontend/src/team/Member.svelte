<script lang="ts">
  import ProgressLinear from "smelte/src/components/ProgressLinear";
  import WowIcon from "../wow/WowIcon.svelte";
  import { TeamStore } from "./team_store";

  export let memberId: number;

  $: row = $TeamStore.row(memberId);
  $: member = row.member;
</script>

<div class="member">
  {#if !$member}
    <ProgressLinear />
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

<style lang="scss" global>
  .member {
    max-width: 10em;

    .icon {
      // padding-right: 4px;
    }
  }
</style>
