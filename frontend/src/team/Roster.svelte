<script lang="ts">
  import Dialog from "smelte/src/components/Dialog";
  import { slide } from "svelte/transition";
  import ProgressLinear from "smelte/src/components/ProgressLinear";
  import Button from "smelte/src/components/Button";
  import Icon from "smelte/src/components/Icon";
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

  let team: {
    tanks: Member[];
    healers: Member[];
    dps: Member[];
  } = { tanks: [], healers: [], dps: [] };

  $: if ($members) {
    AllRoleMembers($members).then((roleMembers) => {
      teamLoaded = true;

      team = {
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

<table class="roster shadow text-sm overflow-x-auto dark:bg-dark-500">
  <thead class="items-center">
    <tr>
      <th />
      <EncounterHeaders />
    </tr>
  </thead>
  {#if !teamLoaded}
    <div class="absolute w-full" transition:slide>
      <ProgressLinear />
    </div>
  {/if}

  <tbody>
    <RoleRow roleName="Tanks" roleType="TANK" />
    {#each team.tanks || [] as member (member.id)}
      <RosterRow
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    <RoleRow roleName="Healers" roleType="HEALER" />
    {#each team.healers || [] as member (member.id)}
      <RosterRow
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    <RoleRow roleName="DPS" roleType="DAMAGE" />
    {#each team.dps || [] as member (member.id)}
      <RosterRow
        memberId={member.id}
        on:configure={() => showEdit(member.id)}
      />
    {/each}
    <tr class="border-gray-200 dark:border-gray-400 border-t border-b px-3">
      <td class="roster-new">
        <Button
          on:click={() => showEdit(undefined)}
          class="flex items-center"
          text
          block
          outlined
        >
          <Icon>add</Icon>
          Add Member
        </Button>
      </td>

      <td />
      {#each $encounters || [] as enc, i (i)}
        <td>
          <Icon
            class="cursor-pointer"
            style="margin-right: 0; color: red"
            on:click={() => removeEncounter(enc.id)}
          />
        </td>
      {/each}

      <td />
    </tr>
  </tbody>
</table>
<Dialog bind:value={open} class="overflow-y-visible">
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
</style>
