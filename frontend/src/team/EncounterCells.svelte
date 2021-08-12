<script lang="ts">
  import { Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import Button, { Icon } from "@smui/button/styled";

  import RosterSpecSelect from "./RosterSpecSelect.svelte";
  import { TeamStore } from "./team_store";

  export let memberId: number;

  $: storeRow = $TeamStore.row(memberId);
  $: member = storeRow.memberInfo;
  $: encounters = $TeamStore.Encounters;
  $: memberEncounters = storeRow.memberEncounters;

  let selectAll: boolean | null = false;

  $: rms = $memberEncounters.then(
    (members) => new Map(members.map((m) => [m.encounter_id, m]))
  );

  $: Promise.all([$encounters, rms]).then(([encs, rms]) => {
    const count = rms.size;
    if (count == 0) {
      selectAll = false;
    } else if (count == encs.length) {
      selectAll = true;
    } else {
      selectAll = null;
    }
  });

  const icons = new Map([
    [true, "check_box"],
    [false, "check_box_outline_blank"],
    [null, "indeterminate_check_box"],
  ]);

  function toggleAll() {
    if (selectAll == null || selectAll == true) {
      memberEncounters.remove();
    } else {
      Promise.all([$member, $encounters]).then(([m, es]) =>
        storeRow.memberEncounters.set([
          ...es.map((e) => ({
            encounter_id: e.id,
            member_id: memberId,
            spec: m.config.primarySpec,
          })),
        ])
      );
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
    <Cell style="overflow: visible">
      <RosterSpecSelect {memberId} encounterId={enc.id} />
    </Cell>
  {/each}
{/await}
<Cell />

<style lang="scss" global>
  .member-encounter-all {
    padding-left: 2px;
    padding-right: 2px;

    button.mdc-button {
      min-width: 24px;
    }
    .material-icons {
      font-size: 24px;
    }
  }
</style>
