<script lang="ts">
  import Button, { Icon } from "@smui/button/styled";
  import Menu from "@smui/menu/styled";
  import { Anchor } from "@smui/menu-surface/styled";
  import List, { Item } from "@smui/list/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { classMap } from "@smui/common/classMap";
  import { CreateAnchor } from "../anchor";
  import { TeamStore } from "./team_store";
  import type { RosterMember } from "./team_api";
  import { LoadingState } from "../store_helpers";

  export let memberId: number;
  export let encounterId: number;

  let specs: number[] = [];
  let primarySpec: number = 0;
  let selected: number | null = null;

  $: memberInfo = $TeamStore.row(memberId).member;
  $: cell = $TeamStore.cell(memberId, encounterId);
  $: rosterMember = cell.rosterMember;

  $: if (rosterMember.state == LoadingState.Uninitialized) {
    $TeamStore.row(memberId).encounterAPI.get();
  }

  $: if ($memberInfo) {
    specs = $memberInfo.config.specs;
    primarySpec = $memberInfo.config.primarySpec;
  }

  $: if ($rosterMember) {
    selected = $rosterMember.spec;
  } else {
    selected = null;
  }

  let selectedIndex;
  $: if (selected != null) {
    if (selected == 0) {
      selectedIndex = 1;
    } else {
      selectedIndex = specs.indexOf(selected) + 2;
    }
  } else {
    selectedIndex = 0;
  }

  let menuOpen = false;

  function localMember(): RosterMember {
    return {
      encounter_id: encounterId,
      member_id: memberId,
      spec: selected,
    };
  }

  function toggle() {
    if (selected) {
      select(null);
    } else {
      select(primarySpec);
    }
  }

  function select(id: number | null) {
    selected = id;
    menuOpen = false;

    if (id != null) {
      cell.rosterMember.update(() => localMember());
    } else {
      cell.rosterMember.remove();
    }
  }

  let anchorElem;
  const anchor = CreateAnchor();
</script>

<div
  class={classMap({
    "roster-spec-select": true,
    "mdc-select": specs.length > 1,
    "mdc-select--activated": menuOpen,
    "smui-select--standard": specs.length > 1,
  })}
>
  {#if specs.length > 1}
    <div class={$anchor} use:Anchor={anchor} bind:this={anchorElem}>
      <Button on:click={() => (menuOpen = true)}>
        <div class="spec-select-anchor">
          {#if $memberInfo == undefined}
            <CircularProgress indeterminate />
          {:else if $rosterMember}
            <WowIcon
              playerClass={$memberInfo.classId}
              spec={selected && selected > 0 ? selected : undefined}
              height={24}
            />
            <Icon class="material-icons spec-select-arrow">arrow_drop_down</Icon
            >
          {:else}
            <Icon class="material-icons">check_box_outline_blank</Icon>
            <Icon class="material-icons spec-select-arrow">arrow_drop_down</Icon
            >
          {/if}
        </div>
      </Button>

      <Menu
        bind:open={menuOpen}
        anchor={false}
        bind:anchorElement={anchorElem}
        class="spec-menu"
      >
        <List dense role="listbox" {selectedIndex}>
          <Item on:SMUI:action={() => select(null)}>
            <Icon class="material-icons">backspace</Icon>
          </Item>
          <Item on:SMUI:action={() => select(0)}>
            {#if $memberInfo == undefined}
              <CircularProgress indeterminate />
            {:else}
              <WowIcon playerClass={$memberInfo.classId} height={24} />
            {/if}
          </Item>
          {#each specs as spec}
            <Item
              on:SMUI:action={() => select(spec)}
              selected={spec == selected}
            >
              <WowIcon {spec} height={24} />
            </Item>
          {/each}
        </List>
      </Menu>
    </div>
  {:else}
    <Button on:click={toggle}>
      {#if selected}
        {#if $memberInfo == undefined}
          <Icon class="material-icons">check_box</Icon>
        {:else}
          <WowIcon
            playerClass={$memberInfo.classId}
            spec={selected && selected > 0 ? selected : undefined}
            height={24}
          />
        {/if}
      {:else}
        <Icon class="material-icons">check_box_outline_blank</Icon>
      {/if}
    </Button>
  {/if}
</div>

<style lang="scss" global>
  .roster-spec-select {
    button.mdc-button {
      min-width: 24px;
    }

    div.mdc-select__anchor {
      height: auto;
      width: auto;
    }

    .class-icon {
      padding-right: 4px;
    }

    .material-icons {
      font-size: 24px;
      height: 24px;
    }

    .spec-menu {
      min-width: 24px;
    }

    .spec-select-anchor {
      display: inline-flex;
    }

    .spec-select-arrow {
      border-left-width: 1px;
      border-left-style: solid;
      color: rgb(200, 200, 200);
      height: 100%;
      margin-left: 2px;
    }
  }
</style>
