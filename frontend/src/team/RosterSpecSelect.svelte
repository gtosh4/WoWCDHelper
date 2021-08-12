<script lang="ts">
  import Button, { Icon } from "@smui/button/styled";
  import Menu from "@smui/menu/styled";
  import { Anchor } from "@smui/menu-surface";
  import List, { Item } from "@smui/list/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { classMap } from "@smui/common/classMap";
  import { CreateAnchor } from "../anchor";
  import { TeamStore } from "./team_store";
  import type { RosterMember } from "./team_api";

  export let memberId: number;
  export let encounterId: number;

  let specs: number[] = [];
  let primarySpec: number = 0;
  let selected: number | null = null;

  $: member = $TeamStore.cell(memberId, encounterId);
  $: memberInfo = member.memberInfo;
  $: rosterMember = member.rosterMember;

  $: infoP = $memberInfo.then((m) => {
    specs = m.config.specs;
    primarySpec = m.config.primarySpec;
    return m;
  });

  $: selectedP = $rosterMember.then((rm) => {
    if (rm) {
      selected = rm.spec;
    } else {
      selected = null;
    }
    return rm;
  });

  let selectedIndex;
  $: if (selected != null) {
    if (selected == 0) {
      selectedIndex = 1;
    } else {
      selectedIndex = specs.indexOf(selected);
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
      member.update(() => localMember());
    } else {
      member.remove();
    }
  }

  let anchorElem;
  const anchor = CreateAnchor();
  anchor.addClass("mdc-select__anchor");
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
        {#await Promise.all([infoP, selectedP])}
          <CircularProgress indeterminate />
        {:then [m, rm]}
          {#if rm}
            <WowIcon
              playerClass={m.classId}
              spec={selected && selected > 0 ? selected : undefined}
              height={24}
            />
          {:else}
            <Icon class="material-icons">check_box_outline_blank</Icon>
          {/if}
        {/await}
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
            {#await infoP}
              <CircularProgress indeterminate />
            {:then member}
              <WowIcon playerClass={member.classId} height={24} />
            {/await}
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
        {#await infoP}
          <Icon class="material-icons">check_box</Icon>
        {:then m}
          <WowIcon
            playerClass={m.classId}
            spec={selected && selected > 0 ? selected : undefined}
            height={24}
          />
        {/await}
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

    .class-icon {
      padding-right: 4px;
    }

    .material-icons {
      font-size: 24px;
    }

    .spec-menu {
      min-width: 24px;
    }
  }
</style>
