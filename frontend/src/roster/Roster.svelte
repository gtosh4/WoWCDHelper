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

  interface member {
    id: number;
    name: string;
    classIcon: string;
    className: string;
  }

  let membersLoaded = false;
  let members: member[] = [];
  fetch("/team/test")
    .then((r) => r.json())
    .then((r) => {
      members = r;
      membersLoaded = true;
    })
    .catch((e) => console.error("error loading members", e));

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
      exec: (m: member) => {
        members = [...members.filter((e) => e.id != m.id)];
      },
    },
    { name: "Clone", icon: "content_copy", exec: () => {} },
  ];

  let hovered: number | null = null;
</script>

<DataTable class="roster-table">
  <Head>
    <Row>
      <Cell columnId="name">
        <Label>Name</Label>
      </Cell>
      <Cell columnId="class" width="42px">
        <Label>Class</Label>
      </Cell>
      <Cell columnId="actions" />
    </Row>
  </Head>

  <Body>
    {#each members as member (member.id)}
      <Row
        on:mouseenter={() => (hovered = member.id)}
        on:mouseleave={() => (hovered = null)}
      >
        <Cell>
          <Textfield variant="outlined" bind:value={member.name} />
        </Cell>
        <Cell>
          <img
            src={member.classIcon}
            alt={member.className}
            class="class-icon"
          />
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
      <Cell colspan="2">
        <Select
          key={(c) => (c && c.id) || ""}
          bind:value={createClass}
          width="10em"
        >
          {#each classes as cls (cls.id)}
            <Option value={cls.name}>{cls.name}</Option>
          {/each}
        </Select>
      </Cell>
    </Row>
  </Body>
</DataTable>

<style lang="scss">
  .class-icon {
    height: 32px;
  }

  :global(.roster-table) {
    & tr {
      height: 42px;
    }
  }

  :global(.roster-action) {
    display: none;
    padding-left: auto;
  }
  :global(.roster-action.hovered) {
    display: inherit;
  }
</style>
