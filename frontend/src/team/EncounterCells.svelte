<script lang="ts">
  import ProgressCircular from "smelte/src/components/ProgressCircular";
  import Icon from "smelte/src/components/Icon";
  import RosterSpecSelect from "./RosterSpecSelect.svelte";

  import { TeamStore } from "./team_store";

  export let memberId: number;

  $: storeRow = $TeamStore.row(memberId);
  $: member = storeRow.member;
  $: encounters = $TeamStore.Encounters;
  $: memberEncounters = storeRow.encounters;

  let selectAll: boolean | null = false;
  const tdClass = "member-encounter relative px-3 font-normal text-right";

  $: {
    const rmCount = $memberEncounters ? $memberEncounters.length : 0;
    const encCount = $encounters ? $encounters.length : 0;
    if (rmCount == 0) {
      selectAll = false;
    } else if (rmCount == encCount) {
      selectAll = true;
    } else {
      selectAll = null;
    }
  }

  const icons = new Map([
    [true, "check_box"],
    [false, "check_box_outline_blank"],
    [null, "indeterminate_check_box"],
  ]);

  function toggleAll() {
    if (selectAll == null || selectAll == true) {
      storeRow.encounterAPI.remove();
    } else {
      storeRow.encounterAPI.put([
        ...$encounters.map((e) => ({
          encounter_id: e.id,
          member_id: memberId,
          spec: $member.config.primarySpec,
        })),
      ]);
    }
  }
</script>

{#await $encounters}
  <td class={tdClass}>
    <ProgressCircular />
  </td>
{:then encounters}
  <td class={`${tdClass} member-encounter-all`}>
    <Icon on:click={toggleAll} class="cursor-pointer">
      {icons.get(selectAll)}
    </Icon>
  </td>

  {#each encounters as enc, i (i)}
    <td class={tdClass}>
      <RosterSpecSelect {memberId} encounterId={enc.id} />
    </td>
  {/each}
{/await}
<td />

<style lang="scss" global>
</style>
