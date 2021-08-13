<script lang="ts">
  import { Row, Cell } from "@smui/data-table/styled";
  import Button, { Icon } from "@smui/button/styled";
  import Textfield from "@smui/textfield/styled";
  import WowIcon from "../wow/WowIcon.svelte";

  import { createEventDispatcher } from "svelte";
  import EncounterCells from "./EncounterCells.svelte";
  import { LoadingState, TeamStore } from "./team_store";

  export let memberId: number;

  $: row = $TeamStore.row(memberId);
  $: member = row.member;

  $: if (row.member.state == LoadingState.Uninitialized) {
    row.member.reload();
  }

  let loadedId;
  let classId = 0;
  let name = "";

  $: if ($member && loadedId != memberId) {
    loadedId = memberId;
    classId = $member.classId;
    name = $member.name;
  }

  function save() {
    if ($member) {
      row.member.put({
        ...$member,
        name,
      });
    }
  }

  function keypress(e) {
    if (e.keyCode === 13) save();
  }

  let hovered = false;

  const dispatch = createEventDispatcher();
</script>

<Row
  class="roster-row"
  on:mouseenter={() => (hovered = true)}
  on:mouseleave={() => (hovered = false)}
>
  <Cell style="vertical-align: bottom;">
    <div class="roster-member">
      <div style="height: 100%; align-items: center; display: inline-flex">
        <WowIcon
          playerClass={classId}
          height={24}
          style="padding-right: 8px;"
        />
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
    </div>
  </Cell>

  <EncounterCells {memberId} />
</Row>

<style lang="scss" global>
  .roster-row {
    height: auto;

    td.mdc-data-table__cell {
      height: auto;
    }

    .roster-member {
      display: inline-flex;
      height: 100%;

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
        height: 100%;

        &.configure--active {
          display: unset;
        }
      }
    }
  }
</style>
