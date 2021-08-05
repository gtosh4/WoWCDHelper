<script lang="ts">
  import Tab, { Label } from "@smui/tab/styled";
  import TabBar from "@smui/tab-bar/styled";
  import Roster from "./team/Roster.svelte";
  import Assignments from "./assignments/Assignments.svelte";

  import { onMount } from "svelte";
  import { HashPathPart } from "./url";

  const tabs = [
    { id: "roster", path: "/#/", component: Roster },
    { id: "assignments", path: "/#/assignments", component: Assignments },
  ];

  let active = tabs[0];

  const TabPath = HashPathPart(1);

  let mounted = false;
  onMount(() => {
    const tab = tabs.filter((t) => t.id == $TabPath);
    if (tab && tab.length > 0) active = tab[0];
    mounted = true;
  });

  $: if (mounted) {
    TabPath.set(active.id);
  }
</script>

<main>
  <TabBar {tabs} let:tab bind:active>
    <Tab {tab} minWidth>
      <Label>{tab.id}</Label>
    </Tab>
  </TabBar>

  {#each tabs as { id, component } (id)}
    <div class="page-tab" class:active={active.id == id}>
      <svelte:component this={component} />
    </div>
  {/each}
</main>

<style lang="scss">
  main {
    width: 100%;
  }

  .page-tab {
    display: none;

    &.active {
      display: inherit;
    }
  }
</style>
