<script lang="ts">
  import { Row, Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";
  import { Encounters } from "./encounters_api";

  export let roleName = "";
  export let roleType = "";
</script>

<Row class="role-row">
  <Cell class="role-label">
    <span class="role-name">{roleName}</span>
    <span class="role-count">(0)</span>
  </Cell>

  {#await $Encounters}
    <Cell>
      <CircularProgress indeterminate />
    </Cell>
  {:then encounters}
    <Cell />
    {#each encounters as enc, i (i)}
      <Cell>
        <span>0</span>
      </Cell>
    {/each}
  {/await}
  <Cell />
</Row>

<style lang="scss" global>
  .roster-entry {
    $line-height: auto;

    height: $line-height;

    td.mdc-data-table__cell {
      height: $line-height;
    }

    .roster-member {
      display: flex;

      .name {
        $width: 12em;

        width: $width;
        height: auto;

        &.name--hovered {
          width: calc(#{$width} - 18px - 8px - 8px);
        }
      }
      .configure {
        display: none;
        min-width: unset;
        width: calc(18px+8px+8px);
        height: $line-height;

        &.configure--active {
          display: unset;
        }
      }
    }
  }
</style>
