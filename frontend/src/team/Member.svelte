<script lang="ts">
  import LinearProgress from "@smui/linear-progress/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { TeamMember } from "./members_api";

  export let memberId: number;

  $: member = TeamMember(memberId);
</script>

<div class="member">
  {#await $member}
    <LinearProgress indeterminate />
  {:then m}
    <WowIcon
      playerClass={m.classId}
      spec={m.config.primarySpec}
      class="icon"
      height={24}
    />
    <span>{m.name}</span>
  {/await}
</div>

<style lang="scss">
  .member {
    max-width: 10em;

    :global(.icon) {
      padding-right: 4px;
    }
  }
</style>
