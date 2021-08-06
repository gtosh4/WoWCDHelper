<script lang="ts">
  import List, { Item } from "@smui/list/styled";
  import Checkbox from "@smui/checkbox/styled";
  import Radio from "@smui/radio/styled";
  import LinearProgress from "@smui/linear-progress/styled";

  import { ClassSpecs } from "./class_specs";
  import Specilization from "./Specilization.svelte";

  export let specs: number[] = [];
  export let classId: number = 0;
  export let primarySpec: number = 0;

  $: classSpecs = ClassSpecs(classId);
</script>

{#await classSpecs}
  <LinearProgress indeterminate />
{:then specList}
  <List dense>
    {#each specList as spec (spec.id)}
      <Item selected={specs.indexOf(spec.id) >= 0}>
        <Checkbox bind:group={specs} value={spec.id} />
        <Specilization specId={spec.id} />
        <Radio
          bind:group={primarySpec}
          value={spec.id}
          disabled={specs.indexOf(spec.id) < 0}
        />
      </Item>
    {/each}
  </List>
{/await}
