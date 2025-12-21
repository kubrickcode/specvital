module.exports = {
  "**/*.{json,yml,yaml,md}": () => "just lint config",
  "**/*.{js,jsx,ts,tsx}": () => "just lint frontend",
  "**/*.go": () => "just lint backend",
  justfile: () => "just lint justfile",
};
