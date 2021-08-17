<script lang="ts">
  import ProgressCircular from "smelte/src/components/ProgressCircular";
  import Button from "smelte/src/components/Button";
  import Icon from "smelte/src/components/Icon";

  import EncounterNameField from "./EncounterNameField.svelte";
  import { TeamStore } from "./team_store";
  import { LoadingState } from "../store_helpers";

  $: encounters = $TeamStore.Encounters;

  const thClass =
    "encounter-header capitalize duration-100 text-gray-600 text-xs hover:text-black dark-hover:text-white p-3 font-normal text-right";

  function newEncounter() {
    $TeamStore.newEncounter({
      id: undefined,
      name: "",
    });
  }
</script>

{#if encounters.state == LoadingState.Loading}
  <th class={thClass}>
    <ProgressCircular />
  </th>
{:else}
  <th class={thClass} />
  {#each $encounters as enc, i (i)}
    <th class={thClass}>
      <div>
        <EncounterNameField encounterId={enc.id} />
      </div>
    </th>
  {/each}
{/if}
<th class={thClass} style="vertical-align: bottom;">
  <Button on:click={newEncounter}>
    <Icon style="margin-right: 0; color: green">add</Icon>
  </Button>
</th>

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
