<script lang="ts">
  export let playerClass: string | number = "";
  export let spec: string | number = "";
  export let spellName: string | number = "";

  let src: string;
  let alt: string;
  $: if (spellName) {
    src = `/wow/spell/${spellName}/icon`;
    alt = spellName.toString();
  } else if (spec && typeof spec == "number") {
    src = `/wow/spec/${spec}/icon`;
    alt = spec.toString();
  } else if (playerClass && spec) {
    src = `/wow/class/${playerClass}/${spec}/icon`;
    alt = spec.toString();
  } else if (playerClass) {
    src = `/wow/class/${playerClass}/icon`;
    alt = playerClass.toString();
  }
</script>

<img {src} {alt} {...$$restProps} />
