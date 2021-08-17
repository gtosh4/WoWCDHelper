<script lang="ts">
  import Card from "smelte/src/components/Card";
  import Textfield from "smelte/src/components/TextField";
  import Button from "smelte/src/components/Button";
  import Icon from "smelte/src/components/Icon";
  import ProgressLinear from "smelte/src/components/ProgressLinear";
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

<Card.Card padded class="member-edit">
  {#await $member}
    <ProgressLinear />
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

  <div slot="actions">
    {#if member}
      <Button on:click={remove}>
        <Icon class="material-icons" style="color: red">delete</Icon>
        <span>Remove</span>
      </Button>
    {/if}
    <Button on:click={save} disabled={!saveEnabled}>
      <Icon class="material-icons" style="color: green">save</Icon>
      <span>Save</span>
    </Button>
    <Button on:click={close}>
      <span>Cancel</span>
    </Button>
  </div>
</Card.Card>

<style lang="scss" global>
  // :global(.member-edit) {
  //   :global(.smui-select--standard.mdc-select--with-leading-icon
  //       .mdc-select__anchor) {
  //     padding-left: 0;
  //   }

  //   :global(.class-select) {
  //     padding-top: 4px;
  //     padding-bottom: 4px;
  //   }
  // }
</style>
