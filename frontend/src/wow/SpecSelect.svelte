<script lang="ts">
  import RadioButton from "smelte/src/components/RadioButton";
  import Checkbox from "smelte/src/components/Checkbox";
  import Icon from "smelte/src/components/Icon";
  import ProgressLinear from "smelte/src/components/ProgressLinear";
  import Specilization from "./Specilization.svelte";
  import Ripple from "smelte/src/components/Ripple";
  import Label from "smelte/src/components/Checkbox/Label.svelte";
  import { slide } from "svelte/transition";

  import { ClassSpecs } from "./class_specs";
  import PrimarySpecRadio from "./PrimarySpecRadio.svelte";

  export let classId: number = 0;
  export let specs: number[] = [];
  export let primarySpec: number = 0;

  $: specInfo = ClassSpecs(classId);

  let loaded = false;
  $: {
    loaded = false;
    specInfo.then(() => (loaded = true));
  }

  $: if (specs.indexOf(primarySpec) < 0) {
    if (specs.length > 0) {
      primarySpec = specs[0];
    } else {
      primarySpec = 0;
    }
  }

  let selectedAll: boolean | null = false;
  $: specInfo.then((infos) => {
    if (infos.length == specs.length) {
      selectedAll = true;
    } else if (specs.length == 0) {
      selectedAll = false;
    } else {
      selectedAll = null;
    }
  });

  const thClass =
    "encounter-header capitalize duration-100 text-gray-600 text-xs hover:text-black dark-hover:text-white p-3 font-normal text-right";
  const trClass =
    "hover:bg-gray-50 dark-hover:bg-dark-400 border-gray-200 dark:border-gray-400 border-t border-b px-3";
  const tdClass = "relative p-3 font-normal text-right";

  function toggle(specId: number) {
    const idx = specs.indexOf(specId);
    if (idx >= 0) {
      specs.splice(idx, 1);
    } else {
      specs.push(specId);
    }
    specs = specs;
  }

  function toggleSelectAll() {
    if (selectedAll == null || selectedAll == true) {
      specs.length = 0;
    } else {
      specInfo.then((infos) => {
        specs.length = 0;
        specs.push(...infos.map((i) => i.id));
      });
    }
    specs = specs;
  }
  const icons = new Map([
    [true, "check_box"],
    [false, "check_box_outline_blank"],
    [null, "indeterminate_check_box"],
  ]);

  $: console.info("spec sel", {
    specInfo,
    loaded,
    selectedAll,
    specs,
    primarySpec,
  });
</script>

<table
  class="spec-select shadow relative text-sm overflow-x-auto dark:bg-dark-500"
>
  <thead>
    <th class={thClass}>
      <Icon on:click={toggleSelectAll} class="cursor-pointer">
        {icons.get(selectedAll)}
      </Icon>
    </th>
    <th class={thClass + " w-full"} />
    <th class={thClass}>Primary</th>
  </thead>
  {#if !loaded}
    <div class="absolute w-full" transition:slide>
      <ProgressLinear />
    </div>
  {/if}

  <tbody>
    {#await specInfo then classSpecs}
      {#each classSpecs as spec (spec.id)}
        <tr class={trClass}>
          <td class={tdClass}>
            <Checkbox
              checked={specs.indexOf(spec.id) >= 0}
              on:change={() => toggle(spec.id)}
            />
          </td>
          <td class={tdClass}>
            <Specilization specId={spec.id} />
          </td>
          <td>
            <PrimarySpecRadio bind:primarySpec bind:specs thisSpec={spec.id} />
          </td>
        </tr>
      {/each}
    {/await}
  </tbody>
</table>

<style lang="scss" global>
  .spec-select {
    width: 100%;
  }
</style>
