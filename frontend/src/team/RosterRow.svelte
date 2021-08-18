<script lang="ts">
  import Textfield, { Underline } from "smelte/src/components/TextField";
  import Icon from "smelte/src/components/Icon";
  import WowIcon from "../wow/WowIcon.svelte";
  import EncounterCells from "./EncounterCells.svelte";

  import { createEventDispatcher } from "svelte";
  import { TeamStore } from "./team_store";

  export let memberId: number;

  $: row = $TeamStore.row(memberId);
  $: member = row.member;

  $: row.member.init();

  let loadedId;
  let classId = 0;
  let name = "";

  $: if ($member && loadedId != memberId) {
    loadedId = memberId;
    classId = $member.classId;
    name = $member.name;
  }

  function save() {
    row.member.update((m) => {
      console.info("save", { m, name });
      if (m) {
        return { ...m, name };
      }
      return m;
    });
  }

  function keypress(e) {
    if (e.keyCode === 13) save();
  }

  let hovered = false;
  $: configHoveredClass = hovered ? "configure--active" : "";
  $: nameHoveredClass = hovered ? "name--hovered" : "";

  const dispatch = createEventDispatcher();

  const trClass =
    "roster-row hover:bg-gray-50 dark-hover:bg-dark-400 border-gray-200 dark:border-gray-400 border-t border-b px-3";
  const tdClass = "relative px-3 font-normal text-right";

  let nameFocused = false;
  function toggleNameFocused() {
    nameFocused = !nameFocused;
  }
  $: nameFocusedClass = nameFocused
    ? "bg-gray-300 dark:bg-dark-400"
    : "bg-gray-100 dark:bg-dark-600";
</script>

<tr
  class={trClass}
  on:mouseenter={() => (hovered = true)}
  on:mouseleave={() => (hovered = false)}
>
  <td class={tdClass}>
    <div class="roster-member inline-flex h-full align-bottom">
      <WowIcon playerClass={classId} class="p-1" />

      <div class="relative text-gray-600 dark:text-gray-100">
        <input
          class={`name rounded-t text-black dark:text-gray-100 caret-primary w-full ${nameFocusedClass} ${nameHoveredClass}`}
          bind:value={name}
          on:blur={save}
          on:keyup={keypress}
          on:focus={toggleNameFocused}
          on:blur={toggleNameFocused}
          placeholder="Name"
        />
        <Underline focused={nameFocused} />
      </div>

      <Icon
        on:click={() => dispatch("configure")}
        class={"configure h-full cursor-pointer " + configHoveredClass}
      >
        settings
      </Icon>
    </div>
  </td>

  <EncounterCells {memberId} />
</tr>

<style lang="scss" global>
  .roster-row {
    height: 36px;

    .roster-member {
      $config-size: 34px;

      .name {
        $width: 12em;

        width: $width;
        // height: auto;

        &.name--hovered {
          width: calc(#{$width} - #{$config-size});
        }
      }
      .configure {
        display: none;
        width: $config-size;
        height: 100%;

        &.configure--active {
          display: unset;
        }
      }
    }
  }
</style>
