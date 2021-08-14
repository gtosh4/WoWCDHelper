<script lang="ts">
  import { Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import Button, { Icon } from "@smui/button/styled";

  import EncounterNameField from "./EncounterNameField.svelte";
  import { TeamStore } from "./team_store";
  import { LoadingState } from "../store_helpers";

  $: encounters = $TeamStore.Encounters;

  function newEncounter() {
    $TeamStore.newEncounter({
      id: undefined,
      name: "",
    });
  }
</script>

{#if encounters.state == LoadingState.Loading}
  <Cell class="encounter-header">
    <CircularProgress indeterminate />
  </Cell>
{:else}
  <Cell class="encounter-header" />
  {#each $encounters as enc, i (i)}
    <Cell class="encounter-header">
      <div>
        <EncounterNameField encounterId={enc.id} />
      </div>
    </Cell>
  {/each}
{/if}
<Cell style="vertical-align: bottom;">
  <Button style="min-width: 32px" on:click={newEncounter}>
    <Icon class="material-icons" style="margin-right: 0; color: green">
      add
    </Icon>
  </Button>
</Cell>

<style lang="scss" global>
  th.encounter-header {
    height: 7em; // 90 degree triangle with 8em hypotenuse with some padding
    white-space: nowrap;
    overflow: visible;
    vertical-align: bottom;

    & > div {
      transform: rotate(-45deg);
      width: 1em;
    }

    .encounter-name {
      width: 8em;
      height: auto;
    }
  }
</style>
