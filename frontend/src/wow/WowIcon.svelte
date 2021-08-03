<script lang="ts">
  export let className: string = "";
  export let specName: string = "";
  export let spellName: string | number = "";

  let src: string;
  let alt: string;
  if (spellName) {
    src = `/wow/spell/${spellName}/icon`;
    alt = spellName.toString();
  } else if (className && specName) {
    src = `/wow/class/${className}/${specName}/icon`;
    alt = specName;
  } else if (className) {
    src = `/wow/class/${className}/icon`;
    alt = className;
  }
</script>

<img {src} {alt} {...$$restProps} />
