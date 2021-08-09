<script lang="ts">
  import { Row, Cell } from "@smui/data-table/styled";
  import Button, { Icon } from "@smui/button/styled";
  import Textfield from "@smui/textfield/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { createEventDispatcher } from "svelte";
  import { Members } from "./members_api";
  import EncounterRow from "./EncounterRow.svelte";

  export let memberId: number;

  $: member = Members.member(memberId);

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
      if (m.name != name) {
        m.name = name;
        return m;
      }
      return undefined;
    });
  }

  function keypress(e) {
    if (e.keyCode === 13) save();
  }

  let hovered = false;

  const dispatch = createEventDispatcher();
</script>

<Row
  class="roster-entry"
  on:mouseenter={() => (hovered = true)}
  on:mouseleave={() => (hovered = false)}
>
  <Cell class="roster-member">
    <div style="height: 100%; align-items: center; display: inline-flex">
      <WowIcon playerClass={classId} height={24} style="padding-right: 8px;" />
    </div>

    <Textfield
      bind:value={name}
      class={hovered ? "name name--hovered" : "name"}
      input$placeholder="Name"
      on:blur={save}
      on:keyup={keypress}
    />

    <Button
      on:click={() => dispatch("configure")}
      class={hovered ? "configure configure--active" : "configure"}
    >
      <Icon class="material-icons" style="margin-right: 0; color: white">
        settings
      </Icon>
    </Button>
  </Cell>

  <EncounterRow {memberId} />
</Row>

<style lang="scss" global>
  .roster-entry {
    $line-height: auto;

    height: $line-height;

    td.mdc-data-table__cell {
      height: $line-height;
    }

    .roster-member {
      display: flex;

      .name {
        $width: 12em;

        width: $width;
        height: auto;

        &.name--hovered {
          width: calc(#{$width} - 18px - 8px - 8px);
        }
      }
      .configure {
        display: none;
        min-width: unset;
        width: calc(18px+8px+8px);
        height: $line-height;

        &.configure--active {
          display: unset;
        }
      }
    }
  }
</style>
