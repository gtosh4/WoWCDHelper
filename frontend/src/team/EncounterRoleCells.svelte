<script lang="ts">
  import { Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";

  import { Encounters } from "./encounters_api";

  export let roleType: string;

  $: rosters = $Encounters.then((encs) =>
    encs.map((enc) => Encounters.encounter(enc.id).roster)
  );
</script>

{#await rosters}
  <Cell>
    <CircularProgress indeterminate />
  </Cell>
{:then rosters}
  <Cell />
  {#each rosters as roster, i (i)}
    <Cell />
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
