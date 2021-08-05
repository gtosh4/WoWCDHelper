<script lang="ts">
  import Card, { Content } from "@smui/card/styled";
  import List, { Item } from "@smui/list/styled";
  import Member from "./Member.svelte";

  import { writable } from "svelte/store";
  import { Spec } from "../wow/api";
  import { CurrentTeam, SortMemberIds } from "./api";
  import { createEventDispatcher } from "svelte";

  export let roleName = "Tank";

  const dispatch = createEventDispatcher();

  const team = CurrentTeam;
  let members = writable(new Set<number>());

  $: [...$team.values()].map((m) => {
    if (!m.config.primary_spec) return;

    const spec = Spec(m.config.primary_spec);
    spec.subscribe((s) => {
      if (!s) return;

      members.update((ms) => {
        if (s.role.name == roleName) {
          ms.add(m.id);
        } else {
          ms.delete(m.id);
        }
        return ms;
      });
    });
  });

  $: sortedMembers = [...$members].sort(SortMemberIds);
</script>

<Card {...$$restProps}>
  <h2>{roleName}</h2>
  <List dense>
    {#each sortedMembers as memberId}
      <Item on:SMUI:action={() => dispatch("edit", memberId)}>
        <Member {memberId} />
      </Item>
    {/each}
  </List>
</Card>

<style lang="scss">
</style>
