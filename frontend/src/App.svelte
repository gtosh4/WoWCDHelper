<script lang="ts">
  import Tab, { Label } from "@smui/tab/styled";
  import TabBar from "@smui/tab-bar/styled";
  import Roster from "./roster/Roster.svelte";
  import Assignments from "./assignments/Assignments.svelte";

  import url from "./url";
  import { onMount } from "svelte";

  const tabs = [
    { id: "roster", path: "/#/", component: Roster },
    { id: "assignments", path: "/#/assignments", component: Assignments },
  ];

  let active = tabs[0];

  let mounted = false;
  onMount(() => {
    const path: string = $url.hash;
    const tabid = path.replace(/^#\//, "").split("/")[0];
    const tab = tabs.filter((t) => t.id == tabid);
    if (tab && tab[0]) active = tab[0];
    mounted = true;
  });

  $: if (mounted) {
    history.pushState(active.path, "", active.path);
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
