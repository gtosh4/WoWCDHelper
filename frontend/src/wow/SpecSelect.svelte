<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table";
  import Radio from "@smui/radio/styled";
  import Checkbox from "@smui/checkbox/styled";
  import LinearProgress from "@smui/linear-progress/styled";

  import Specilization from "./Specilization.svelte";

  import { ClassSpecs } from "./class_specs";

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

  function toggleSelectAll() {
    if (selectedAll == null || selectedAll == true) {
      specs.length = 0;
    } else {
      specInfo.then((infos) => {
        specs.length = 0;
        specs.push(...infos.map((i) => i.id));
      });
    }
  }
</script>

<DataTable class="spec-select">
  <Head>
    <Row>
      <Cell>
        <Checkbox
          bind:value={selectedAll}
          indeterminate={selectedAll == null}
          on:change={toggleSelectAll}
        />
      </Cell>
      <Cell style="width: 100%" />
      <Cell>Primary</Cell>
    </Row>
  </Head>

  <Body>
    {#await specInfo then classSpecs}
      {#each classSpecs as spec (spec.id)}
        <Row>
          <Cell>
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
    {/await}
  </Body>
  <LinearProgress bind:closed={loaded} indeterminate slot="progress" />
</DataTable>

<style lang="scss">
  :global(.spec-select) {
    width: 100%;
  }
</style>
