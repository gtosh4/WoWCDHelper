<script lang="ts">
  import ProgressCircular from "smelte/src/components/ProgressCircular";
  import Button from "smelte/src/components/Button";
  import Icon from "smelte/src/components/Icon";

  import EncounterNameField from "./EncounterNameField.svelte";
  import { TeamStore } from "./team_store";
  import { LoadingState } from "../store_helpers";

  export let height = 24;

  $: encounters = $TeamStore.Encounters;

  const thClass = `encounter-header capitalize duration-100 text-gray-600 text-xs hover:text-black dark-hover:text-white px-3 font-normal text-right h-${height}`;

  function newEncounter() {
    $TeamStore.newEncounter({
      id: undefined,
      name: "",
    });
  }
</script>

{#if encounters.state == LoadingState.Loading}
  <th class={thClass}>
    <ProgressCircular size={26} />
  </th>
{:else}
  <th class={thClass} />
  {#each $encounters as enc, i (i)}
    <th class={thClass}>
      <div class="relative">
        <div class="transform -rotate-45 absolute -bottom-8">
          <EncounterNameField
            encounterId={enc.id}
            replace={{
              "w-full": `w-${height}`,
              "\\S*bg-\\S*": "bg-transparent",
              "p[xt]-\\d+": "",
            }}
          />
        </div>
      </div>
    </th>
  {/each}
{/if}
<th class={thClass} style="vertical-align: bottom;">
  <Button on:click={newEncounter} outlined small text>
    <Icon>add</Icon>
  </Button>
</th>

<style lang="scss" global>
</style>
