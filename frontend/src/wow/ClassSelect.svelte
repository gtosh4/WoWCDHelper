<script lang="ts">
  import Button, { Label, Icon } from "@smui/button/styled";
  import Menu from "@smui/menu/styled";
  import List, { Item } from "@smui/list/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import PlayerClass from "./PlayerClass.svelte";

  import { ClassList } from "./api";
  import { classMap } from "@smui/common/classMap";

  export let value = 0;

  let menuOpen = false;

  function select(id) {
    value = id;
    menuOpen = false;
  }

  $: selectedIndex = ClassList.then((l) =>
    l.map((cls) => cls.id).indexOf(value)
  );
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
      {#if value > 0}
        <PlayerClass playerClass={value} />
      {:else}
        <Label>Class</Label>
      {/if}
      <Icon class="material-icons" style="margin: 0;">arrow_drop_down</Icon>
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
        {#each cs as cls, idx (cls.id)}
          <Item on:SMUI:action={() => select(cls.id)} selected={idx == value}>
            <PlayerClass playerClass={cls.id} />
          </Item>
        {/each}
      </List>
    {/await}
  </Menu>
</div>

<style lang="scss">
</style>
