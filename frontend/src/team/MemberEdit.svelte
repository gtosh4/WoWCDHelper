<script lang="ts">
  import Card, { Actions, ActionButtons } from "@smui/card/styled";
  import Textfield from "@smui/textfield/styled";
  import Select, { Option } from "@smui/select/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { TeamMember } from "./api";
  import { Classes, SortClassByName } from "../wow/api";
  import { createEventDispatcher } from "svelte";

  export let memberId: number | undefined = undefined;

  $: console.log("edit1", memberId);

  const dispatch = createEventDispatcher();

  let memberName = "";
  let classId = ""; // Use strings: https://github.com/hperrin/svelte-material-ui/issues/252

  $: member = memberId ? TeamMember(memberId) : undefined;

  $: if ($member) {
    memberName = $member.name;
  }
  $: if ($member) {
    classId = `${$member.classId}`;
  }

  $: console.log("edit2", { memberName, classId, member: $member });

  function save() {
    if ($member) {
      member.update((m) => {
        m.name = memberName;
        m.classId = +classId;
        return m;
      });
    } else {
    }
    dispatch("close");
  }

  function cancel() {
    dispatch("close");
  }
</script>

<Card padded class="member-edit">
  <Textfield bind:value={memberName} label="Name" />

  <Select
    key={(m) => (m && m.id) || 0}
    bind:value={classId}
    label="Class"
    list$dense
  >
    <WowIcon
      slot="leadingIcon"
      playerClass={classId}
      class="class-icon"
      height={24}
    />
    {#await $Classes then classes}
      {#each [...classes.values()].sort(SortClassByName) as cls (cls.id)}
        <Option value={`${cls.id}`}>
          <WowIcon playerClass={cls.id} class="class-icon" height={24} />
          {cls.name}
        </Option>
      {/each}
    {/await}
  </Select>

  <Actions>
    <ActionButtons>
      <Button on:click={save}>
        {#if $$slots.saveLabel}
          <slot name="saveLabel" />
        {:else}
          <Icon class="material-icons">save</Icon>
          <Label>Save</Label>
        {/if}
      </Button>
      <Button on:click={cancel}>
        {#if $$slots.cancelLabel}
          <slot name="cancelLabel" />
        {:else}
          <Label>Cancel</Label>
        {/if}
      </Button>
    </ActionButtons>
  </Actions>
</Card>

<style lang="scss">
  :global(.class-icon) {
    padding-right: 4px;
  }

  :global(.member-edit) {
    :global(.smui-select--standard.mdc-select--with-leading-icon
        .mdc-select__anchor) {
      padding-left: 0;
    }
  }
</style>
