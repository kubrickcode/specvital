export default {
  "**/*.{json,yml,yaml,md}": (files) => files.map((f) => `just lint-file "${f}"`),
  "**/*.go": (files) => files.map((f) => `just lint-file "${f}"`),
  "**/[Jj]ustfile": () => "just --fmt --unstable",
};
