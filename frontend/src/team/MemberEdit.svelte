<script lang="ts">
  import Card, { Actions, ActionButtons } from "@smui/card/styled";
  import Textfield from "@smui/textfield/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import ClassSelect from "../wow/ClassSelect.svelte";

  import { CurrentTeam, TeamMember } from "./api";
  import type { Member } from "./api";
  import { createEventDispatcher } from "svelte";
  import SpecSelect from "../wow/SpecSelect.svelte";

  export let memberId: number | undefined = undefined;

  const dispatch = createEventDispatcher();

  $: member = memberId ? TeamMember(memberId) : undefined;

  let localId,
    localName = "",
    localClassId = 0,
    localSpecs: number[] = [],
    localPrimarySpec = 0;

  $: if ($member) {
    $member.then((m) => {
      if (localId != m.id) {
        localId = m.id;
        localName = m.name;
        localClassId = m.classId;
        localSpecs = m.config.specs;
        localPrimarySpec = m.config.primarySpec;
      }
    });
  }

  function save() {
    const newMember: Member = {
      id: localId,
      name: localName,
      classId: localClassId,
      config: {
        specs: localSpecs,
        primarySpec: localPrimarySpec,
      },
    };

    if ($member) {
      $member = newMember;
    } else {
      CurrentTeam.addMember(newMember);
    }
    cancel();
  }

  function cancel() {
    localId = undefined;
    localClassId = 0;
    localName = "";
    localClassId = 0;
    localSpecs = [];
    localPrimarySpec = 0;

    dispatch("close");
  }
</script>

<Card padded class="member-edit">
  {#await $member}
    <LinearProgress indeterminate />
  {/await}
  <Textfield
    style="width: 100%;"
    helperLine$style="width: 100%;"
    bind:value={localName}
    label="Name"
  />

  <ClassSelect bind:value={localClassId} />

  <SpecSelect
    classId={localClassId}
    bind:specs={localSpecs}
    bind:primarySpec={localPrimarySpec}
  />

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
