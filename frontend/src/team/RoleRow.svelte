<script lang="ts">
  import { Row, Cell } from "@smui/data-table/styled";
  import RoleCell from "./RoleCell.svelte";

  import { TeamStore } from "./team_store";

  export let roleName = "";
  export let roleType = "";

  $: encounters = $TeamStore.Encounters;
</script>

<Row class="role-row">
  <Cell class="role-label">
    <span class="role-name">{roleName}</span>
    <span class="role-count">(0)</span>
  </Cell>

  <Cell />
  {#each $encounters as enc, i (i)}
    <RoleCell {roleType} encounterId={enc.id} />
  {/each}
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
