<script lang="ts">
  import DataTable from "smelte/src/components/DataTable";
  import type { Event, EventInstance } from "../team/team_api";
  import { TeamStore } from "../team/team_store";

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

  const columns = [
    { value: (i: item) => `${i.instance.offset_sec}` },
    { value: (i: item) => `${i.event.label}` },
  ];
</script>

<DataTable
  class="event-table"
  data={items}
  {columns}
  pagination={false}
  sortable={false}
/>

<style lang="scss" global>
  .event-table {
    height: 100%;
  }
</style>
