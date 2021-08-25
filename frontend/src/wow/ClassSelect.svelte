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
    showList = false;
  }

  let showList = false;
</script>

<div>
  {#await ClassList}
    <ProgressLinear />
  {:then cs}
    <Select dense autocomplete {showList}>
      <Button slot="select" on:click={() => (showList = !showList)}>
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
</style>
