<script lang="ts">
  import Textfield from "@smui/textfield/styled";
  import { TeamStore } from "./team_store";

  export let encounterId: number;

  let loadedEncounter: number | null = null;
  let localName = "";

  $: column = $TeamStore.column(encounterId);
  $: encounter = column.encounter;

  $: if ($encounter && loadedEncounter != $encounter.id) {
    localName = $encounter.name || "";
    loadedEncounter = $encounter.id;
  }

  $: console.info("enf", { encounterId, localName, loadedEncounter });

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
  {...$$restProps}
/>

<style lang="scss" global>
  .encounter-name {
    height: auto;
  }
</style>
