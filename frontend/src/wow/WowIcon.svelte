<script lang="ts">
  export let playerClass: string | number = "";
  export let spec: string | number = "";
  export let spellName: string | number = "";
  export let height: number | undefined = 26;
  export { extraClass as class };

  let extraClass: string = "";

  let src: string;
  let alt: string;
  let typeClass: string;
  $: if (spellName) {
    src = `/wow/spell/${spellName}/icon`;
    alt = spellName.toString();
    typeClass = "spell-icon";
  } else if (spec && typeof spec == "number") {
    src = `/wow/spec/${spec}/icon`;
    alt = spec.toString();
    typeClass = "spec-icon";
  } else if (playerClass && spec) {
    src = `/wow/class/${playerClass}/${spec}/icon`;
    alt = spec.toString();
    typeClass = "spec-icon";
  } else if (playerClass) {
    src = `/wow/class/${playerClass}/icon`;
    alt = playerClass.toString();
    typeClass = "class-icon";
  }
</script>

<img
  {src}
  {alt}
  class={`${typeClass} ${extraClass}`}
  {height}
  width={height}
  on:click
  {...$$restProps}
/>
