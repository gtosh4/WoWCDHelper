<script lang="ts">
  import { Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import Button, { Icon } from "@smui/button/styled";

  import RosterSpecSelect from "./RosterSpecSelect.svelte";
  import { TeamStore } from "./team_store";

  export let memberId: number;

  $: storeRow = $TeamStore.row(memberId);
  $: member = storeRow.member;
  $: encounters = $TeamStore.Encounters;
  $: memberEncounters = storeRow.encounters;

  let selectAll: boolean | null = false;

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
  <Cell>
    <CircularProgress indeterminate />
  </Cell>
{:then encounters}
  <Cell class="member-encounter-all">
    <Button on:click={toggleAll}>
      <Icon class="material-icons">
        {icons.get(selectAll)}
      </Icon>
    </Button>
  </Cell>

  {#each encounters as enc, i (i)}
    <Cell class="member-encounter">
      <RosterSpecSelect {memberId} encounterId={enc.id} />
    </Cell>
  {/each}
{/await}
<Cell />

<style lang="scss" global>
  .member-encounter-all {
    padding-left: 2px;
    padding-right: 16px;

    button.mdc-button {
      min-width: 24px;

      .mdc-button__icon {
        margin: 0;
      }
    }
    .material-icons {
      font-size: 24px;
    }
  }

  .member-encounter {
    overflow-y: visible;
    padding: 0;
  }
</style>
