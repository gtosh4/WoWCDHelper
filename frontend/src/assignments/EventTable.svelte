<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table/styled";
  import Dialog from "@smui/dialog/styled";
  import { TeamStore } from "../team/team_store";

  export let encounterId: number;

  $: column = $TeamStore.column(encounterId);
  $: events = column.events;

  $: events.init();

  $: items = ($events || [])
    .flatMap((evt) => evt.instances.map((ei) => ({ event: evt, instance: ei })))
    .sort((a, b) => {
      if (a.instance.offset_sec != b.instance.offset_sec) {
        return a.instance.offset_sec - b.instance.offset_sec;
      } else {
        return a.event.id - b.event.id;
      }
    });

  $: console.info("event table", { events: $events, items: items });

  let open;
  function close() {}
</script>

<DataTable class="event-table">
  <Body>
    {#each items as item, i (i)}
      <Row>
        <Cell>{item.instance.offset_sec}</Cell>
        <Cell>{item.event.label}</Cell>
        <Cell />
      </Row>
    {/each}
  </Body>
</DataTable>
<Dialog bind:open class="add-member" on:MDCDialog:closed={close} />

<style lang="scss" global>
  .event-table {
    height: 100%;
  }
</style>
