<script lang="ts">
  import Card, { Actions, ActionButtons } from "@smui/card/styled";
  import Textfield from "@smui/textfield/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import ClassSelect from "../wow/ClassSelect.svelte";
  import SpecSelect from "../wow/SpecSelect.svelte";

  import { TeamId } from "./team_api";
  import type { Member } from "./team_api";
  import { createEventDispatcher } from "svelte";
  import { TeamStore } from "./team_store";

  export let memberId: number | undefined = undefined;

  const dispatch = createEventDispatcher();

  $: row = memberId ? $TeamStore.row(memberId) : undefined;
  $: member = row ? row.member : undefined;

  let localId,
    localName = "",
    localClassId = 0,
    localSpecs: number[] = [],
    localPrimarySpec = 0;

  function localMember(): Member {
    return {
      id: localId,
      team: $TeamId,
      name: localName,
      classId: localClassId,
      config: {
        specs: localSpecs,
        primarySpec: localPrimarySpec,
      },
    };
  }

  let specsClassId;

  $: if (memberId != undefined) {
    if (localId != $member.id) {
      localId = $member.id;
      localName = $member.name;
      localClassId = $member.classId;
      localSpecs.length = 0;
      $member.config.specs.forEach((s) => localSpecs.push(s));
      specsClassId = $member.classId;
      localPrimarySpec = $member.config.primarySpec;
    }
  }

  $: if (memberId == undefined) {
    localId = undefined;
    localClassId = 0;
    localName = "";
    localClassId = 0;
    localSpecs.length = 0;
    localPrimarySpec = 0;
  }

  $: if (specsClassId != localClassId) {
    localSpecs.length = 0;
    localPrimarySpec = 0;
  }

  function remove() {
    if (member) {
      member.remove();
    }
    close();
  }

  function save() {
    const newMember = localMember();

    if ($member) {
      $member = newMember;
    } else {
      $TeamStore.newMember(newMember);
    }
    close();
  }

  function close() {
    dispatch("close");
    localId = undefined;
    localClassId = 0;
    localName = "";
    localClassId = 0;
    localSpecs.length = 0;
    localPrimarySpec = 0;
  }

  $: saveEnabled =
    localName != "" &&
    localClassId > 0 &&
    localSpecs.length > 0 &&
    localPrimarySpec > 0;
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
      {#if member}
        <Button on:click={remove}>
          {#if $$slots.removeLabel}
            <slot name="removeLabel" />
          {:else}
            <Icon class="material-icons" style="color: red">delete</Icon>
            <Label>Remove</Label>
          {/if}
        </Button>
      {/if}
      <Button on:click={save} disabled={!saveEnabled}>
        {#if $$slots.saveLabel}
          <slot name="saveLabel" />
        {:else}
          <Icon class="material-icons" style="color: green">save</Icon>
          <Label>Save</Label>
        {/if}
      </Button>
      <Button on:click={close}>
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
  :global(.member-edit) {
    :global(.smui-select--standard.mdc-select--with-leading-icon
        .mdc-select__anchor) {
      padding-left: 0;
    }

    :global(.class-select) {
      padding-top: 4px;
      padding-bottom: 4px;
    }
  }
</style>
