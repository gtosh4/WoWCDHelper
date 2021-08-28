<script lang="ts">
  import EncounterSelect from "./EncounterSelect.svelte";
  import EventTable from "./EventTable.svelte";

  import { PathPart } from "../url";
  import AssignPalette from "./AssignPalette.svelte";

  $: encounterPath = PathPart(2);
  $: encounterId = +$encounterPath;
</script>

<div class="w-full">
  <div class="grid w-full justify-items-center">
    <EncounterSelect />
  </div>
  <div class="flex">
    <div class="flex-grow pr-4">
      <EventTable {encounterId} />
    </div>
    <div class="flex-shrink">
      <AssignPalette {encounterId} />
    </div>
  </div>
</div>

<style lang="scss" global>
</style>
