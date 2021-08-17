<script lang="ts">
  import AppBar from "smelte/src/components/AppBar";
  import { Tabs, Tab } from "smelte/src/components/Tabs";
  import Textfield from "smelte/src/components/TextField";
  import ProgressCircular from "smelte/src/components/ProgressCircular";
  import Roster from "./team/Roster.svelte";
  import Assignments from "./assignments/Assignments.svelte";

  import { onMount } from "svelte";
  import { PathPart } from "./url";
  import { TeamStore } from "./team/team_store";

  const tabs = [
    { id: "roster", path: "", text: "Roster" },
    {
      id: "assignments",
      path: "assignments",
      text: "Assignments",
    },
  ];

  let selected = tabs[0].id;

  const TabPath = PathPart(1);

  let mounted = false;
  onMount(() => {
    const tab = tabs.filter((t) => t.id == $TabPath || t.path == $TabPath);
    if (tab && tab.length > 0) selected = tab[0].id;
    mounted = true;
  });

  $: if (mounted) {
    TabPath.set(selected);
  }

  const team = $TeamStore.Team;
  let localName = "";
  let loadedTeam: string | null = null;
  $: if ($team && loadedTeam != $team.id) {
    localName = $team.name;
    loadedTeam = $team.id;
  }

  function save() {
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
  <AppBar class="team-bar static inline-flex">
    <section class="flex-initial">
      {#await $team}
        <ProgressCircular />
      {:then _team}
        <Textfield
          bind:value={localName}
          on:blur={save}
          on:keyup={keypress}
          placeholder="Roster Name"
          dense
        />
      {/await}
    </section>

    <section class="flex-initial">
      <Tabs items={tabs} bind:active={selected} class="w-56" />
    </section>
  </AppBar>

  {#if selected == "roster"}
    <Roster />
  {:else if selected == "assignments"}
    <Assignments />
  {/if}
</main>

<style lang="scss" global>
  main {
    width: 100%;
  }

  :root {
    // https://wow.tools/dbc/?dbc=chrclasses&build=9.1.0.39653#page=1

    --wow-class-1: rgb(198, 155, 109);
    --wow-class-warrior: var(--wow-class-1);
    --wow-class-2: rgb(244, 140, 186);
    --wow-class-paladin: var(--wow-class-2);
    --wow-class-3: rgb(170, 211, 114);
    --wow-class-hunter: var(--wow-class-3);
    --wow-class-4: rgb(255, 244, 104);
    --wow-class-rogue: var(--wow-class-4);
    --wow-class-5: rgb(255, 255, 255);
    --wow-class-priest: var(--wow-class-5);
    --wow-class-6: rgb(196, 30, 58);
    --wow-class-death-knight: var(--wow-class-6);
    --wow-class-7: rgb(0, 112, 221);
    --wow-class-shaman: var(--wow-class-7);
    --wow-class-8: rgb(63, 199, 235);
    --wow-class-mage: var(--wow-class-8);
    --wow-class-9: rgb(135, 136, 238);
    --wow-class-warlock: var(--wow-class-9);
    --wow-class-10: rgb(0, 255, 152);
    --wow-class-monk: var(--wow-class-10);
    --wow-class-11: rgb(255, 124, 10);
    --wow-class-druid: var(--wow-class-11);
    --wow-class-12: rgb(163, 48, 201);
    --wow-class-demon-hunter: var(--wow-class-12);
  }
</style>
