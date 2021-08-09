<script lang="ts">
  import Textfield from "@smui/textfield/styled";
  import { Encounters } from "./encounters_api";

  export let encounterId: number;

  let loadedEncounter: number | null = null;
  let localName = "";

  $: encounter = Encounters.encounter(encounterId);
  $: $encounter.then((e) => {
    if (loadedEncounter != e.id) {
      localName = e.name || "";
      loadedEncounter = e.id;
    }
  });

  function save() {
    encounter.update((e) => {
      if (e.name != localName) {
        e.name = localName;
        return e;
      }
      return undefined;
    });
  }

  function keypress(e) {
    if (e.keyCode === 13) save();
  }
</script>

<Textfield
  class="encounter-name"
  bind:value={localName}
  on:blur={() => save()}
  on:keyup={keypress}
/>
