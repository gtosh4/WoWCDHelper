<script lang="ts">
  import Card, { Content } from "@smui/card/styled";
  import EncounterSelect from "./EncounterSelect.svelte";
  import EventTable from "./EventTable.svelte";

  import { PathPart } from "../url";
  import AssignPalette from "./AssignPalette.svelte";

  $: encounterPath = PathPart(2);
</script>

<div class="assignments">
  <Card>
    <EncounterSelect />
    <Content class="assignments-content">
      <div style="flex-grow: 1; margin-right: 16px">
        <EventTable encounterId={+$encounterPath} />
      </div>
      <div style="flex-shrink: 1;">
        <AssignPalette encounterId={+$encounterPath} />
      </div>
    </Content>
  </Card>
</div>

<style lang="scss" global>
  .assignments {
    .assignments-content {
      display: inline-flex;
    }
  }
</style>
