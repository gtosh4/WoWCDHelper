import Color from "color";

function c(id: number): Color {
  return Color(
    getComputedStyle(document.documentElement).getPropertyValue(
      `--wow-class-${id}`
    )
  );
}

export const ClassColors = {
  1: c(1),
  2: c(2),
  3: c(3),
  4: c(4),
  5: c(5),
  6: c(6),
  7: c(7),
  8: c(8),
  9: c(9),
  10: c(10),
  11: c(11),
  12: c(12),
};
