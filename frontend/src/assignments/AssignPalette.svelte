<script lang="ts">
  import Assign from "./Assign.svelte";

  import { TeamStore } from "../team/team_store";
  import { derived } from "svelte/store";

  export let encounterId: number;

  $: column = $TeamStore.column(encounterId);
  $: memberList = column.members;
  $: members = derived(memberList, (ms) =>
    ms ? ms.filter((m) => m != undefined) : []
  );
</script>

<ul class="assign-palette mdc-list">
  {#each $members || [] as member, i (member.member_id)}
    <li class="assign-member">
      <span class="mdc-list-item__content">
        <Assign memberId={member.member_id} {encounterId} />
      </span>
    </li>
    {#if i < $members.length - 1}
      <li class="mdc-list-divider member-divider" role="separator" />
    {/if}
  {/each}
</ul>

<style lang="scss" global>
  .assign-palette {
    height: 100%;

    .member-divider {
      padding-top: 2px;
      padding-bottom: 2px;
    }
  }
</style>
