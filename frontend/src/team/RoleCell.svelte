<script lang="ts">
  import { Cell } from "@smui/data-table/styled";
  import CircularProgress from "@smui/circular-progress/styled";

  import { RoleMembers } from "./member_filter";
  import { TeamStore } from "./team_store";

  export let roleType = "";
  export let encounterId: number;

  $: column = $TeamStore.column(encounterId);
  $: members = column.members;

  let count: Promise<number>;
  $: if ($members) {
    count = RoleMembers($members, roleType).then((rms) =>
      rms ? rms.length : 0
    );
  }
</script>

<Cell class="role-count-cell">
  {#await count}
    <CircularProgress indeterminate />
  {:then count}
    <span>{count}</span>
  {/await}
</Cell>
