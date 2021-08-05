<script lang="ts">
  import { Spec } from "./api";

  import WowIcon from "./WowIcon.svelte";

  export let specId = 0;

  const spec = Spec(specId);
</script>

<div {...$$restProps}>
  <WowIcon spec={specId} class="spec-icon" height={24} />
  <span>{$spec.name}</span>
</div>
