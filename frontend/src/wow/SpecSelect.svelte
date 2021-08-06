<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table";
  import Radio from "@smui/radio/styled";
  import Checkbox from "@smui/checkbox/styled";
  import LinearProgress from "@smui/linear-progress/styled";

  import Specilization from "./Specilization.svelte";

  import { ClassSpecs } from "./class_specs";
  import type { Specialization as APISpecialization } from "./api";

  export let classId: number = 0;
  export let specs: number[] = [];
  export let primarySpec: number = 0;

  let classSpecs: APISpecialization[] = [];
  let loaded = false;
  $: {
    classSpecs = [];
    loaded = false;

    ClassSpecs(classId).then((ss) => {
      classSpecs = ss;
      loaded = true;
    });
  }

  $: {
    if (specs.indexOf(primarySpec) < 0 && specs.length > 0) {
      primarySpec = specs[0];
    }
  }
</script>

<DataTable class="spec-select">
  <Head>
    <Row>
      <Cell checkbox><Checkbox /></Cell>
      <Cell style="width: 100%" />
      <Cell>Primary</Cell>
    </Row>
  </Head>

  <Body>
    {#each classSpecs as spec (spec.id)}
      <Row>
        <Cell checkbox>
          <Checkbox
            bind:group={specs}
            checked={specs.indexOf(spec.id) >= 0}
            value={spec.id}
          />
        </Cell>
        <Cell>
          <Specilization specId={spec.id} />
        </Cell>
        <Cell>
          <Radio
            bind:group={primarySpec}
            value={spec.id}
            disabled={specs.indexOf(spec.id) < 0}
          />
        </Cell>
      </Row>
    {/each}
  </Body>
  <LinearProgress bind:closed={loaded} indeterminate slot="progress" />
</DataTable>

<style lang="scss">
  :global(.spec-select) {
    width: 100%;
  }
</style>
