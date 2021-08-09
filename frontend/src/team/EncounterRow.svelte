<script lang="ts">
  import { Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import Button, { Icon } from "@smui/button/styled";

  import { Encounters } from "./encounters_api";
  import RosterSpecSelect from "./RosterSpecSelect.svelte";
  import { Members } from "./members_api";

  export let memberId: number;

  $: member = Members.member(memberId);
  $: encounters = member.encounters;

  let selectAll: boolean | null = false;

  $: rms = $encounters.then(
    (vs) => new Map(vs.map((m) => [m.encounter_id, m]))
  );

  $: Promise.all([$Encounters, rms]).then(([encs, rms]) => {
    const count = rms.size;
    if (count == 0) {
      selectAll = false;
    } else if (count == encs.length) {
      selectAll = true;
    } else {
      selectAll = null;
    }
  });

  function toggleAll() {
    if (selectAll == null || selectAll == true) {
      encounters.remove();
    } else {
      $member.then((m) =>
        encounters.set({
          encounter_id: undefined,
          member_id: memberId,
          spec: m.config.primarySpec,
        })
      );
    }
  }
</script>

{#await $Encounters}
  <Cell>
    <CircularProgress indeterminate />
  </Cell>
{:then encounters}
  <Cell>
    <Button class="member-encounter-all" on:click={toggleAll}>
      <Icon class="material-icons">
        {#if selectAll == null}
          indeterminate_check_box
        {:else if selectAll == true}
          check_box
        {:else}
          check_box_outline_blank
        {/if}
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
    button.mdc-button {
      min-width: 24px;
    }
    .material-icons {
      font-size: 24px;
    }
  }
</style>
