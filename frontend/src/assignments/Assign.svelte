<script lang="ts">
  import { TeamStore } from "../team/team_store";
  import { ClassColors } from "../wow/class_colors";

  let classes = "";

  export let memberId: number;
  export let encounterId: number;
  export let assignId: number | undefined = undefined;
  export { classes as class };

  $: cell = $TeamStore.cell(memberId, encounterId);

  let style = "";
  $: {
    const m = {
      "background-color": ClassColors[$cell.member.classId].fade(0.5),
    };
    m["color"] = m["background-color"].isDark() ? "white" : "black";

    style = Object.entries(m)
      .map(([k, v]) => `${k}: ${v};`)
      .join(" ");
  }

  let hovered = false;

  function cardClass(): string {
    let s = "rounded inline-flex flex-col overflow-hidden duration-200 ease-in"; // card
    s += " shadow-sm hover:shadow"; // card.hover
    s += " flex justify-start";
    if (hovered) {
      s += " hovered";
    }
    s += ` ${classes}`;
    return s;
  }
</script>

<div
  class={cardClass()}
  {style}
  on:mouseenter={() => (hovered = true)}
  on:mouseleave={() => (hovered = false)}
>
  <i class="material-icons">drag_indicator</i>
  <span style="flex: flex-grow">{$cell.member.name}</span>
  {#if assignId}
    <i class="material-icons close-icon">close</i>
  {/if}
</div>

<style lang="scss" global>
  // .assign-member-item {
  //   display: flex;
  //   flex-direction: row;
  //   flex-wrap: nowrap;
  //   align-items: center;
  //   padding-top: 1px;
  //   padding-bottom: 1px;

  //   .close-icon {
  //     display: none;
  //   }

  //   &.hovered {
  //     .close-icon {
  //       display: unset;
  //     }
  //   }
  // }
</style>
