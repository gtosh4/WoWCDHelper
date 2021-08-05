<script lang="ts">
  import LayoutGrid, { Cell } from "@smui/layout-grid/styled";
  import Button, { Label, Icon } from "@smui/button/styled";
  import Dialog from "@smui/dialog/styled";

  import RosterRole from "./RosterRole.svelte";
  import MemberEdit from "./MemberEdit.svelte";

  const roles = ["Tank", "Healer", "Damage"];

  let open = false;
  let memberId: number | undefined;

  function showEdit(id?: number) {
    open = true;
    memberId = id;
  }

  function close() {
    open = false;
    memberId = undefined;
  }
</script>

<LayoutGrid>
  {#each roles as roleName (roleName)}
    <Cell align="top">
      <RosterRole
        {roleName}
        padded
        on:edit={(event) => showEdit(event.detail)}
      />
    </Cell>
  {/each}
</LayoutGrid>
<Button on:click={() => showEdit(undefined)}>
  <Icon class="material-icons">add</Icon>
  <Label>Add Member</Label>
</Button>
<Dialog bind:open class="add-member">
  <MemberEdit on:close={close} {memberId} />
</Dialog>

<style lang="scss">
  :global(.add-member) {
    overflow-y: visible;

    :global(.mdc-dialog__surface) {
      overflow-y: visible;
    }
  }
</style>
