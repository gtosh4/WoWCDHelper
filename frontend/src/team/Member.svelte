<script lang="ts">
  import WowIcon from "../wow/WowIcon.svelte";
  import { TeamMember } from "./api";

  export let memberId: number;

  $: member = TeamMember(memberId);

  $: console.log("member", $member);
</script>

<div class="member">
  {#if $member}
    <WowIcon
      playerClass={$member.classId}
      spec={$member.config.primary_spec}
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
