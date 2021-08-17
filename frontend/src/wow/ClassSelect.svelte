<script lang="ts">
  import Button from "smelte/src/components/Button";
  import ProgressLinear from "smelte/src/components/ProgressLinear";
  import Select from "smelte/src/components/Select";
  import { ListItem } from "smelte/src/components/List";

  import PlayerClass from "./PlayerClass.svelte";

  import { ClassList } from "./api";

  export let value = 0;

  function select(id) {
    value = id;
  }

  let showList = false;
</script>

<div class="class-select">
  {#await ClassList}
    <ProgressLinear />
  {:then cs}
    <Select dense autocomplete {showList}>
      <Button slot="select" on:click={() => (showList = true)}>
        <PlayerClass playerClass={value} />
      </Button>
      <div
        slot="options"
        class="absolute left-0 bg-white rounded shadow w-full z-20 dark:bg-dark-500"
      >
        <ul class="py-2 rounded">
          {#each cs as cls (cls.id)}
            <ListItem
              selected={cls.id == value}
              on:click={() => select(cls.id)}
            >
              <PlayerClass playerClass={cls.id} />
            </ListItem>
          {/each}
        </ul>
      </div>
    </Select>
  {/await}
</div>

<style lang="scss" global>
  .class-select {
    .class-icon {
      padding-right: 4px;
    }
  }
</style>
