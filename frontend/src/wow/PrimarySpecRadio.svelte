<script lang="ts">
  import Icon from "smelte/src/components/Icon";
  import Ripple from "smelte/src/components/Ripple";

  export let primarySpec: number;
  export let specs: number[];
  export let thisSpec: number;

  $: disabled = specs.indexOf(thisSpec) < 0;

  function select() {
    if (disabled) return;
    primarySpec = thisSpec;
  }

  $: rippleColor = !disabled ? "primary" : "gray";
  $: selected = primarySpec == thisSpec;
</script>

{#if !disabled}
  <div
    class="inline-flex block items-center mb-2 cursor-pointer z-0"
    on:click={select}
  >
    <input class="hidden" type="radio" role="radio" {selected} />
    <div class="relative">
      <Ripple color={rippleColor} noHover>
        {#if selected}
          <Icon class="text-primary-500">radio_button_checked</Icon>
        {:else}
          <Icon class="text-gray-600">radio_button_unchecked</Icon>
        {/if}
      </Ripple>
    </div>
  </div>
{/if}
