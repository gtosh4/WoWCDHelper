<script lang="ts">
  import Card from "@smui/card/styled";
  import List, { Item } from "@smui/list/styled";
  import Member from "./Member.svelte";

  import { createEventDispatcher } from "svelte";
  import { RoleMembers } from "./role_members";

  export let roleName = "Tank";

  const dispatch = createEventDispatcher();

  const members = RoleMembers(roleName);
</script>

<Card {...$$restProps}>
  <h2>{roleName}</h2>
  <List dense>
    {#each $members as memberId}
      <Item on:SMUI:action={() => dispatch("edit", memberId)}>
        <Member {memberId} />
      </Item>
    {/each}
  </List>
</Card>

<style lang="scss">
</style>
