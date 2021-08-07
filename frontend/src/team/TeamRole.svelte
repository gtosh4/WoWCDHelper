<script lang="ts">
  import Card from "@smui/card/styled";
  import List, { Item } from "@smui/list/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import Member from "./Member.svelte";

  import { createEventDispatcher } from "svelte";
  import { RoleMembers } from "./role_members";

  export let roleName = "Tank";

  const dispatch = createEventDispatcher();

  const members = RoleMembers(roleName);
</script>

<Card {...$$restProps}>
  <h2>
    {roleName}
    {#await $members then ms}
      ({ms.length})
    {/await}
  </h2>
  {#await $members}
    <LinearProgress indeterminate />
  {:then ms}
    <List dense>
      {#each ms as memberId}
        <Item on:SMUI:action={() => dispatch("edit", memberId)}>
          <Member {memberId} />
        </Item>
      {/each}
    </List>
  {/await}
</Card>

<style lang="scss">
</style>
