<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table/styled";
  import Dialog from "@smui/dialog/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import Lazy from "svelte-lazy";
  import MemberEdit from "./MemberEdit.svelte";
  import RosterRow from "./RosterRow.svelte";
  import EncounterHeaders from "./EncounterHeaders.svelte";

  import RoleRow from "./RoleRow.svelte";
  import { TeamStore } from "./team_store";
  import { AllRoleMembers, SortMembers } from "./member_filter";
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

  let team: Promise<{
    tanks: Member[];
    healers: Member[];
    dps: Member[];
  }> = new Promise(() => {});

  $: if ($members) {
    team = AllRoleMembers($members).then((roleMembers) => {
      teamLoaded = true;

      return {
        tanks: (roleMembers.get("TANK") || []).sort(SortMembers),
        healers: (roleMembers.get("HEALER") || []).sort(SortMembers),
        dps: (roleMembers.get("DAMAGE") || []).sort(SortMembers),
      };
    });
  }

  function removeEncounter(id: number) {
    $TeamStore.column(id).encounter.remove();
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

        <Cell />
        {#if $encounters}
          {#each $encounters as enc, i (i)}
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
        {/if}

        <Cell />
      </Row>
    {/await}
  </Body>
  <LinearProgress bind:closed={teamLoaded} indeterminate slot="progress" />
</DataTable>
<Dialog bind:open class="add-member" on:MDCDialog:closed={close}>
  <Lazy>
    <MemberEdit on:close={close} memberId={editMemberId} />
  </Lazy>
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
