<script lang="ts">
  import type { Event, EventInstance } from "../team/team_api";
  import { TeamStore } from "../team/team_store";
  import EventAssignRow from "./EventAssignRow.svelte";

  export let encounterId: number;

  $: column = $TeamStore.column(encounterId);
  $: events = column.events;

  $: events.init();

  interface item {
    event: Event;
    instance: EventInstance;
  }

  $: items = ($events || [])
    .flatMap((evt): item[] =>
      evt.instances.map((ei) => ({ event: evt, instance: ei }))
    )
    .sort((a, b) => {
      if (a.instance.offset_sec != b.instance.offset_sec) {
        return a.instance.offset_sec - b.instance.offset_sec;
      } else {
        return a.event.id - b.event.id;
      }
    });

  $: console.info("event table", { events: $events, items: items });

  const tclass =
    "shadow relative text-sm coverflow-x-auto dark:bg-dark-500 w-full";
</script>

<table class={tclass}>
  {#each items as item, i (i)}
    <EventAssignRow event={item.event} instance={item.instance} />
  {/each}
</table>

<style lang="scss" global>
</style>
