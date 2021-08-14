<script lang="ts">
  import Card, { Content } from "@smui/card/styled";
  import EncounterSelect from "./EncounterSelect.svelte";
  import EventTable from "./EventTable.svelte";

  import { PathPart } from "../url";

  $: encounterPath = PathPart(2);
</script>

<div class="assignments">
  <Card>
    <EncounterSelect />
    <Content>
      <EventTable encounterId={+$encounterPath} />
    </Content>
  </Card>
</div>
