<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table/styled";
  import Dialog from "@smui/dialog/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import MemberEdit from "./MemberEdit.svelte";
  import RosterRow from "./RosterRow.svelte";
  import EncounterHeaders from "./EncounterHeaders.svelte";

  import { Members, SortMembers } from "./members_api";
  import type { Member } from "./team_api";
  import { writable } from "svelte/store";
  import { Spec } from "../wow/api";
  import { Encounters } from "./encounters_api";
  import RoleRow from "./RoleRow.svelte";

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
  let team = writable({
    tanks: [] as Member[],
    healers: [] as Member[],
    dps: [] as Member[],
  });

  $: $Members.then((t) => {
    const newTeam = {
      tanks: [] as Member[],
      healers: [] as Member[],
      dps: [] as Member[],
    };
    teamLoaded = false;
    Promise.all(
      [...t.values()].map((member) =>
        Spec(member.config.primarySpec).then((spec) => {
          switch (spec.role.type) {
            case "TANK":
              newTeam.tanks.push(member);
              break;

            case "HEALER":
              newTeam.healers.push(member);
              break;

            default:
              newTeam.dps.push(member);
          }
        })
      )
    ).then(() => {
      newTeam.tanks.sort(SortMembers);
      newTeam.healers.sort(SortMembers);
      newTeam.dps.sort(SortMembers);

      $team = newTeam;
      teamLoaded = true;
    });
  });

  function removeEncounter(id: number) {
    Encounters.encounter(id).remove();
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
    <RoleRow roleName="Tanks" roleType="TANK" />
    {#each $team.tanks as member (member.id)}
      <RosterRow
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    <RoleRow roleName="Healers" roleType="HEALER" />
    {#each $team.healers as member (member.id)}
      <RosterRow
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    <RoleRow roleName="DPS" roleType="DAMAGE" />
    {#each $team.dps as member (member.id)}
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

      {#await $Encounters}
        <Cell />
      {:then encounters}
        <Cell />
        {#each encounters as enc, i (i)}
          <Cell>
            <Button
              style="min-width: 32px"
              on:click={() => removeEncounter(enc.id)}
            >
              <Icon class="material-icons" style="margin-right: 0; color: red">
                remove_circle
              </Icon>
            </Button>
          </Cell>
        {/each}
      {/await}

      <Cell />
    </Row>
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
