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
  import { LoadingState } from "../store_helpers";

  export let memberId: number | undefined = undefined;

  const dispatch = createEventDispatcher();

  $: row = memberId ? $TeamStore.row(memberId) : undefined;
  $: member = row ? row.member : undefined;

  $: console.log("edit", { row, member, m: $member, lm: localMember() });

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

<Card.Card class="overflow-y-visible">
  <div slot="title">
    <Textfield class="w-full" bind:value={localName} label="Name" />
  </div>

  {#if member.state != LoadingState.Loaded}
    <ProgressLinear />
  {/if}

  <ClassSelect bind:value={localClassId} />

  <SpecSelect
    classId={localClassId}
    bind:specs={localSpecs}
    bind:primarySpec={localPrimarySpec}
  />

  <div slot="actions">
    {#if member}
      <Button small on:click={remove}>
        <Icon small color="red">delete</Icon>
        Remove
      </Button>
    {/if}
    <Button small on:click={save} disabled={!saveEnabled}>
      <Icon small color="green">save</Icon>
      Save
    </Button>
    <Button small on:click={close}>
      <Icon small>undo</Icon>
      Cancel
    </Button>
  </div>
</Card.Card>

<style lang="scss" global>
</style>
