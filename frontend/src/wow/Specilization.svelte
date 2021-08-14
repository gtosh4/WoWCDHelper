<script lang="ts">
  import { Spec } from "./api";

  import WowIcon from "./WowIcon.svelte";

  export let specId = 0;

  const spec = Spec(specId);
</script>

<div {...$$restProps} class="wow-spec">
  <WowIcon spec={specId} class="spec-icon" height={24} />
  <span class="spec-name">
    {#await spec}
      ...
    {:then s}
      {s.name}
    {/await}
  </span>
</div>

<style lang="scss" global>
  .wow-spec {
    display: flex;

    .spec-name {
      padding-left: 4px;
    }
  }
</style>
