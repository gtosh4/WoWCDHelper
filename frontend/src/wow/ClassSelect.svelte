<script lang="ts">
  import Button, { Label, Icon } from "@smui/button/styled";
  import Menu from "@smui/menu/styled";
  import List, { Item } from "@smui/list/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import PlayerClass from "./PlayerClass.svelte";

  import { ClassList } from "./api";
  import { classMap } from "@smui/common/classMap";
  import WowIcon from "./WowIcon.svelte";

  export let value = 0;
  export let small = false;

  let menuOpen = false;

  function select(id) {
    value = id;
    menuOpen = false;
  }
</script>

<div
  class={classMap({
    "class-select": true,
    "mdc-select": true,
    "mdc-select--activated": menuOpen,
    "smui-select--standard": true,
  })}
>
  <div class="anchor mdc-select__anchor">
    <Button on:click={() => (menuOpen = true)}>
      {#if small}
        <WowIcon playerClass={value} height={24} />
      {:else}
        {#if value > 0}
          <PlayerClass playerClass={value} />
        {:else}
          <Label>Class</Label>
        {/if}
        <Icon class="material-icons" style="margin: 0;">arrow_drop_down</Icon>
      {/if}
    </Button>
  </div>

  <Menu bind:open={menuOpen}>
    {#await ClassList}
      <LinearProgress indeterminate />
    {:then cs}
      <List
        dense
        role="listbox"
        selectedIndex={cs.map((c) => c.id).indexOf(value)}
      >
        {#each cs as cls (cls.id)}
          <Item
            on:SMUI:action={() => select(cls.id)}
            selected={cls.id == value}
          >
            <PlayerClass playerClass={cls.id} />
          </Item>
        {/each}
      </List>
    {/await}
  </Menu>
</div>

<style lang="scss" global>
  .class-select {
    .class-icon {
      padding-right: 4px;
    }
  }
</style>
