<script lang="ts">
  import DataTable, {
    Head,
    Body,
    Row,
    Cell,
    Label,
  } from "@smui/data-table/styled";
  import Textfield from "@smui/textfield/styled";
  import Select, { Option } from "@smui/select/styled";
  import SegmentedButton, {
    Segment,
    Icon,
  } from "@smui/segmented-button/styled";

  import type { Member } from "./api";
  import { CurrentTeam } from "./api";
  import WowIcon from "../wow/WowIcon.svelte";

  let createName: string = "";
  let createClass: string = "";

  let classes = [];
  fetch("/wow/class-info")
    .then((r) => r.json())
    .then((r) => (classes = r))
    .catch((e) => console.error("error loading classes", e));

  const actions = [
    {
      name: "Delete",
      icon: "delete",
      exec: (m: Member) => {
        CurrentTeam.deleteMember(m.id);
      },
    },
    { name: "Clone", icon: "content_copy", exec: () => {} },
  ];

  let hovered: number | null = null;

  let addMember = {
    name: "Add",
    icon: "add",
    color: "green",
    enabled: false,
    exec: () => {
      CurrentTeam.addMember({
        id: 0,
        name: createName,
        className: createClass,
      });
      createName = "";
      createClass = "";
    },
  };
  let clearAdd = {
    name: "Clear",
    icon: "clear",
    color: "red",
    enabled: false,
    exec: () => {
      createName = "";
      createClass = "";
    },
  };
  $: if (createName && createClass) {
    addMember = { ...addMember, enabled: true };
  }
  $: if (createName || createClass) {
    clearAdd = { ...clearAdd, enabled: true };
  }
</script>

<DataTable class="roster-table">
  <Head>
    <Row>
      <Cell columnId="name">
        <Label>Name</Label>
      </Cell>
      <Cell columnId="class">
        <Label>Class</Label>
      </Cell>
      <Cell columnId="actions" />
    </Row>
  </Head>

  <Body>
    {#each $CurrentTeam as member (member.id)}
      <Row
        on:mouseenter={() => (hovered = member.id)}
        on:mouseleave={() => (hovered = null)}
      >
        <Cell>
          <Textfield variant="outlined" bind:value={member.name} />
        </Cell>
        <Cell>
          <WowIcon className={member.className} />
        </Cell>
        <Cell>
          <SegmentedButton
            segments={actions}
            let:segment
            singleSelect
            key={(segment) => segment.name}
            class={hovered == member.id
              ? "roster-action hovered"
              : "roster-action"}
          >
            <Segment
              {segment}
              title={segment.name}
              on:click$preventDefault={segment.exec(member)}
            >
              <Icon class="material-icons" style="width: 1em; height: auto;">
                {segment.icon}
              </Icon>
            </Segment>
          </SegmentedButton>
        </Cell>
      </Row>
    {/each}

    <Row class="add-member">
      <Cell>
        <Textfield variant="outlined" bind:value={createName} />
      </Cell>
      <Cell>
        <Select key={(c) => (c && c.id) || ""} bind:value={createClass}>
          <Option value="" />
          {#each classes as cls (cls.id)}
            <Option value={cls.name}>{cls.name}</Option>
          {/each}
        </Select>
      </Cell>
      <Cell>
        <SegmentedButton
          segments={[addMember, clearAdd]}
          let:segment
          singleSelect
          key={(segment) => segment.name}
          class="class-select"
        >
          <Segment
            {segment}
            title={segment.name}
            disabled={!segment.enabled}
            color={segment.color}
            on:click$preventDefault={segment.exec()}
          >
            <Icon class="material-icons" style="width: 1em; height: auto;">
              {segment.icon}
            </Icon>
          </Segment>
        </SegmentedButton>
      </Cell>
    </Row>
  </Body>
</DataTable>
<Select key={(c) => (c && c.id) || ""} bind:value={createClass} width="10em">
  {#each classes as cls (cls.id)}
    <Option value={cls.name}>{cls.name}</Option>
  {/each}
</Select>

<style lang="scss">
  .class-icon {
    height: 32px;
  }

  :global(button.mdc-segmented-button__segment:disabled) {
    cursor: default;
    pointer-events: none;
    background-color: transparent;
    color: rgba(255, 255, 255, 0.38);
  }

  :global(.mdc-select__menu.class-select) {
    position: absolute;
  }

  :global(.roster-table tr) {
    height: 42px;
  }

  :global(.roster-table .mdc-text-field) {
    height: 100%;
  }

  :global(.roster-action) {
    display: none;
    padding-left: auto;
  }
  :global(.roster-action.hovered) {
    display: inherit;
  }
</style>
