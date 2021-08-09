<script lang="ts">
  import Textfield from "@smui/textfield/styled";
  import { Encounters } from "./encounters_api";

  export let encounterId: number;

  let loadedEncounter = false;
  let localName = "";

  $: encounter = Encounters.encounter(encounterId);
  $: $encounter.then((e) => {
    if (!loadedEncounter) {
      localName = e.name || "";
      loadedEncounter = true;
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
