<script lang="ts">
  import { Row, Cell } from "@smui/data-table/styled";
  import Button, { Icon } from "@smui/button/styled";
  import Textfield from "@smui/textfield/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { createEventDispatcher } from "svelte";
  import { TeamMember } from "./api";

  export let memberId: number;

  $: member = TeamMember(memberId);

  let loadedId;
  let classId = 0;
  let name = "";

  $: $member.then((m) => {
    if (loadedId != memberId) {
      loadedId = memberId;
      classId = m.classId;
      name = m.name;
    }
  });

  function save() {
    member.update((m) => {
      m.name = name;
      return m;
    });
  }

  const dispatch = createEventDispatcher();
</script>

<Row class="roster-entry">
  <Cell class="roster-member">
    <div style="height: 100%; align-items: center; display: inline-flex">
      <WowIcon playerClass={classId} height={24} style="padding-right: 8px;" />
    </div>

    <Textfield
      bind:value={name}
      style="width: 10em"
      input$placeholder="Name"
      on:blur={save}
    />

    <Button
      on:click={() => dispatch("configure")}
      color="white"
      style="min-width: 32px"
    >
      <Icon class="material-icons" style="margin-right: 0">settings</Icon>
    </Button>
  </Cell>
  <Cell />
</Row>

<style lang="scss">
  :global(.roster-entry) {
    $line-height: 32px;

    height: $line-height;

    :global(td.mdc-data-table__cell) {
      height: $line-height;
    }

    :global(.roster-member) {
      display: flex;

      :global(label.smui-text-field--standard) {
        height: $line-height - 4px;
      }
    }
  }
</style>
