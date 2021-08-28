<script lang="ts">
  import Assign from "./Assign.svelte";

  import { TeamStore } from "../team/team_store";
  import { derived } from "svelte/store";

  export let encounterId: number;

  $: column = $TeamStore.column(encounterId);
  $: column.memberAPI.get();

  $: memberList = column.members;
  $: members = derived(memberList, (ms) =>
    ms ? ms.filter((m) => m != undefined) : []
  );
</script>

<ul class="assign-palette py-2 rounded divide-y">
  {#each $members || [] as member (member.member_id)}
    <li class="">
      <Assign class="w-full" memberId={member.member_id} {encounterId} />
    </li>
  {/each}
</ul>

<style lang="scss" global>
  .assign-palette {
    height: 100%;
  }
</style>
