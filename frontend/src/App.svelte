<script lang="ts">
  import TopAppBar, { Row, Section } from "@smui/top-app-bar/styled";
  import Tab, { Label } from "@smui/tab/styled";
  import TabBar from "@smui/tab-bar/styled";
  import Textfield from "@smui/textfield/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import Roster from "./team/Roster.svelte";
  import Assignments from "./assignments/Assignments.svelte";
  import Lazy from "svelte-lazy";

  import { onMount } from "svelte";
  import { PathPart } from "./url";
  import { TeamStore } from "./team/team_store";

  const tabs = [
    { id: "roster", path: "", component: Roster },
    { id: "assignments", path: "assignments", component: Assignments },
  ];

  let active = tabs[0];

  const TabPath = PathPart(1);

  let mounted = false;
  onMount(() => {
    const tab = tabs.filter((t) => t.id == $TabPath || t.path == $TabPath);
    if (tab && tab.length > 0) active = tab[0];
    mounted = true;
  });

  $: if (mounted) {
    TabPath.set(active.id);
  }

  const team = $TeamStore.Team;
  let localName = "";
  let loadedTeam: string | null = null;
  $: if ($team && loadedTeam != $team.id) {
    localName = $team.name;
    loadedTeam = $team.id;
  }

  function save() {
    console.log("starting save", { localName });
    if ($team && $team.name != localName) {
      const t = { ...$team, name: localName };
      $team = t;
    }
  }

  function keypress(e) {
    if (e.keyCode === 13) save();
  }
</script>

<main>
  <TopAppBar variant="static" dense class="team-bar">
    <Row>
      <Section>
        {#await $team}
          <CircularProgress indeterminate />
        {:then _team}
          <Textfield
            bind:value={localName}
            on:blur={save}
            on:keyup={keypress}
            style="height: 44px"
            input$placeholder="Roster Name"
          />
        {/await}
      </Section>

      <Section align="end">
        <TabBar {tabs} let:tab bind:active style="width: unset">
          <Tab {tab} minWidth>
            <Label>{tab.id}</Label>
          </Tab>
        </TabBar>
      </Section>
    </Row>
  </TopAppBar>

  {#each tabs as { id, component } (id)}
    <div class="page-tab" class:active={active.id == id}>
      <Lazy>
        <svelte:component this={component} />
      </Lazy>
    </div>
  {/each}
</main>

<style lang="scss" global>
  main {
    width: 100%;
  }

  .team-bar {
    display: inline-flex;
    background-color: #333333;
    margin-bottom: 4px;
  }

  .page-tab {
    display: none;

    &.active {
      display: inherit;
    }
  }
</style>
