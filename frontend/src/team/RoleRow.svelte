<script lang="ts">
  import RoleCell from "./RoleCell.svelte";

  import { TeamStore } from "./team_store";

  export let roleName = "";
  export let roleType = "";

  $: encounters = $TeamStore.Encounters;
</script>

<tr
  class="role-row hover:bg-gray-50 dark-hover:bg-dark-400 border-gray-200 dark:border-gray-400 border-t border-b px-3"
>
  <td class="role-label">
    <span class="role-name">{roleName}</span>
    <span class="role-count">(0)</span>
  </td>

  <td />
  {#each $encounters || [] as enc, i (i)}
    <RoleCell {roleType} encounterId={enc.id} />
  {/each}
  <td />
</tr>

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
