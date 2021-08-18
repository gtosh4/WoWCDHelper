<script lang="ts">
  import ProgressCircular from "smelte/src/components/ProgressCircular";

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

<td class="role-count-cell relative p-3 font-normal text-right">
  {#await count}
    <ProgressCircular size={26} />
  {:then count}
    <span>{count}</span>
  {/await}
</td>
