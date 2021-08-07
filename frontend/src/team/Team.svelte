<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table/styled";
  import Dialog from "@smui/dialog/styled";
  import LinearProgress from "@smui/linear-progress/styled";
  import Button, { Icon, Label } from "@smui/button/styled";
  import MemberEdit from "./MemberEdit.svelte";
  import RosterEntry from "./RosterEntry.svelte";

  import { CurrentTeam, SortMembers } from "./api";
  import type { Member } from "./api";
  import { writable } from "svelte/store";
  import { Spec } from "../wow/api";

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

  $: $CurrentTeam.then((t) => {
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
</script>

<DataTable class="team">
  <Head>
    <Row>
      <Cell />
      <Cell />
    </Row>
  </Head>

  <Body>
    {#each $team.tanks as member (member.id)}
      <RosterEntry
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    {#each $team.healers as member (member.id)}
      <RosterEntry
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    {#each $team.dps as member (member.id)}
      <RosterEntry
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}

    <Row class="roster-new">
      <Cell>
        <Button on:click={() => showEdit(undefined)}>
          <Icon class="material-icons">add</Icon>
          <Label>Add Member</Label>
        </Button>
      </Cell>
      <Cell />
    </Row>
  </Body>
  <LinearProgress bind:closed={teamLoaded} indeterminate slot="progress" />
</DataTable>
<Dialog bind:open class="add-member">
  <MemberEdit on:close={close} memberId={editMemberId} />
</Dialog>

<style lang="scss">
  :global(.roster-new) {
    $line-height: 32px;

    height: $line-height;

    :global(td.mdc-data-table__cell) {
      height: $line-height;
    }
  }

  :global(.add-member) {
    overflow-y: visible;

    :global(.mdc-dialog__surface) {
      overflow-y: visible;
    }
  }
</style>
