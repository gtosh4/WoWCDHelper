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

  const itemClass =
    "focus:bg-gray-50 dark-focus:bg-gray-700 hover:bg-gray-transDark relative overflow-hidden duration-100 p-2 cursor-pointer text-gray-700 dark:text-gray-100 flex items-center z-10";
</script>

<ul class="assign-palette py-2 rounded divide-y">
  {#each $members || [] as member (member.member_id)}
    <li class={itemClass}>
      <Assign memberId={member.member_id} {encounterId} />
    </li>
  {/each}
</ul>

<style lang="scss" global>
  .assign-palette {
    height: 100%;
  }
</style>
