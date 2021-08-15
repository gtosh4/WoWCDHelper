<script lang="ts">
  import Button, { Icon } from "@smui/button/styled";
  import List, { Item, Text } from "@smui/list/styled";
  import Menu from "@smui/menu/styled";
  import { Anchor } from "@smui/menu-surface/styled";

  import { CreateAnchor } from "../anchor";
  import { classMap } from "@smui/common/classMap";
  import { PathPart } from "../url";
  import { TeamStore } from "../team/team_store";
  import EncounterNameField from "../team/EncounterNameField.svelte";

  $: encounters = $TeamStore.Encounters;

  $: encounterPath = PathPart(2);

  function select(id: number) {
    $encounterPath = `${id}`;
  }

  $: console.info("encounter select", { ep: $encounterPath });

  let menuOpen;
  $: isSelect = $encounters ? $encounters.length > 1 : false;

  let anchorElem;
  const anchor = CreateAnchor();
</script>

<div
  class={classMap({
    "encounter-select": true,
    "mdc-select": isSelect,
    "mdc-select--activated": menuOpen,
    "smui-select--standard": isSelect,
  })}
>
  {#if isSelect}
    <div class={$anchor} use:Anchor={anchor} bind:this={anchorElem}>
      <Button on:click={() => (menuOpen = true)}>
        <Icon class="material-icons encounter-select-arrow">
          arrow_drop_down
        </Icon>
      </Button>

      <Menu
        bind:open={menuOpen}
        anchor={false}
        bind:anchorElement={anchorElem}
        class="encounter-menu"
      >
        <List
          dense
          role="listbox"
          selectedIndex={$encounters.findIndex((e) => e.id == +$encounterPath)}
        >
          {#each $encounters as enc}
            <Item
              on:SMUI:action={() => select(enc.id)}
              selected={enc.id == +$encounterPath}
            >
              <Text>{enc.name}</Text>
            </Item>
          {/each}
        </List>
      </Menu>
    </div>
  {/if}
  <EncounterNameField
    encounterId={+$encounterPath}
    input$placeholder="Encounter"
  />
</div>

<style lang="scss" global>
  .encounter-select {
    display: inline-flex;

    .encounter-select-arrow {
      color: white;
    }

    button.mdc-button {
      min-width: unset;
      padding-left: 4px;
      padding-right: 4px;
    }

    .material-icons {
      font-size: 24px;
    }

    .mdc-select__anchor {
      width: unset;
    }
  }
</style>
