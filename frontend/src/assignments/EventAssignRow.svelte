<script lang="ts">
  import type { Event, EventInstance } from "../team/team_api";

  export let event: Event;
  export let instance: EventInstance;

  const rclass =
    "hover:bg-gray-50 dark-hover:bg-dark-400 border-gray-200 dark:border-gray-400 border-t border-b px-3";
  const cclass = "relative p-3 font-normal text-right";
</script>

<tr class={rclass}>
  <td class={cclass}>{instance.offset_sec}</td>
  <td class={cclass}>{event.label}</td>
  <td class={cclass + " w-full"} />
</tr>
