<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table/styled";
  import Dialog from "@smui/dialog/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import MemberEdit from "./MemberEdit.svelte";
  import RosterRow from "./RosterRow.svelte";
  import EncounterHeaders from "./EncounterHeaders.svelte";

  import RoleRow from "./RoleRow.svelte";
  import { TeamStore } from "./team_store";
  import { RoleMembers, SortMembers } from "./member_filter";
  import type { Member } from "./team_api";

  let open = false;
  let editMemberId: number | undefined;

  function showEdit(id?: number) {
    open = true;
    editMemberId = id;
  }

  function close() {
    open = false;
    editMemberId = undefined;
  }

  let teamLoaded = false;

  $: members = $TeamStore.Members;
  $: encounters = $TeamStore.Encounters;

  $: team = Promise.all([
    RoleMembers($members, "TANK").then((ms) => ({
      tanks: ms.sort(SortMembers),
    })),
    RoleMembers($members, "HEALER").then((ms) => ({
      healers: ms.sort(SortMembers),
    })),
    RoleMembers($members, "DAMAGE").then((ms) => ({
      dps: ms.sort(SortMembers),
    })),
  ]).then((roleMembers) => {
    teamLoaded = true;

    return roleMembers.reduce((team, role) => ({ ...team, ...role }), {}) as {
      tanks: Member[];
      healers: Member[];
      dps: Member[];
    };
  });

  function removeEncounter(id: number) {
    $TeamStore.column(id).encounterInfo.remove();
  }
</script>

<DataTable class="roster">
  <Head>
    <Row>
      <Cell />
      <EncounterHeaders />
    </Row>
  </Head>

  <Body>
    {#await team then team}
      <RoleRow roleName="Tanks" roleType="TANK" />
      {#each team.tanks as member (member.id)}
        <RosterRow
          memberId={member.id}
          on:configure={() => showEdit(member.id)}
        />
      {/each}
      <RoleRow roleName="Healers" roleType="HEALER" />
      {#each team.healers as member (member.id)}
        <RosterRow
          memberId={member.id}
          on:configure={() => showEdit(member.id)}
        />
      {/each}
      <RoleRow roleName="DPS" roleType="DAMAGE" />
      {#each team.dps as member (member.id)}
        <RosterRow
          memberId={member.id}
          on:configure={() => showEdit(member.id)}
        />
      {/each}
      <Row>
        <Cell class="roster-new">
          <Button on:click={() => showEdit(undefined)}>
            <Icon class="material-icons">add</Icon>
            <Label>Add Member</Label>
          </Button>
        </Cell>

        {#await $encounters}
          <Cell />
        {:then encounters}
          <Cell />
          {#each encounters as enc, i (i)}
            <Cell>
              <Button
                style="min-width: 32px"
                on:click={() => removeEncounter(enc.id)}
              >
                <Icon
                  class="material-icons"
                  style="margin-right: 0; color: red"
                >
                  remove_circle
                </Icon>
              </Button>
            </Cell>
          {/each}
        {/await}

        <Cell />
      </Row>
    {/await}
  </Body>
  <LinearProgress bind:closed={teamLoaded} indeterminate slot="progress" />
</DataTable>
<Dialog bind:open class="add-member" on:MDCDialog:closed={close}>
  <MemberEdit on:close={close} memberId={editMemberId} />
</Dialog>

<style lang="scss" global>
  .roster {
    .mdc-data-table__table-container {
      overflow: visible;
    }

    .roster-new {
      $line-height: 32px;

      height: $line-height;

      td.mdc-data-table__cell {
        height: $line-height;
      }
    }
  }
  .add-member {
    overflow-y: visible;

    .mdc-dialog__surface {
      overflow-y: visible;
    }
  }
</style>
