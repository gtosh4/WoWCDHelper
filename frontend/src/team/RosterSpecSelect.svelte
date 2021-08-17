<script lang="ts">
  import { TeamStore } from "./team_store";
  import Icon from "smelte/src/components/Icon";
  import { ListItem } from "smelte/src/components/List";
  import ProgressCircular from "smelte/src/components/ProgressCircular";
  import WowIcon from "../wow/WowIcon.svelte";
  import { fly } from "svelte/transition";
  import { quadOut, quadIn } from "svelte/easing";

  export let memberId: number;
  export let encounterId: number;

  const cell = $TeamStore.cell(memberId, encounterId);
  const rosterMember = cell.rosterMember;
  rosterMember.init();
  const memberInfo = cell.row.member;
  memberInfo.init();

  $: spec = $rosterMember ? $rosterMember.spec : null;

  let showMenu = false;
  $: useMenu = $memberInfo && $memberInfo.config.specs.length > 1;

  function select(selected: number | null) {
    if (selected == null) {
      rosterMember.remove();
    } else {
      rosterMember.set({
        member_id: memberId,
        encounter_id: encounterId,
        spec: selected,
      });
    }
  }

  function toggle(evt) {
    evt.stopPropagation();
    if (useMenu) {
      showMenu = true;
    } else {
      if (spec) {
        select(null);
      } else {
        select($memberInfo ? $memberInfo.config.primarySpec : 0);
      }
    }
  }

  const inProps = { y: 10, duration: 200, easing: quadIn };
  const outProps = { y: -10, duration: 100, easing: quadOut, delay: 100 };
</script>

<div class="cursor-pointer relative flex items-center">
  {#if !$memberInfo}
    <ProgressCircular />
  {:else}
    {#if spec == null}
      <Icon on:click={toggle}>check_box_outline_blank</Icon>
    {:else}
      <WowIcon playerClass={$memberInfo.classId} {spec} on:click={toggle} />
    {/if}
    {#if showMenu}
      <div
        class="absolute w-auto top-0 rounded shadow z-20 bg-white dark:bg-dark-500"
        in:fly={inProps}
        out:fly={outProps}
      >
        <ul class="py-2 rounded">
          <ListItem dense on:click={() => select(null)}>
            <Icon>backspace</Icon>
          </ListItem>

          <ListItem dense on:click={() => select(0)}>
            <WowIcon playerClass={$memberInfo.classId} />
          </ListItem>

          {#each $memberInfo.config.specs as spec}
            <ListItem dense on:click={() => select(spec)}>
              <WowIcon playerClass={$memberInfo.classId} {spec} />
            </ListItem>
          {/each}
        </ul>
      </div>
    {/if}
  {/if}
</div>
<svelte:window on:click={() => (showMenu = false)} />

<style lang="scss" global>
</style>
