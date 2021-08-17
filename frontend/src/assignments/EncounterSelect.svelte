<script lang="ts">
  import Button from "smelte/src/components/Button";
  import Icon from "smelte/src/components/Icon";
  import Select from "smelte/src/components/Select";

  import { PathPart } from "../url";
  import { TeamStore } from "../team/team_store";
  import EncounterNameField from "../team/EncounterNameField.svelte";

  $: encounters = $TeamStore.Encounters;

  $: encounterPath = PathPart(2);

  function select(id: number) {
    $encounterPath = `${id}`;
  }

  $: console.info("encounter select", { ep: $encounterPath });

  $: isSelect = $encounters ? $encounters.length > 1 : false;

  function itemClass(id: number): string {
    let c =
      "focus:bg-gray-50 dark-focus:bg-gray-700 hover:bg-gray-transDark relative overflow-hidden duration-100 p-2 cursor-pointer text-gray-700 dark:text-gray-100 flex items-center z-10";
    if (id == +$encounterPath) {
      c += " bg-gray-200 dark:bg-primary-transLight";
    }
    return c;
  }
</script>

<div class="encounter-select">
  {#if isSelect}
    <Select dense>
      <Button slot="select">
        <Icon class="encounter-select-arrow">arrow_drop_down</Icon>
      </Button>

      <div
        slot="options"
        class="absolute left-0 bg-white rounded shadow w-full z-20 dark:bg-dark-500"
      >
        <ul class="py-2 rounded">
          {#each $encounters as enc}
            <li class={itemClass(enc.id)} on:click={() => select(enc.id)}>
              {enc.name}
            </li>
          {/each}
        </ul>
      </div>
    </Select>
  {/if}
  <EncounterNameField
    encounterId={+$encounterPath}
    input$placeholder="Encounter"
  />
</div>

<style lang="scss" global>
  .encounter-select {
    display: inline-flex;

    // .encounter-select-arrow {
    //   color: white;
    // }

    // button.mdc-button {
    //   min-width: unset;
    //   padding-left: 4px;
    //   padding-right: 4px;
    // }

    // .material-icons {
    //   font-size: 24px;
    // }

    // .mdc-select__anchor {
    //   width: unset;
    // }
  }
</style>
