# Changelog

## [0.0.1](https://github.com/kubrickcode/specvital/compare/v0.0.0...v0.0.1) (2026-02-13)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **ci:** fix CI dependency installation for monorepo workspace ([e0c67d1](https://github.com/kubrickcode/specvital/commit/e0c67d1c27ebdf682626bf0aa7f2dc6dbf95dbc1))
- **ci:** fix eslint command for pnpm workspace context ([15192cc](https://github.com/kubrickcode/specvital/commit/15192cc886e43122b82589169d82d4999caeba4c))
- **ci:** fix frontend independent install and lint in pnpm workspace ([16f21ea](https://github.com/kubrickcode/specvital/commit/16f21eaceb88e6fe7901985244e9e1c712a19bdb))
- **ci:** move reusable workflows from subdirectory to flat structure ([cc1a547](https://github.com/kubrickcode/specvital/commit/cc1a5471701fbfb144cc09dd405fd7f91a9ebe77))
- fix install-playwright failing to find package from root ([67c379a](https://github.com/kubrickcode/specvital/commit/67c379a012e617c10e49c781de913cfce7569866))
- fixed the issue where generated-types.ts would change when running the just lint command. ([24c21c2](https://github.com/kubrickcode/specvital/commit/24c21c2c155773ef9be67dbd9c46ee3f20b46fe3))
- pin pnpm and node versions in frontend deploy workflow ([47a2cfa](https://github.com/kubrickcode/specvital/commit/47a2cfa34c97954f902b63636cc1f280f03cd384))

#### 📚 Documentation

- consolidate README, LICENSE, NOTICE into monorepo root ([934c564](https://github.com/kubrickcode/specvital/commit/934c5645f08044af11dec0327ddac1442669230c))
- rewrite CLAUDE.md files to reflect monorepo structure ([93ae8ad](https://github.com/kubrickcode/specvital/commit/93ae8ad70e1fa4d23b61234baa5a1b60b9093e23))

#### 💄 Styles

- format code ([d4a6edb](https://github.com/kubrickcode/specvital/commit/d4a6edb5219d3a61b9a83e45aca42f111b72ca88))

#### ♻️ Refactoring

- add spec-generator mock runner and remove unused smee ([5e21bc4](https://github.com/kubrickcode/specvital/commit/5e21bc4389f8e43ce86800d64e8a58d62824a178))
- centralize Railway deploy configs into root infra/railway/ ([4be3fe9](https://github.com/kubrickcode/specvital/commit/4be3fe9776e3033fc2d874b3e044f80d088aa0c6))
- consolidate .vscode settings and simplify justfile for single DB ([77e9903](https://github.com/kubrickcode/specvital/commit/77e99038d7bce5241dda8e37fe012b4518bd1179))
- consolidate bootstrap/install recipes into root justfile ([cadd6a9](https://github.com/kubrickcode/specvital/commit/cadd6a96d3a68740ad92caa9b1f4ac3ac5ccc584))
- consolidate dev tool configs from sub-projects into root ([f698527](https://github.com/kubrickcode/specvital/commit/f698527defbc959af8f83af950da45c02c92a3c6))
- consolidate global justfile logic from sub-projects into root ([6c7a9b0](https://github.com/kubrickcode/specvital/commit/6c7a9b0b0a6a0fe52b6c98d95246adebdf1723b8))
- consolidate migrate-local into single root migrate ([415c585](https://github.com/kubrickcode/specvital/commit/415c5859196fadad7606f86f2b274471c1646cd6))
- consolidate sub-project CI workflows into root ([5067082](https://github.com/kubrickcode/specvital/commit/50670824413ef078f09e3be4b91b8d687fd03d4a))
- consolidate sub-project release workflows into single orchestrator ([506f265](https://github.com/kubrickcode/specvital/commit/506f265ca35c81acb075615d80db98c478646f48))
- deduplicate docs and flatten docs/docs nesting ([3c22afd](https://github.com/kubrickcode/specvital/commit/3c22afda8ae8348f4257d48ca29a2009d00fff2d))
- rationalize monorepo package structure ([8640b73](https://github.com/kubrickcode/specvital/commit/8640b732b1587abf857d7d4daaef196973b41488))
- set up Go workspace and migrate module paths to monorepo structure ([446cb42](https://github.com/kubrickcode/specvital/commit/446cb42703a41476f9bbca87d266a19ee6366fb4))
- unify lint/lint-file rules into per-language helper recipes ([e087e3a](https://github.com/kubrickcode/specvital/commit/e087e3aeb0b5a8d1432fa5877ea474884be1f924))

#### ✅ Tests

- slim E2E suite to Tier 1 and migrate removed tests to Vitest components ([60d5c95](https://github.com/kubrickcode/specvital/commit/60d5c950a93d31379a72c289cfa458cd441178da))

#### 🔧 CI/CD

- add frontend Vitest and worker unit test workflows ([1f8918a](https://github.com/kubrickcode/specvital/commit/1f8918a7d0fa27f2b3a9a757b40ac30ef54c120d))
- add Vercel CLI frontend deploy workflow and rename deploy-web to deploy-backend ([00d7316](https://github.com/kubrickcode/specvital/commit/00d731665185afa2f46eed5e4356cdae217fc5da))
- unify workers deploy workflow to Railway CLI v4 syntax ([7754390](https://github.com/kubrickcode/specvital/commit/77543906d5814ea4719de84ac321055cce6a7147))

#### 🔨 Chore

- add install-atlas to bootstrap ([a1ad5f4](https://github.com/kubrickcode/specvital/commit/a1ad5f42d14b0b1cb3b8a9da11a5432dcc2bfdc0))
- deps all ([b649c44](https://github.com/kubrickcode/specvital/commit/b649c4499deed9660f098fae8296bb4a0c33af18))
- dump schema ([157a9f2](https://github.com/kubrickcode/specvital/commit/157a9f2230088b0ca9dcc569f9660d3909a2d5eb))
- remove individual CHANGELOG ([bdbafba](https://github.com/kubrickcode/specvital/commit/bdbafba7570d97a01a8e9a58746fca86b481c96e))
- remove obsolete update-core commands and update docs for monorepo ([7047502](https://github.com/kubrickcode/specvital/commit/70475026f8c1c4293590482b50baa77f21e67060))
- remove unnecessary command ([ef1bad4](https://github.com/kubrickcode/specvital/commit/ef1bad41ccc1d21bf3acc7d715f2900e66a25552))
- remove unnecessary file ([9bdfb5c](https://github.com/kubrickcode/specvital/commit/9bdfb5cf44491d0af720d2e9e319e4532d060e2f))
- set up pnpm workspace and consolidate duplicate JS tooling ([d594c16](https://github.com/kubrickcode/specvital/commit/d594c16372807ef424afd4d50154685096350d52))

> Unified changelog for all SpecVital packages. For pre-monorepo history, see individual package changelogs.

## [web/v1.6.3](https://github.com/specvital/web/compare/v1.6.2...v1.6.3) (2026-02-09)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **infra:** set Railway deploy region to US East ([db548d2](https://github.com/specvital/web/commit/db548d204eb6f32bc7f316c019e6dc8bf38bcdf6))

## [worker/v1.3.1](https://github.com/specvital/worker/compare/v1.3.0...v1.3.1) (2026-02-09)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **infra:** set Railway deploy region to US East ([b0163e2](https://github.com/specvital/worker/commit/b0163e24668ad43bf49e222125e78aa56cbecd7f))

## [web/v1.6.2](https://github.com/specvital/web/compare/v1.6.1...v1.6.2) (2026-02-05)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **analysis:** resolve polling delay and elapsed time display issues ([15542fb](https://github.com/specvital/web/commit/15542fb7e4c0e14f2339c1893388ced9fab12a69))
- **auth:** overhaul React Query cache policy to resolve stale data bugs ([dfa6aae](https://github.com/specvital/web/commit/dfa6aaeaecae6fab4dc9ac9f86fae7d3248403fe))
- **dashboard:** summary cards not updating on task completion ([5eab7ee](https://github.com/specvital/web/commit/5eab7eef6de5fcfc001d61b06e938c0fefafc244))

### 🔧 Maintenance

#### ♻️ Refactoring

- **background-tasks:** replace sessionStorage-based task-store with server API ([7456bdf](https://github.com/specvital/web/commit/7456bdf1da20aeef0788fd6c945a443a1b8c7dc0))

#### 🔨 Chore

- claude code execution command modified to always run in a new terminal ([8c97ed6](https://github.com/specvital/web/commit/8c97ed60ab7643c060d53f732c3f02bec8c89d09))
- sync ai & container config from kubrickcode/ai-config-toolkit ([38bd132](https://github.com/specvital/web/commit/38bd132bc6634cabeb32be8a8f60c76796891f92))

## [worker/v1.3.0](https://github.com/specvital/worker/compare/v1.2.4...v1.3.0) (2026-02-05)

### 🎯 Highlights

#### ✨ Features

- **analysis:** add batch storage infrastructure for streaming pipeline ([5cf4008](https://github.com/specvital/worker/commit/5cf40082de45d66d301d32a8991fbd33a85acc54))
- **analysis:** add logging for streaming analysis processing ([72ad88e](https://github.com/specvital/worker/commit/72ad88eaf16c7a783d4cfa2b79eec3b5a492177a))
- **analysis:** integrate streaming parsing mode ([bce6b44](https://github.com/specvital/worker/commit/bce6b44a351cdfaf988b89b91c79a7f9ffb3e8cc))
- **parser:** implement streaming parser for memory-efficient large repo analysis ([d25fdb0](https://github.com/specvital/worker/commit/d25fdb0801705114b8c0d935a338adb2c2ce247d))

#### 🐛 Bug Fixes

- **fairness:** per-user concurrency limit was applied globally ([9f2a693](https://github.com/specvital/worker/commit/9f2a693d84e33de5bce034fb8d58c682b7097885))

#### ⚡ Performance

- **analysis:** optimize timeout and connection settings for large repositories ([05338c3](https://github.com/specvital/worker/commit/05338c3c315b468bb59c7a2b61730541366bf613))
- **analysis:** reduce N roundtrips to 1 by batching TestFile inserts ([4a3a3bb](https://github.com/specvital/worker/commit/4a3a3bb236ff43d8b2841adc402b85588ba7da68))

### 🔧 Maintenance

#### ♻️ Refactoring

- **specview:** remove Phase 1 V2/V3 architecture dead code ([406c778](https://github.com/specvital/worker/commit/406c7781c3a387d4c6e6b9f6203bb70da49d6cf6))

#### 🔨 Chore

- claude code execution command modified to always run in a new terminal ([beec3c0](https://github.com/specvital/worker/commit/beec3c0da0cdc689070fffd4dafc438595928a26))
- sync ai & container config from kubrickcode/ai-config-toolkit ([06bd1ab](https://github.com/specvital/worker/commit/06bd1ab30c17d327846a82fe770a9e51f13fd9b1))
- update-core ([1c2210b](https://github.com/specvital/worker/commit/1c2210b0ad66a7cd4dfc995d586021d9f0efade8))

## [web/v1.6.1](https://github.com/specvital/web/compare/v1.6.0...v1.6.1) (2026-02-04)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **i18n:** apply missing i18n to toast messages ([9be4e2f](https://github.com/specvital/web/commit/9be4e2f88e203a3e80c25adbf9bc0b694e2f6590))
- **ui:** fix loading spinner appearing below visual center ([516cf00](https://github.com/specvital/web/commit/516cf005070cf3ca92144d1d04b2923d49e965e7))

## [core/v1.6.0](https://github.com/specvital/core/compare/v1.5.0...v1.6.0) (2026-02-04)

### 🎯 Highlights

#### ✨ Features

- **domain-hints:** add DomainHints quality validation command ([cd95f41](https://github.com/specvital/core/commit/cd95f419722e1d5b6d773c96b765a9b3e2296ad0))
- **domain-hints:** add noise filter for string literal methods, function args, and URL patterns ([def841a](https://github.com/specvital/core/commit/def841a868f4c437a24b3694d91608f11224c2ca))
- **domain-hints:** add noise filtering for domain hints extraction ([3349196](https://github.com/specvital/core/commit/3349196122187db45393d52cc4354e558c1bf34d))
- **domain-hints:** apply noise filter to JavaScript/TypeScript extractor ([dce6467](https://github.com/specvital/core/commit/dce64679ee9ea4d7271145bb7597e1211bfe8e1e))
- **domain-hints:** apply noise filter to remaining 7 language extractors ([deacdda](https://github.com/specvital/core/commit/deacdda09e17635bf0bff20de06d38ef941d8800))
- **domain-hints:** apply noise filter to Rust/Kotlin extractor ([a8646e2](https://github.com/specvital/core/commit/a8646e295e0de253c70b8e79c50f4276f26a1c26))
- **domain-hints:** implement DomainHints extraction for C# test files ([80590f2](https://github.com/specvital/core/commit/80590f29c6aa21ce025ea8f352afacc222a47a2c))
- **domain-hints:** implement DomainHints extraction for Java/Kotlin test files ([4dfb02a](https://github.com/specvital/core/commit/4dfb02ae9a2e47ae8e6aed89b8026ff6e764c020))
- **domain-hints:** implement DomainHints extraction for Python test files ([7007f0d](https://github.com/specvital/core/commit/7007f0d6062bac2a9ac9a05fdbe9a4be784bb2fd))
- **domain-hints:** implement DomainHints extraction for Ruby/PHP test files ([2374cbf](https://github.com/specvital/core/commit/2374cbf0b068a1e990d7843245c29b63b5c1e624))
- **domain-hints:** implement DomainHints extraction for Rust/Swift/C++ test files ([ae87dbf](https://github.com/specvital/core/commit/ae87dbf52158153cf2b585ea073571da16e9926b))
- **domain-hints:** implement Go test file domain hints extraction ([fd3cb93](https://github.com/specvital/core/commit/fd3cb930cc7980a6b577bf953d4b2bdcccd608c8))
- **domain-hints:** implement JS/TS test file domain hints extraction ([00e034a](https://github.com/specvital/core/commit/00e034a812509fccbbec244796daec68bf8a302e))
- **domain:** add DomainHints type for AI-based domain classification ([7efd5a5](https://github.com/specvital/core/commit/7efd5a51b7046aef93a1901799dc8be18401664d))
- **parser:** add FileResult type for streaming API ([753049b](https://github.com/specvital/core/commit/753049bec6ed965d98eee09aa756897623ba3066))
- **parser:** add ScanStreaming public API with documentation ([f8011b3](https://github.com/specvital/core/commit/f8011b3b985a724c3d550c964cd914c239e09866))
- **parser:** implement ScanStream method for streaming parsing ([a45e1ad](https://github.com/specvital/core/commit/a45e1ad9a29726acdb5683f4bbfd93928686fa70))

#### 🐛 Bug Fixes

- **detection:** support vitest globals mode test file detection ([1049be5](https://github.com/specvital/core/commit/1049be5e01ed68bc88aa669a0148ed41ed0b0887))
- **domain-hints:** add import noise and Java Object method filtering ([38ba020](https://github.com/specvital/core/commit/38ba020335d29982dcf649a9fde0b61e6a1be6f3))
- **domain-hints:** filter all single-character calls as noise ([53ee28f](https://github.com/specvital/core/commit/53ee28fd5f39cbab8f814068e32685eea028b55e))
- **domain-hints:** filter C# stdlib imports and nameof as noise ([72df4fa](https://github.com/specvital/core/commit/72df4fa938ba6d7d1e40fbaa753aedf63936e82a))
- **domain-hints:** filter C++ dot-prefix noise patterns and stdlib imports ([db8338d](https://github.com/specvital/core/commit/db8338d9297e3518db419923d1f44b52d7c05a3b))
- **domain-hints:** filter Go builtin functions as noise ([223c2c9](https://github.com/specvital/core/commit/223c2c95f99fb77def4e5c7c9b4d8617e9d52a11))
- **domain-hints:** filter Go stdlib imports as noise ([74aa9ce](https://github.com/specvital/core/commit/74aa9ceb525ccee6c126dc723ee2167ca1a474f3))
- **domain-hints:** filter inline comments and short (≤2 char) standalone calls as noise ([6cebcb2](https://github.com/specvital/core/commit/6cebcb25858a5fc5056b188cd2f60ee2ac6aaff5))
- **domain-hints:** filter Kotlin stdlib collection functions as noise ([9d664bc](https://github.com/specvital/core/commit/9d664bc0fd7c932c4bba62ca7dbb023ed6601084))
- **domain-hints:** filter multiline generic imports and Object base methods in C# ([1d5399f](https://github.com/specvital/core/commit/1d5399faca7aee9ef2e473af9f8174978429e9d9))
- **domain-hints:** filter out generic callback variable name fn as noise ([d2a13f0](https://github.com/specvital/core/commit/d2a13f0adcd501db0bf3cd340d0604c2605519ce))
- **domain-hints:** filter unbalanced parentheses patterns as noise ([33ebdd8](https://github.com/specvital/core/commit/33ebdd86328f66c0fc421261c971db251c566a8e))
- **domain-hints:** fix Rust multiline use statement parsing bug ([a97e31e](https://github.com/specvital/core/commit/a97e31e0f33c529fa9d3e05acaffa531ae307d0a))
- **domain-hints:** improve DomainHints extraction data quality ([e4b9f2d](https://github.com/specvital/core/commit/e4b9f2d82a58c0b24bce090910cfa68e3a10f56d))
- **parser:** detect dynamic subtests in Go test parser ([0fca289](https://github.com/specvital/core/commit/0fca289bc58037ebce1ee6ac435ae0ecbd649513))
- **pytest:** support unittest.TestCase style class detection ([d126e03](https://github.com/specvital/core/commit/d126e03d8219eaafa0fb8d96b4828f9fdf15e7ed))
- **scanner:** add test/, tests/ directory patterns to JS test file discovery ([4de8e36](https://github.com/specvital/core/commit/4de8e36d0439f957c6ca49f699c2fc07b78bfe77))

### 🔧 Maintenance

#### 📚 Documentation

- add specvital-specialist agent ([527833f](https://github.com/specvital/core/commit/527833fc4fd2b11d1f57b7e8f79585c2ca41922e))
- **domain-hints:** document DomainHints option and data structure ([c726d23](https://github.com/specvital/core/commit/c726d23051d5c7cf3efbe907aaf4df4bbcb8bc99))

#### ♻️ Refactoring

- **domain-hints:** extract noise filter to common module ([7fec069](https://github.com/specvital/core/commit/7fec069074be2f443585f5633b021c576c0dd2c0))
- **domain-hints:** remove Variables field from DomainHints ([8b83b95](https://github.com/specvital/core/commit/8b83b95520e0c828a63a72de262fcde7df622395))
- **parser:** refactor file discovery to channel-based streaming ([dc6536f](https://github.com/specvital/core/commit/dc6536fef465d0a7cf5adceb577fa8c4acda0d12))
- **parser:** unify Scan API implementation to streaming-based ([bd5f311](https://github.com/specvital/core/commit/bd5f311db014b91078f707e154db1537bff3067c))

#### ✅ Tests

- **domain-hints:** add integration tests for DomainHints extraction using grafana repository ([21ad5f1](https://github.com/specvital/core/commit/21ad5f18c53194182e19983c2691816ee6e4fb50))
- **domain-hints:** add tRPC v10.45.2 noise pattern validation ([898e6a8](https://github.com/specvital/core/commit/898e6a8dbb81ef3039838d48fe0d68172a7bb11f))
- **integration:** add chi and httpx repositories to integration tests ([b9d12af](https://github.com/specvital/core/commit/b9d12aff3144845152f229bf1f9b0d654a6edeb0))
- **integration:** add chi repository to integration tests ([c98eaf6](https://github.com/specvital/core/commit/c98eaf6a29aa9ca6981cb34aec9877f629ad338e))
- **integration:** add echo repository to integration tests ([4180d3d](https://github.com/specvital/core/commit/4180d3dd3e0a3fa986c81b9a5a3550b77b686ef5))
- **integration:** add gson repository ([32f5e2e](https://github.com/specvital/core/commit/32f5e2e133747b64153a09414f25efb13c54d4ca))
- **integration:** add jackson-core repository ([bb3357e](https://github.com/specvital/core/commit/bb3357ebcbbff28dd0efc1348063713efbf9c145))
- **integration:** add pydantic repository ([93bf2db](https://github.com/specvital/core/commit/93bf2dbd9187496f801d6f858a3f2d6f62a66434))
- **integration:** add zod repository to integration tests ([423e6e4](https://github.com/specvital/core/commit/423e6e48b712f8520cf1089a4900428d38ba4ebc))

#### 🔨 Chore

- change license from MIT to Apache 2.0 ([2743e6d](https://github.com/specvital/core/commit/2743e6db12adc6419b573c318ea8b69fd2a16ecb))
- remove commit file ([535243c](https://github.com/specvital/core/commit/535243c753a2500b7cc90a00d1212b4c2f6a6aa6))
- snapshot update ([5e006a8](https://github.com/specvital/core/commit/5e006a83e8788a0e80d3ebe841aa6453109ceeb8))
- snapshot-update ([e415f16](https://github.com/specvital/core/commit/e415f16bd7189ed1f9b1f65120d1ca34c98bf169))
- sync ai-config-toolkit ([cc69e74](https://github.com/specvital/core/commit/cc69e745c353c623798fca0c4bd33ac66cea54f5))
- sync-docs ([2587358](https://github.com/specvital/core/commit/2587358377a5d6b907bd7b93210a472f391551ba))

## [web/v1.6.0](https://github.com/specvital/web/compare/v1.5.0...v1.6.0) (2026-02-03)

### 🎯 Highlights

#### ✨ Features

- **analysis:** add analysis history API endpoint ([2dc0326](https://github.com/specvital/web/commit/2dc032673a1c84fdabaef639d85e843a271c45bd))
- **analysis:** add commit history selection dropdown UI ([bb215fa](https://github.com/specvital/web/commit/bb215fafe1b9cc394ac8018712ac7d06329a1350))
- **analysis:** add commit switching with URL state management ([b9170c9](https://github.com/specvital/web/commit/b9170c924548a04ae3d785d64470cec06d369de4))
- **analysis:** add commit-specific analysis query API ([da00d30](https://github.com/specvital/web/commit/da00d304f103a06d7e7cdcb160bf327514166f77))
- **spec-view:** add dynamic cost estimation based on cache prediction ([3645d52](https://github.com/specvital/web/commit/3645d52ac69724c3d7cf329e75f803adfbb245b9))

#### 🐛 Bug Fixes

- **analysis:** fix immediate completion toast on reanalysis ([354d49a](https://github.com/specvital/web/commit/354d49a209daaf117af6e6099c7d6a0e3ab8a2f3))
- **analysis:** fix new commits banner not showing on analysis detail page ([27090f9](https://github.com/specvital/web/commit/27090f9edadf305ae7d4e8f0c22e683d684cb071))
- **analysis:** fix toast appearing immediately on commit update ([3417432](https://github.com/specvital/web/commit/34174324dd8acdb5554b62f563545133055e5391))
- **dashboard:** fix AI Spec badge disappearing when new commits exist ([79df048](https://github.com/specvital/web/commit/79df048eae03e3df81a517b4129421f97f7bb138))
- **spec-view:** fix version switching not working when selecting same version from different commits ([052e6d0](https://github.com/specvital/web/commit/052e6d06d4ff2d7a5cc1b7c336537ae3ebe09d6e))

### 🔧 Maintenance

#### ♻️ Refactoring

- **i18n:** move messages resources into i18n module ([d6c0f67](https://github.com/specvital/web/commit/d6c0f6756b2a8fe3cecf710ab81bf6026a2b2096))

## [web/v1.5.0](https://github.com/specvital/web/compare/v1.4.1...v1.5.0) (2026-02-02)

### 🎯 Highlights

#### ✨ Features

- **spec-view:** group version history by commit SHA ([42b852b](https://github.com/specvital/web/commit/42b852b0662972ccd61459cc601e2d778aa0538d))

#### 🐛 Bug Fixes

- **spec-view:** fix multiple documents showing as "latest" in version history ([b217af4](https://github.com/specvital/web/commit/b217af4486ccf86224c0bcacc6293dd3d241f095))

## [web/v1.4.1](https://github.com/specvital/web/compare/v1.4.0...v1.4.1) (2026-02-02)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix release failure due to semantic-release plugin version mismatch ([0426f2d](https://github.com/specvital/web/commit/0426f2dee52237c2a8063ebd9d32379722cc8b39))

## [web/v1.4.0](https://github.com/specvital/web/compare/v1.3.1...v1.4.0) (2026-02-02)

### 🎯 Highlights

#### ✨ Features

- **account:** add subscription plan and usage status page ([db64e74](https://github.com/specvital/web/commit/db64e74ea05d23a85c3a8d7cbde445cf796ef17e))
- **account:** add upgrade button to plan section ([8f242fa](https://github.com/specvital/web/commit/8f242fa0fc5ae3a2e2406894a9d35f7c85b0f2bb))
- add close button to toast messages ([2d5552b](https://github.com/specvital/web/commit/2d5552b477279f5a087936914e2291f816f27055))
- **analysis-header:** convert metadata to collapsible panel ([6e5ff2a](https://github.com/specvital/web/commit/6e5ff2a68015387b44898d8248b0830174dbe35f))
- **analysis-header:** improve GitHub link icon clarity ([9450bda](https://github.com/specvital/web/commit/9450bda74c6d60d683561e3388aa0610045cd1f1))
- **analysis-header:** unify View on GitHub button style to outline ([fd53d55](https://github.com/specvital/web/commit/fd53d55de475f68844a390f510e2b09c95cbbe56))
- **analysis:** add progress information to analysis waiting state ([fe631d8](https://github.com/specvital/web/commit/fe631d87fd6c9a93d22971cd8313dd1f4575aa37))
- **analysis:** add visual emphasis for AI Spec tab ([cc29695](https://github.com/specvital/web/commit/cc29695bff6973c40803b9c11b26e30d2d8a97e7))
- **analysis:** redesign tab-based UI information hierarchy ([c87c8fe](https://github.com/specvital/web/commit/c87c8fea225261e822aa0bbb91321429c34a6ec9))
- **analysis:** remove auto-reanalysis on page access and add update banner ([7c360c5](https://github.com/specvital/web/commit/7c360c5bb13a430e0070dfeb591d5d358404318d))
- **analysis:** separate AI spec button and view mode toggle roles ([2889dc8](https://github.com/specvital/web/commit/2889dc8125d9077c156ae7e8598fa917fcbc78fe))
- **analysis:** separate initial loading state from analysis progress UI ([4d659f3](https://github.com/specvital/web/commit/4d659f31f242c3516f0fd439c073ec0710ec59c3))
- **analyzer:** add parser_version comparison for conditional re-analysis ([0260122](https://github.com/specvital/web/commit/02601224a390c36fce7177c685ca66f90c9b3a75))
- **analyzer:** add parser_version query from system_config ([ad04a81](https://github.com/specvital/web/commit/ad04a817e232f113cac4c49210b1f159e4cc085d))
- **analyzer:** add parserVersion field to API response ([f8427db](https://github.com/specvital/web/commit/f8427dbe5b0d6bd61f43bd6792fafaf569536591))
- **analyzer:** add rate limiting for anonymous users on analyze API ([107f387](https://github.com/specvital/web/commit/107f387c0d8fa953d95d7143d93abd3f7709c11f))
- **analyzer:** add startedAt field to AnalyzingResponse ([a6219e7](https://github.com/specvital/web/commit/a6219e7acd5da62825fe6a6e7f51730a89fea4d4))
- **analyzer:** lookup user plan tier in handler for queue selection ([1572ee8](https://github.com/specvital/web/commit/1572ee8a65056e154818ad72630d0b33026b0b97))
- **background-tasks:** add Dashboard active tasks section and extract shared components ([fc99ce5](https://github.com/specvital/web/commit/fc99ce53e0c9b1e76950d3765097d26a9e9455ea))
- **background-tasks:** add global task store with persistence ([8664cbc](https://github.com/specvital/web/commit/8664cbcc2eceb88b33cac613add978d37127b427))
- **background-tasks:** enhance task badge visibility and improve loading feedback consistency ([985975f](https://github.com/specvital/web/commit/985975f0f1b9fd7401548e60496376f31392aae2))
- **background-tasks:** integrate Account Badge and Tasks Dropdown ([ecb4434](https://github.com/specvital/web/commit/ecb4434fdf2d9b85bd5ab19cfd6bdbaf3efeba41))
- **dashboard:** add AiSpecBadge component for repository cards ([46cc3e2](https://github.com/specvital/web/commit/46cc3e282cfc6c29c9c02ab17d0fba239c1d5b20))
- **dashboard:** add AiSpecSummary schema for repository card badge ([1127f49](https://github.com/specvital/web/commit/1127f49a7359954799aa03507afe509f90befe83))
- **dashboard:** integrate AI Spec summary into repository card API response ([382aa47](https://github.com/specvital/web/commit/382aa47e642b23fecdc64c254474e5b42e83c756))
- **dashboard:** integrate AiSpecBadge into RepositoryCard component ([20a0196](https://github.com/specvital/web/commit/20a0196a6d97ee37d90f8da703b657971084648b))
- **docs:** add docs landing page with how-it-works TOC ([265ddbd](https://github.com/specvital/web/commit/265ddbd214a273227bb420719499840b794c1e0d))
- **docs:** add github-access concept page ([6254a2c](https://github.com/specvital/web/commit/6254a2c5dfb162f054146e261fae9fdbdc08b982))
- **docs:** add queue-processing concept page ([486c417](https://github.com/specvital/web/commit/486c417d1b8bc0dd9234d23c37b5509e88afeeec))
- **docs:** add sidebar navigation infrastructure for how-it-works ([aae09b4](https://github.com/specvital/web/commit/aae09b4df46c15476bd5e6cc51ebb9d423bd5607))
- **docs:** add specview-generation concept page ([4a627dc](https://github.com/specvital/web/commit/4a627dcb447f5acc5e9c188a2bdceeb75f3ae0c7))
- **docs:** add test writing guide documentation page ([aaeec9f](https://github.com/specvital/web/commit/aaeec9f975787d7e0a9edead3a8c620ebb34f1b0))
- **docs:** add test-detection concept page ([cf889ee](https://github.com/specvital/web/commit/cf889eea958b4aac96f8339d3c7b5c15133519bf))
- **docs:** add usage-billing concept page ([e3e3e98](https://github.com/specvital/web/commit/e3e3e98dc86e4459f9ce4efc75d3c21af3a3dc80))
- **docs:** remove how-it-works category and simplify docs structure ([a82f007](https://github.com/specvital/web/commit/a82f007726c0aa302b4ff158c681b6ff004fd4da))
- **e2e:** set up Playwright-based E2E UI test infrastructure ([5951acf](https://github.com/specvital/web/commit/5951acfe4067bd39473945ba60e3cb30d7bd6a85))
- **feedback:** add shared animation primitives for async waiting states ([e4a181f](https://github.com/specvital/web/commit/e4a181fab1ed15df1cc495e98a6b3de5ab0f948c))
- **frontend:** add Document View MVP for spec-view feature ([b69f275](https://github.com/specvital/web/commit/b69f2759b057baaaf60f64cab93b13faeeca662f))
- **frontend:** display parser version in analysis view ([3b089bb](https://github.com/specvital/web/commit/3b089bb7b25a94e786f7e26437b3c515bac56cbf))
- **frontend:** enable React Compiler and remove manual memoization ([482d080](https://github.com/specvital/web/commit/482d080e4d060662a4fd0baadedd6cf528361ae3))
- **global-search:** add Cmd+K search dialog foundation ([00f19b2](https://github.com/specvital/web/commit/00f19b2695512d48503adf9901dba402dadecc3a))
- **global-search:** add Header search trigger and complete mobile responsive design ([d274e75](https://github.com/specvital/web/commit/d274e7554eabd4d6d0d8caded2bb75e4f5903896))
- **global-search:** add recent items with localStorage persistence ([343776e](https://github.com/specvital/web/commit/343776ea8f8b5261bb143ecc134f29e2e76d49a9))
- **global-search:** add repository search with fuzzy matching ([593daeb](https://github.com/specvital/web/commit/593daeb3ae12c7e43d78cc686cf75b8467cc05fe))
- **global-search:** add static actions and page navigation ([4c73684](https://github.com/specvital/web/commit/4c73684f4776d6190ff21d3b4ae967a80563d7ce))
- **header:** apply gradient style to New Analysis button ([8b32beb](https://github.com/specvital/web/commit/8b32bebf6bd567e168279d35e0e6b017eafe9c8d))
- **header:** reorder navigation tabs to Docs before Pricing ([922fc3d](https://github.com/specvital/web/commit/922fc3dca5c567b93d0199b4d8a616ac2c549ca2))
- **header:** unify header button heights to 32px ([2139e8b](https://github.com/specvital/web/commit/2139e8b14b62c700373b1ce303294f9c12646097))
- **header:** unify header button styles into single variant ([cd777e5](https://github.com/specvital/web/commit/cd777e5a7ce0978a5d6da65011d6081d0346cdc4))
- **header:** unify New Analysis and utility button styles ([9e5e98b](https://github.com/specvital/web/commit/9e5e98b6270bc2aaba4b535055d34b18cb6c4b4d))
- **inline-stats:** conditionally display focused/xfail status ([a298b51](https://github.com/specvital/web/commit/a298b51ded305f34157566d1c91ba6088338f7af))
- **loading:** add missing loading feedback for server fetching states ([6badf51](https://github.com/specvital/web/commit/6badf516cc0ed94c6f9ce55e3074e14010201954))
- **loading:** add skeleton loading for explore/account page entry ([4e84237](https://github.com/specvital/web/commit/4e84237c43e780d550517cddc61150837ffc4bac))
- **mini-bar:** render all 5 test status color segments with improved accessibility ([05d0146](https://github.com/specvital/web/commit/05d014617c64169a0f149de193060fdbefc45767))
- **mobile-nav:** add more menu to mobile bottom bar for Docs and Pricing access ([2f52cd4](https://github.com/specvital/web/commit/2f52cd4b23841196d38a605bea553b60831a32de))
- **navigation:** add Docs navigation entry to header ([2818699](https://github.com/specvital/web/commit/281869977bcf58ca9e101bf96ce145c96a00edb5))
- **pricing:** add DB-driven pricing API ([66dc783](https://github.com/specvital/web/commit/66dc78332c6eff88996288dc23ff24edef7d1ca5))
- **pricing:** add pricing page with plan comparison and early access promotion ([9838ec4](https://github.com/specvital/web/commit/9838ec46012864903b050948c95f8aaedb4f5a2c))
- **pricing:** remove Pro plan highlight styles and badge ([0fd2505](https://github.com/specvital/web/commit/0fd250514c736d06376a272797779d77515eeb20))
- **queue:** add plan tier-based queue selection for priority processing ([586d31e](https://github.com/specvital/web/commit/586d31e847078e8fcc480e63748e39b3711d7e93))
- show dropdown icon inside Markdown export button ([f78ff37](https://github.com/specvital/web/commit/f78ff375b1c7474aa418ad6f6fae0edf2498656c))
- **sidebar:** apply silver gradient to active navigation items ([f6d249b](https://github.com/specvital/web/commit/f6d249be4d55286e57609c85d1f6406983f1bfbd))
- **spec-view:** add 3-step pipeline visualization for spec generation progress ([84affa1](https://github.com/specvital/web/commit/84affa153c418572bbf4e7972cb7bcd3e33cd3ce))
- **spec-view:** add auth and ownership verification to generation status API ([e1e403a](https://github.com/specvital/web/commit/e1e403a9cbd5b419faa12ab39b8c616b9098d005))
- **spec-view:** add auto-expand for TOC items on scroll ([18b5dca](https://github.com/specvital/web/commit/18b5dca54ba15ab61ba9af7460a076e438466877))
- **spec-view:** add availableLanguages and version to SpecDocument API ([15ebcd8](https://github.com/specvital/web/commit/15ebcd83c95152669ff25e770a9c34bc7cd555eb))
- **spec-view:** add behaviorCacheStats field to SpecDocument ([d4c9f84](https://github.com/specvital/web/commit/d4c9f841e4eec0806af98aa798ca449bdec05415))
- **spec-view:** add cache hit rate display UI ([1eee8b7](https://github.com/specvital/web/commit/1eee8b7315b89e21365a27fcc30f0598dbff61a5))
- **spec-view:** add cache reuse selection UI for spec generation ([e5a7a53](https://github.com/specvital/web/commit/e5a7a53b1c17e02531ea37658708cde204330263))
- **spec-view:** add delete logic for forced spec regeneration ([2309625](https://github.com/specvital/web/commit/2309625df942738e7c87c25d3e23cf1285f4761b))
- **spec-view:** add infrastructure for AI Spec API access control ([9d9f5ed](https://github.com/specvital/web/commit/9d9f5ed3900e971b850e722d535bfb283f6e44e5))
- **spec-view:** add language query parameter for document retrieval ([4dc51bf](https://github.com/specvital/web/commit/4dc51bfdfcb8a54fecb49a459e78e817e9067516))
- **spec-view:** add language switch dropdown UI to ExecutiveSummary ([5b76830](https://github.com/specvital/web/commit/5b76830a0c07c0cd67a0728bab215f172c4ba971))
- **spec-view:** add language-specific cache key separation for React Query ([9f058bf](https://github.com/specvital/web/commit/9f058bf32de2ca580bda776af79ab9b0b2778e4a))
- **spec-view:** add Level 2 API infrastructure for spec document ([9eb5092](https://github.com/specvital/web/commit/9eb5092b80240e91809c43e97de6c8470ab90881))
- **spec-view:** add native names to language selection data ([e2a1677](https://github.com/specvital/web/commit/e2a1677fd744b1e873c15e22e965878728034332))
- **spec-view:** add polish and error handling for production readiness ([f3b0145](https://github.com/specvital/web/commit/f3b01452b83f8a7d1d48c1be264559e599ab627e))
- **spec-view:** add quota confirmation dialog before SpecView generation ([613d465](https://github.com/specvital/web/commit/613d46546fbf0b105015184bd577efc1ea4f0da4))
- **spec-view:** add reading progress indicator for document scroll ([5e09a1d](https://github.com/specvital/web/commit/5e09a1d0c50adcfe366ed25ddfc5e27d988f1160))
- **spec-view:** add regeneration UI to ExecutiveSummary component ([720df60](https://github.com/specvital/web/commit/720df606ce66439eb3d6bd09fd7975f193e31b7a))
- **spec-view:** add repository-based AI Spec query API ([d9bf4a6](https://github.com/specvital/web/commit/d9bf4a6578d7417e96f555d0061fb293e30c2f52))
- **spec-view:** add repository-based version history support in frontend ([46dadc1](https://github.com/specvital/web/commit/46dadc1260c4e4f632afcc484f169d18190224c4))
- **spec-view:** add search and filtering for document view ([79a3740](https://github.com/specvital/web/commit/79a3740af84a8df58b1786d48d6bc38dd6c1f1d7))
- **spec-view:** add search match count feedback badge ([83c7e5f](https://github.com/specvital/web/commit/83c7e5fda768bd07944ddabf98dd623021d8f0a1))
- **spec-view:** add spec document markdown export utility ([9017353](https://github.com/specvital/web/commit/901735345e86737b42fdd4a403772b6242f4d8ab))
- **spec-view:** add Spec Export button to ExecutiveSummary ([1f9d120](https://github.com/specvital/web/commit/1f9d120cef770bb662d8c018974d4c96a3d662ec))
- **spec-view:** add status legend component for test statuses ([f99c040](https://github.com/specvital/web/commit/f99c0408a6be1aca80d5f4664727b17ff6697372))
- **spec-view:** add TOC sidebar and navigation features ([0113c49](https://github.com/specvital/web/commit/0113c4927b8ff5174de5baf156a9baf0c6a74320))
- **spec-view:** add version history API ([1c8c24c](https://github.com/specvital/web/commit/1c8c24c3d8305431ed5c9b4e5104d00983ae04bd))
- **spec-view:** add version history dropdown for viewing previous document versions ([723ec73](https://github.com/specvital/web/commit/723ec731b290691fd3b3ecce8fcaf1e73a6c42d5))
- **spec-view:** add visual indicator banner for old version viewing ([7d08bb4](https://github.com/specvital/web/commit/7d08bb44827be8bcdb1820d72f16e56a54689346))
- **spec-view:** apply last used language as default for new analyses ([f3ede0b](https://github.com/specvital/web/commit/f3ede0bf050c6078fd4f2c27fe90800a0bd5e7da))
- **spec-view:** implement error screens for AI Spec access control ([bbab443](https://github.com/specvital/web/commit/bbab4439761ed4eb3896bf893d20102257b316aa))
- **spec-view:** implement permission verification for AI Spec document access ([50a308f](https://github.com/specvital/web/commit/50a308f8a4c11cd5080519996afb6e9922a3f97d))
- **spec-view:** improve AI spec document generation UX ([7a1d0fe](https://github.com/specvital/web/commit/7a1d0feab82ec642954ead0394f73820f9052562))
- **spec-view:** improve card UI and extract FrameworkBadge as shared component ([689f41c](https://github.com/specvital/web/commit/689f41cd069d5e0a76522bf3c90d497d80814783))
- **spec-view:** integrate document mode into analysis page and connect River queue ([c09339e](https://github.com/specvital/web/commit/c09339ed5b0b3c753186d14ff380436f56da9b2c))
- **spec-view:** integrate quota check into SpecView generation request ([87512a6](https://github.com/specvital/web/commit/87512a6d4d3ea8f15e1036fd5d87cba2982ae122))
- **spec-view:** integrate quota reservation with spec generation request ([8211e06](https://github.com/specvital/web/commit/8211e06c34921b538cd1e9116ccc825b78b63d76))
- **spec-view:** lookup user plan tier in handler for queue selection ([cc9b6d2](https://github.com/specvital/web/commit/cc9b6d29d844c054adcef2703d591051797e7763))
- **spec-view:** pass user_id to spec generation queue ([44bc914](https://github.com/specvital/web/commit/44bc914dfb7045b64adde039e9da9e081d003390))
- **spec-view:** remove AI model name Badge from Executive Summary ([42f9d69](https://github.com/specvital/web/commit/42f9d69bb827979bd21359c4ecd9c72fbe3a0271))
- **spec-view:** replace language select with searchable Combobox ([3a6992f](https://github.com/specvital/web/commit/3a6992ffadae134303601012c31ce598ebf23635))
- **spec-view:** replace time-based fake progress bar with status-based spinner UI ([f9bbb3b](https://github.com/specvital/web/commit/f9bbb3b91624763860dc192bbe2aeb905d661d23))
- **spec-view:** restrict anonymous Spec Generation access with login prompt UI ([ddd9411](https://github.com/specvital/web/commit/ddd9411cbda76d72c662299c6aed9a165c72e7b2))
- **spec-view:** separate free viewing and paid generation with Two-Tier language dropdown ([6d7d83b](https://github.com/specvital/web/commit/6d7d83b863b5a43e77e512c0f6229f4d25990111))
- **spec-view:** show banner for spec documents from previous commits ([50232f4](https://github.com/specvital/web/commit/50232f44e44cdedc5c031989a29dd1121056bbab))
- **spec-view:** show original test name via ResponsiveTooltip ([d6b2693](https://github.com/specvital/web/commit/d6b2693de4867a0ffe946e38c47edd4115a23741))
- **spec-view:** simplify domain section header layout and optimize mobile padding ([54e63ee](https://github.com/specvital/web/commit/54e63ee69e6aff3dbf6c3628952e9baac4c3cd26))
- **subscription:** auto-assign subscription plan on user signup ([8e15e80](https://github.com/specvital/web/commit/8e15e80a527f61f1c073cb705be08ec9a7547075))
- **ui:** add maintenance mode feature ([fadcfbd](https://github.com/specvital/web/commit/fadcfbdb56939926dd13cdbf80dd247aacbbbd9f))
- **ui:** add smooth expand/collapse animation to accordion components ([3409211](https://github.com/specvital/web/commit/340921193ff5ae81126f2b5cf9054cd1e9bd7702))
- **ui:** add subtle background color to outline buttons and toggles ([53acb44](https://github.com/specvital/web/commit/53acb440cb6e48dffb83bd47d2d2b4a128878126))
- **ui:** add unanalyzed variant to RepositoryCard ([609961c](https://github.com/specvital/web/commit/609961cd3d1f5aa4da2472ae313b479578c54f07))
- **ui:** consolidate Empty State framework list to shared constants ([b970352](https://github.com/specvital/web/commit/b97035275245d015bca5eb23062ebb79c7050166))
- **ui:** redesign analysis page stats with minimal layout ([cff5b21](https://github.com/specvital/web/commit/cff5b214ee4059fcd97462e78d1cc60399965f1b))
- **ui:** replace quota modal header icon with AI badge ([24eca3b](https://github.com/specvital/web/commit/24eca3b793362edef392e09d8d8d12ff0b33ec1f))
- **ui:** show applied filters and add reset button on empty filter results ([6575f32](https://github.com/specvital/web/commit/6575f32332014c3e330b7641b25940023554f2fc))
- **ui:** unify quota usage UI in spec generation modal ([1eea40e](https://github.com/specvital/web/commit/1eea40e0abc7ac75b6d21cdecdd3be5fe9b049c9))
- **ui:** unify Star icon and terminology to Bookmark ([373fb5e](https://github.com/specvital/web/commit/373fb5ea1c3d468e19cd0de21b8044d1ad4d13e3))
- **usage:** add quota check API endpoints ([579e96d](https://github.com/specvital/web/commit/579e96d032b1a0b8cf8e744ef050b87a17a3e11e))
- **usage:** add quota reservation repository ([a006201](https://github.com/specvital/web/commit/a00620145257e7bb0bc81e7db47122d22b6e3d7b))
- **usage:** add usage aggregation query infrastructure ([eed6a62](https://github.com/specvital/web/commit/eed6a62ad0a0f04ece6dfcca6ba9e0c516c785de))
- **usage:** display reserved quota info on client ([6a6dd51](https://github.com/specvital/web/commit/6a6dd51aa30cd8dbfa8a798ff9005935e682729f))
- **usage:** include reserved amount in quota check ([e498b1e](https://github.com/specvital/web/commit/e498b1e405e53d213d605b27f996de84f04d9923))

#### 🐛 Bug Fixes

- **account:** add missing unit labels to usage display ([2b8c45f](https://github.com/specvital/web/commit/2b8c45f4b196faf3908ff39775ccc02cccc2defb))
- **account:** fix awkward grammar in usage reset date display ([b364f51](https://github.com/specvital/web/commit/b364f512e02436683a41e166cfe190f77a798323))
- **analysis:** eliminate flicker during tab switch and filter URL state changes ([7ac573a](https://github.com/specvital/web/commit/7ac573a6da145c9b35068709984e2a9f40015974))
- **analysis:** fix accordion overlap in List View ([21a7fb8](https://github.com/specvital/web/commit/21a7fb83eabb6d83accbd6c1446b1563fd621d40))
- **analysis:** fix Export button size mismatch and improve mobile layout ([28d8eba](https://github.com/specvital/web/commit/28d8ebab2030ef52cf850ba884f450b9070d5382))
- **analysis:** fix test card header overflow on mobile ([0e54ecd](https://github.com/specvital/web/commit/0e54ecd5a48a7198a9152bec91d4fc3c02d1f5e8))
- **analysis:** fix truncated text not viewable in analysis page ([8678cf1](https://github.com/specvital/web/commit/8678cf116fa61331ec2879d33c928869edd0db86))
- **analysis:** reflect progress status when clicking Update Now in update-banner ([eaa96d8](https://github.com/specvital/web/commit/eaa96d81b56657078ae598ab9df092d8a7e0b18e))
- **analyzer:** use valid UUID format in test mock data ([3c47b7d](https://github.com/specvital/web/commit/3c47b7d9f4ed7f5027012cc20308059228f67de4))
- **api:** fix repository list disappearing when changing sort option ([a232332](https://github.com/specvital/web/commit/a232332fcd096944605db6f395d165fae3a43977))
- **auth:** eliminate flicker during home/dashboard page transitions ([76e9724](https://github.com/specvital/web/commit/76e972428e4b2db7f76c4426d7c98d7901dd3d62))
- **auth:** fix OAuth login not returning to original page ([9a961ed](https://github.com/specvital/web/commit/9a961ed5981223d90ac82a3f8360cd729a56c1c1))
- **background-tasks:** resolve task badge UI overlap in profile icon area ([70f332b](https://github.com/specvital/web/commit/70f332b3cdce160e16ed48465d0d3a89ff7ed6a3))
- **card:** fix long repository names being truncated without full view option ([3fd7433](https://github.com/specvital/web/commit/3fd7433dd262b51b2554132fb27e72a04d17d63b))
- community tab not showing analyses from non-logged-in users ([0279580](https://github.com/specvital/web/commit/02795801467f64670fa237ff58a7d765adf833c8))
- **dashboard:** fix bookmark terminology inconsistency in emptyState ([9131f3b](https://github.com/specvital/web/commit/9131f3b0f3a43861909c52f315bc504d67dd42af))
- **dashboard:** fix incomplete last row in card grid layout ([88105a4](https://github.com/specvital/web/commit/88105a41cb657d01791b7ba214c8a2c42947f771))
- **docs:** fix asymmetric left/right spacing on mobile docs page ([72b5bab](https://github.com/specvital/web/commit/72b5bab73041eea849b29bbbf3dd392c7fac2555))
- **docs:** fix card content overflow on mobile writing guide page ([1931cf4](https://github.com/specvital/web/commit/1931cf4aa45eba6f049a39679c6ab405ae312aee))
- **e2e:** fix 34 failing E2E tests ([524f1a5](https://github.com/specvital/web/commit/524f1a546c49683f739f81360e58a842db402acd))
- **e2e:** update selectors to match analysis page UI redesign ([bbab230](https://github.com/specvital/web/commit/bbab230b21681a309659847a97029233e0414936))
- **explore:** fix unanalyzed repository card layout inconsistency ([8551743](https://github.com/specvital/web/commit/85517433387a10afd283e4cbd6eb993ca8bc17b5))
- **global-search:** align search trigger button height with adjacent header buttons ([8ae7f30](https://github.com/specvital/web/commit/8ae7f306aabe5720634ea2be549fa9af38cbba89))
- **global-search:** keep command palette open on theme toggle ([2cb19b0](https://github.com/specvital/web/commit/2cb19b0ac7a3dce0007dc4620119820df7db831f))
- **global-search:** resolve hydration mismatch in search trigger button ([d060463](https://github.com/specvital/web/commit/d060463d7bc78986190366806669ae9bcef915f6))
- **header:** fix tab text wrapping vertically at medium screen sizes ([c9f7605](https://github.com/specvital/web/commit/c9f76057ba2ad8710c8a0647fd87ecfc0f28b246))
- **i18n:** add action hint to new commits status message ([5ef030b](https://github.com/specvital/web/commit/5ef030bf367de4bd82cbf4f05425b8b27eab7621))
- **i18n:** add missing translations for spec-view and analysis components ([359b43b](https://github.com/specvital/web/commit/359b43b82f699573a3ad27810bc5b434b3b2f316))
- **i18n:** fix misleading text in new commits badge on repository card ([cd45a07](https://github.com/specvital/web/commit/cd45a07a9b1fc64f7565830c1882587214a0353c))
- **i18n:** improve spec regeneration modal wording ([cdfcd0c](https://github.com/specvital/web/commit/cdfcd0c795ae61491a3fc7c47e5fdbeffae9e042))
- **i18n:** spec document dates ignore app locale, use browser default ([2f82f98](https://github.com/specvital/web/commit/2f82f98918e7b1b7cf1dbe05842d8da8ce50d92d))
- **i18n:** spec regeneration modal description mismatches actual behavior ([c419e8e](https://github.com/specvital/web/commit/c419e8efdb1bff6927568d628ac24fb9afa7e9f8))
- missing loading feedback on dashboard navigation after background task completion ([c88c609](https://github.com/specvital/web/commit/c88c609b000fc43eb0b684dda5c50e576624c237))
- **mobile-nav:** fix bottom bar items cut off on narrow screens ([33908a7](https://github.com/specvital/web/commit/33908a7d66a1044d6bffcd1e7bcaf360f3e55c96))
- **mobile:** fix floating sidebar buttons appearing above search modal ([5c3ea59](https://github.com/specvital/web/commit/5c3ea59ad9976361a9d23ccfe95e46d28b764f98))
- **pricing:** correct subscription plan prices ([0610be9](https://github.com/specvital/web/commit/0610be9b3ed99b3c74331a6a229aead722d6adb8))
- **pricing:** fix "Current Plan" incorrectly shown on Free card for Pro users ([f0b52e6](https://github.com/specvital/web/commit/f0b52e600a4a5f75cb0b1fc1be6c2a52da9b5b3f))
- **pricing:** fix uneven spacing on first/last FAQ accordion items ([3d6e461](https://github.com/specvital/web/commit/3d6e461485516b1985e59e2313768669c58bafed))
- **pricing:** improve FAQ terminology and explanations ([6eceb13](https://github.com/specvital/web/commit/6eceb13bb225ca174a4a294a1cff74203e817126))
- **pricing:** improve terminology and unit display in pricing and quota UI ([bfa8484](https://github.com/specvital/web/commit/bfa84841cc964bdcad791e12bfd0049af3204951))
- **pricing:** remove misleading CTA text for paid plans when logged out ([6c320fc](https://github.com/specvital/web/commit/6c320fcbb3789c383191b4a950fba556551596dd))
- **pricing:** resolve vertical alignment inconsistency and improve plan descriptions ([82389d2](https://github.com/specvital/web/commit/82389d2164ceba172c45f37289698d4e472122dd))
- **pricing:** resolve visual instability from inconsistent card sizes ([d8099d1](https://github.com/specvital/web/commit/d8099d1c9b54b1773a4c4a712412561e7fc9bb3a))
- **queue:** change queue name separator from colon to underscore for River compatibility ([8a6d4c1](https://github.com/specvital/web/commit/8a6d4c16b5d90d7182803bf58ed93611505e30d9))
- **queue:** separate dedicated queues for each River worker ([81dd507](https://github.com/specvital/web/commit/81dd50726aee14b730c6cc2bf9f9119ec70277d9))
- **spec-view:** add truncate and tooltip for long text in TOC sidebar ([50c8bb5](https://github.com/specvital/web/commit/50c8bb539ac1454f74e6150d89b1f9f6a1b65052))
- **spec-view:** add user-scoping to generation status API ([7e3c31d](https://github.com/specvital/web/commit/7e3c31dd5b6dbd6f79ee81e20b0ba54cc74a91c6))
- **spec-view:** auto-refresh document and close modal on generation completion ([252b40b](https://github.com/specvital/web/commit/252b40b9b4ed4c21d95060c999f7250cb6e8ca10))
- **spec-view:** dashboard not updating and missing toast on spec generation completion ([1fa8ef7](https://github.com/specvital/web/commit/1fa8ef7b4f363cf99f4a0da058ce0b217aa41c76))
- **spec-view:** fix 500 error and modal closing immediately on regeneration ([01a31a6](https://github.com/specvital/web/commit/01a31a6d69f21d9405068ab80ac5c09831f7ab57))
- **spec-view:** fix AI Spec document not displaying after reanalysis ([750ab45](https://github.com/specvital/web/commit/750ab4567a94cf66c17e3e16c8a99128915d9e5d))
- **spec-view:** fix AI Spec per-user personalization logic ([abd6b81](https://github.com/specvital/web/commit/abd6b81de82d29f2416c77f20d3a74835d18fba7))
- **spec-view:** fix Behavior items not navigable via Tab key ([4fdf29d](https://github.com/specvital/web/commit/4fdf29dbf66de1ccb9159b42a31d8a09a98b3d68))
- **spec-view:** fix document reverting to previous language after regeneration ([bd638e7](https://github.com/specvital/web/commit/bd638e754729107b10bf02ffcd9cc9bd020a8bb9))
- **spec-view:** fix document view mobile layout overflow ([ba7983e](https://github.com/specvital/web/commit/ba7983e4016ec58468e387c3b0796f9cd25bee99))
- **spec-view:** fix regeneration UI hidden when executiveSummary is missing ([30c68cd](https://github.com/specvital/web/commit/30c68cdbcfc7ba6b9579901501c25edbe3c62f8e))
- **spec-view:** fix TOC sidebar navigation not working with virtualized list ([e0b71c2](https://github.com/specvital/web/commit/e0b71c2f3258bf1d1288fde63b70de0f03daffae))
- **spec-view:** fix TOC sidebar truncating domains when list is long ([0b74941](https://github.com/specvital/web/commit/0b74941128a12ea40f773e7264b8f7cc171a679c))
- **spec-view:** improve generation state management and prevent duplicate generation per language ([a4f6112](https://github.com/specvital/web/commit/a4f61123f56650240a3951ea96697c996fc7ab49))
- **spec-view:** improve mobile FAB visibility and card layout ([7b3f54f](https://github.com/specvital/web/commit/7b3f54fbfe7c45bd0859afb7968e97cdfff234d8))
- **spec-view:** improve tooltip label readability for original test name ([e0e6404](https://github.com/specvital/web/commit/e0e6404840684511bbb635be3a7ad96a52c19a2c))
- **spec-view:** keep polling during River job retry ([cec2f4b](https://github.com/specvital/web/commit/cec2f4b899f2cab856334a3f3e0f1f4683724492))
- **spec-view:** quota confirm dialog always opening with English instead of selected language ([2548515](https://github.com/specvital/web/commit/254851576702e68178b92268559e6f98f9df4fe6))
- **spec-view:** resolve multiple spec generation loading feedback UI defects ([b2ba622](https://github.com/specvital/web/commit/b2ba622e94539356a9e34de436ae71ff460744bf))
- **spec-view:** restore domain card structure in virtualized view ([9c41475](https://github.com/specvital/web/commit/9c414750b7a218ed6bf770d4a33bc239e25d8d6d))
- **spec-view:** restore missing gaps between cards in virtualized view ([a155369](https://github.com/specvital/web/commit/a155369cc0a7252167778d0966b37ca22aa38ad4))
- **spec-view:** return 403 instead of 500 for users without subscription ([778a4c1](https://github.com/specvital/web/commit/778a4c11b8401cbc14729d1bf63c52c72aef2b1f))
- **spec-view:** show cache reuse option when regenerating spec document ([943f5bd](https://github.com/specvital/web/commit/943f5bd04ec4f79734bbc364d867187dd8956f45))
- **spec-view:** show documents immediately after AI generation completion ([469dcce](https://github.com/specvital/web/commit/469dcce2960391fe1da8ce22f84333970e845d45))
- **spec-view:** suppress old version banner flash during document regeneration ([79d3798](https://github.com/specvital/web/commit/79d379893eca3c1832ef733d08802af7aec32984))
- **spec-view:** unify version badge and dropdown height with adjacent buttons ([88c14f4](https://github.com/specvital/web/commit/88c14f4fa800a8043d117fcc97bec72198eae176))
- stale quota displayed in re-analyze modal ([e855165](https://github.com/specvital/web/commit/e8551656583dba477e2aa5403a99d989bcd7fad8))
- **subscription:** fix users getting 2 months quota when signing up late in month ([ec64e42](https://github.com/specvital/web/commit/ec64e4286aa93e4963adb17a2afb96ed989a30cd))
- **tasks-dropdown:** long repository names unidentifiable due to truncation ([6fe9bf9](https://github.com/specvital/web/commit/6fe9bf9adff308386af65cb671e9c156606b6b21))
- **ui:** add cursor pointer to analysis tab buttons ([099ebb9](https://github.com/specvital/web/commit/099ebb90355750ef44cb49654182c9c6ab99f7f5))
- **ui:** add missing cursor-pointer to modal close button ([070c520](https://github.com/specvital/web/commit/070c5200e0047f02ef1dae2845c8f6670ef35f22))
- **ui:** apply pointer cursor to dropdown menu and command palette items ([1b2d3d7](https://github.com/specvital/web/commit/1b2d3d7ee9edd7b9b9b9dce56d73cde2cbe97679))
- **ui:** dropdown category headers indistinguishable from selectable items ([b443746](https://github.com/specvital/web/commit/b443746229564cdf0f052748fb5271792c79a291))
- **ui:** fix stepper connector alignment in spec generation modal ([107ea1b](https://github.com/specvital/web/commit/107ea1bb5ef841212977e7a8bae08efc33e91655))
- **ui:** fix tooltip secondary text invisible on dark backgrounds ([e03279f](https://github.com/specvital/web/commit/e03279f5d9e92a120c62a8afb0cdba08873ce95b))
- **ui:** fix vertical alignment of repository name and bookmark button in dashboard ([27ff606](https://github.com/specvital/web/commit/27ff60628d81372a6bb1183d64a91a5f0d3bfa3a))
- **ui:** prevent pointless spec regenerate on same commit ([b2fada3](https://github.com/specvital/web/commit/b2fada361f0349be77407ccd9da57db594bd30cd))
- **ui:** remove inaccurate warning in spec regeneration modal ([345d1f1](https://github.com/specvital/web/commit/345d1f1c1bb34352fcff3e62fae391b1d1d23b12))
- **ui:** unify card grid layout between Dashboard and Explore pages ([9846f73](https://github.com/specvital/web/commit/9846f73f75d5b5d28921b31b87e2c8dec85026a5))

#### ⚡ Performance

- **analysis:** reduce polling interval for better status capture ([ef8a75c](https://github.com/specvital/web/commit/ef8a75c95b61246a5dda80e19953d52d6932e052))
- **auth:** remove backend API call from homepage redirect ([8284e0c](https://github.com/specvital/web/commit/8284e0cd8609746843fc9f4b8c12983692ea35c5))
- **spec-view:** implement window-level virtualization for large document performance ([71fce34](https://github.com/specvital/web/commit/71fce346fd2b382207cef2a6c1404d97af8ce8aa))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **ci:** add debug step to diagnose E2E workflow server startup failure ([085c01b](https://github.com/specvital/web/commit/085c01b41a12f2217507830e3b4bdd5af108db84))
- **ci:** add missing GitHub App env vars for server startup ([fd8d4cd](https://github.com/specvital/web/commit/fd8d4cd8c0b17e097b87767c05dcded51dde8eff))
- **ci:** fix E2E workflow server startup errors ([b254989](https://github.com/specvital/web/commit/b2549896c28650c5c6261c46660b1f79d6b28c4f))
- **ci:** fix E2E workflow server startup failure ([2bf0079](https://github.com/specvital/web/commit/2bf0079e5af44e54e4ed2202501daecdbe870bd7))
- **ci:** fix ENCRYPTION_KEY length (31 bytes → 32 bytes) ([6315b50](https://github.com/specvital/web/commit/6315b5040e61e2063c880224798914b944602285))
- **ci:** fix ENCRYPTION_KEY length error (33 bytes → 32 bytes) ([9abe930](https://github.com/specvital/web/commit/9abe930b0cf7eeb80f0ee44671524ca4e9111f85))
- **ci:** fix RSA key format (YAML literal block → escaped string) ([44ce90c](https://github.com/specvital/web/commit/44ce90cb2a6511301585e71c17a0fab0258b6e8f))
- **ci:** fix schema load and add server startup debugging ([35bd1eb](https://github.com/specvital/web/commit/35bd1eb5ff628e23dc58b20fcfc79a9dc450c834))
- **ci:** use go run directly instead of air in CI to fix server startup failure ([d4f23b2](https://github.com/specvital/web/commit/d4f23b26a0359057abf9f40ef907e4aa6d6c6476))
- **e2e:** fix 15 E2E test failures and resolve post-test hang issue ([a4fbec9](https://github.com/specvital/web/commit/a4fbec9cee361adfda1c29fb7323d500e00dfe02))
- **e2e:** fix E2E test selectors and timing issues in CI ([1b80107](https://github.com/specvital/web/commit/1b801077f8d70277b7cc6492bf1611d5192e7e09))
- **e2e:** fix E2E tests failing in CI environment ([95c4355](https://github.com/specvital/web/commit/95c4355a82b8b0422fb0828739400841ff500feb))
- **e2e:** fix global-search and spec-view E2E test failures ([44e3a75](https://github.com/specvital/web/commit/44e3a7550ddcb06bc93ebf778a88e283228a7535))
- **e2e:** improve E2E test stability in CI environment ([1c0a372](https://github.com/specvital/web/commit/1c0a372d4a70e71ec54d2db562ac8431b6e85c6a))
- fix Railway environment variables rate limit issue ([cda35b0](https://github.com/specvital/web/commit/cda35b0ceb387fa938c92699f3dd352d1afe55e3))

#### 📚 Documentation

- add E2E test maintenance rule to CLAUDE.md ([77b6208](https://github.com/specvital/web/commit/77b6208a67992073e0e1ac97b25977d6009969c0))
- add specvital-specialist agent ([cdbb296](https://github.com/specvital/web/commit/cdbb29622129ee5bbef2f6b854303f821ca55042))
- correct usage docs inaccuracies and add missing cache feature documentation ([b59fd6c](https://github.com/specvital/web/commit/b59fd6c2fdcf90125696b99638cbf1242e09317c))
- **github-access:** remove GitHub access documentation page ([aa8a8ec](https://github.com/specvital/web/commit/aa8a8ec23761201d61504bbf5ddce7e61f09c507))
- **queue-processing:** remove queue processing documentation page ([e3d2dcd](https://github.com/specvital/web/commit/e3d2dcdd2642886e74554023df6bfaea099bf0ec))
- remove implementation-detail content from documentation pages ([4fd3300](https://github.com/specvital/web/commit/4fd3300b1929675a88470c6f1b264c3e56a2e836))
- simplify spec documentation ([de90097](https://github.com/specvital/web/commit/de9009743cff24e9f054ab0b1481b5c3b897726b))
- **usage-billing:** remove redundant check usage section ([f38addb](https://github.com/specvital/web/commit/f38addb8a86024bb75e35aa9b4a0c76f0675a2c5))

#### 💄 Styles

- **docs:** remove decorative icons from documentation page headers ([120cb2b](https://github.com/specvital/web/commit/120cb2b2eeef8b1283a0ad0601b187645a01fd55))
- format code ([e136032](https://github.com/specvital/web/commit/e136032899a3b4315fb1215fbdbec6901884b94a))
- translate Korean comments to English ([c46d61f](https://github.com/specvital/web/commit/c46d61f809e9047b7083c646c3894f8e3f38414b))

#### ♻️ Refactoring

- **analysis:** move Export button from Header to Tests tab toolbar ([47b9772](https://github.com/specvital/web/commit/47b977246fc7c8a4cfa86a2d03a3bb6f4504893a))
- **background-tasks:** remove unused polling-manager infrastructure ([188f8d1](https://github.com/specvital/web/commit/188f8d108feabfe8fd35ae9ef8ebea0d68793abe))
- **dashboard:** migrate use-reanalyze hook to TanStack Query polling ([2e1e6df](https://github.com/specvital/web/commit/2e1e6df471072917beb6ea9897badf7430e0646a))
- **inline-stats:** apply conditional rendering to todo status ([833fe07](https://github.com/specvital/web/commit/833fe07bdea51b37ee6003350a8f24fdac1bf572))
- **query:** replace invalidation-registry with native invalidateQueries ([34e77b9](https://github.com/specvital/web/commit/34e77b9465d5130ad4fa722a2d8ddd38dbd67cc9))
- remove dead code (unused exports, duplicate functions) ([851d430](https://github.com/specvital/web/commit/851d4304fbba6aab8a173878725eb6f2394685fa))
- remove unused dead code ([83970dd](https://github.com/specvital/web/commit/83970ddb448c9bf7ac432b6a27f364f0b80d6b65))
- **spec-view:** declutter spec generation modal by removing redundant elements ([9995b44](https://github.com/specvital/web/commit/9995b44771907d3e794ff1a1b3d1f52a62385c5f))
- **spec-view:** remove duplicate filter UI from document view ([0307d3b](https://github.com/specvital/web/commit/0307d3bd5ae8b0ed0553c0978150254922047db0))
- **spec-view:** replace isForceRegenerate boolean with generationMode enum ([1b94e07](https://github.com/specvital/web/commit/1b94e07a3bf666ff6bd4f39edde5905c1dc9a340))
- **spec-view:** replace useEffect-based state sync with derived state for spec generation ([2283c0c](https://github.com/specvital/web/commit/2283c0cf0f92a589287cce9fb379ceb7695b7069))
- **state:** migrate all stores to zustand ([080228f](https://github.com/specvital/web/commit/080228f7dca9e97981b5cf0d66dff87a5381c784))
- **status-counts:** support individual counting for all 5 test statuses ([930468a](https://github.com/specvital/web/commit/930468af78952c5649c114de78a236dc19a39146))
- sync queries with 4-tier schema via test_files table JOIN ([aa965df](https://github.com/specvital/web/commit/aa965df0368bdc223b16972060fe38eed909f279))
- sync with collector→worker rename and analyzer separation ([cd5a738](https://github.com/specvital/web/commit/cd5a7388b0029a420b565738abc2aadab01fc0d3))
- **ui:** consolidate FilterEmptyState into shared component ([4025a7c](https://github.com/specvital/web/commit/4025a7cb86f2a33182f2156c6507ddbc874e391d))
- **ui:** reuse RepositoryCard component in MyReposTab ([98a44e3](https://github.com/specvital/web/commit/98a44e3979de8aef3eb94b11cdae1b091016777e))
- **ui:** reuse RepositoryCard component in OrgReposTab ([47ce763](https://github.com/specvital/web/commit/47ce7634b938abe0e8c78d45a3c41ae5ff5739ca))

#### ✅ Tests

- **backend:** update tests to match sortBy mismatch graceful restart behavior ([a13c721](https://github.com/specvital/web/commit/a13c7218b80adbb130086db48511004128261449))
- **e2e:** add API mocking infrastructure and mocked E2E tests ([d116c53](https://github.com/specvital/web/commit/d116c53fb1236817fde2a242e5031ebcae2b8ab8))
- **e2e:** add authenticated page E2E tests ([53f7aee](https://github.com/specvital/web/commit/53f7aee9cae312478c5ddc8fe185cfde972ad6d0))
- **e2e:** add E2E tests for docs pages and spec view features ([8b9520d](https://github.com/specvital/web/commit/8b9520d2551216b75796bdf7cc35da34c522c061))
- **e2e:** add E2E tests for focused/xfail conditional display ([9bbab0b](https://github.com/specvital/web/commit/9bbab0bae6c0acab242cd033ccfb73476c126977))
- **e2e:** add E2E tests for sorting, bookmark, and AI spec features ([91c19a6](https://github.com/specvital/web/commit/91c19a67587da7d1b091c34d1b657d48e095932d))
- **e2e:** add polling behavior tests for analysis and dashboard ([b44ff6c](https://github.com/specvital/web/commit/b44ff6cc14979ee4c5722d2f4ebad62165f6a400))
- **e2e:** add Spec Generation, Background Tasks, and Analysis Polling UI tests ([0dea5a4](https://github.com/specvital/web/commit/0dea5a4c259e1abbdeea41e46142115f15fa6d69))
- **e2e:** expand authenticated E2E tests and add new test suites ([de7c22e](https://github.com/specvital/web/commit/de7c22e5051407dc89a25212ce3cdba128fb9096))
- **e2e:** implement 14 UI E2E test scenarios ([9283dd9](https://github.com/specvital/web/commit/9283dd91dcad0ae2128864cc16a844832cf372f2))
- **e2e:** implement 14 UI E2E test scenarios ([f676787](https://github.com/specvital/web/commit/f676787e08f7e5a4f982d7f89911aa380c47f384))
- **e2e:** implement skipped plan limits tests ([95627ab](https://github.com/specvital/web/commit/95627ab5bff6598684b2af88dcfda9291b9ddc86))
- **e2e:** implement skipped tests and add Analysis page E2E tests ([37e360c](https://github.com/specvital/web/commit/37e360c5db3fad3dc5bde8cdb86db2bff96eb5cf))
- **e2e:** remove obsolete Re-analyze button test ([1f566a2](https://github.com/specvital/web/commit/1f566a25cbce3006982610d04295513a2c9b841d))

#### 🔧 CI/CD

- add E2E tests workflow with sharded parallel execution ([9929d13](https://github.com/specvital/web/commit/9929d131c80ea2329d801d5b3ea86419a7abed55))
- migrate Railway deployment to IaC ([c9ee36f](https://github.com/specvital/web/commit/c9ee36f15e70cc1aecd46d6df9e236a052903d65))

#### 🔨 Chore

- add e2e test command ([ec42d0b](https://github.com/specvital/web/commit/ec42d0b9f4a76f8c2d76c44def9f4c013aa331b0))
- add e2e-ui docs in gitignore ([6d07cfb](https://github.com/specvital/web/commit/6d07cfbbab0ddb917caa284533c9055d3b1a9bdf))
- add mock spec-generator run command ([4650d26](https://github.com/specvital/web/commit/4650d268ce22540d8f29b4fcada024acf2d85069))
- add seed data ([2b4d7a4](https://github.com/specvital/web/commit/2b4d7a436e7ba9b1b9bf63e859c27b34cf37d477))
- change license from MIT to Apache 2.0 ([c407415](https://github.com/specvital/web/commit/c407415e7c3b0999c2af6c5cf6783f87573bc87e))
- dump schema ([c9c8bc2](https://github.com/specvital/web/commit/c9c8bc214326ddf3626c4a8561f253f4d6acf8e9))
- dump schema ([0dc9090](https://github.com/specvital/web/commit/0dc90901374ca1c376640a8bd69dee64adc06f3f))
- dump schema ([bbfca33](https://github.com/specvital/web/commit/bbfca33c1df0d7306d091e5a858426d875446ff8))
- dump schema ([1978405](https://github.com/specvital/web/commit/197840535c61b9019a635a8e57540f1831725455))
- dump schema ([3335616](https://github.com/specvital/web/commit/33356163e7ebe343dee7195a4d918f2d8e566844))
- dump schema ([1ffbd05](https://github.com/specvital/web/commit/1ffbd0572b28fda8b30d71c33dbb3c54174a4547))
- dump schema ([d8968d7](https://github.com/specvital/web/commit/d8968d7ea6b1a3b38adf9c888af9debb9a1dd05c))
- dump schema ([14893e0](https://github.com/specvital/web/commit/14893e03db04c3f2c40935632f0577a565057a32))
- dump schema ([3c78248](https://github.com/specvital/web/commit/3c782487452eb78e14b61176b021b0b30ebe9f47))
- hotfix port for load container ([4cdc61d](https://github.com/specvital/web/commit/4cdc61da13505b77799f38a6d913e71c4ef71959))
- remove commit_message.md ([1fb66e6](https://github.com/specvital/web/commit/1fb66e632c495d37283e246b110209a22697538c))
- remove spec-view feature and document AI integration patterns ([a5475bf](https://github.com/specvital/web/commit/a5475bf1ef3056b5d782c325ae04a7fe1fd97eb8))
- remove unnecessary backfill script ([61ad7a9](https://github.com/specvital/web/commit/61ad7a9117a2f80550148c62c4709b39743737c3))
- **spec-view:** add run-spec-generator command for local development ([8f84b48](https://github.com/specvital/web/commit/8f84b4833f657daf626df05d61fa0198af40a794))
- sync ai-config-toolkit ([4eb1b3f](https://github.com/specvital/web/commit/4eb1b3fa31e9748488c3811167979df3215c685c))
- sync ai-config-toolkit ([0290023](https://github.com/specvital/web/commit/02900237dbc6c37b8ebc55115072facae8457ad0))
- sync ai-config-toolkit ([2342fe7](https://github.com/specvital/web/commit/2342fe7df26d3e218f01e3a312d30b1b2fc581e6))
- sync docs ([7bda672](https://github.com/specvital/web/commit/7bda672fe78136150284a86772fd8fe9b5b13b87))
- sync seed specview monthly limits with production ([37dd459](https://github.com/specvital/web/commit/37dd459c3ca9447ece4d16a239c2880844ae8fd8))
- sync-docs ([5df1c72](https://github.com/specvital/web/commit/5df1c7297b9fec1327c2b874cde59129f1044495))
- unify oapi-codegen version to v2.5.1 ([e3fff92](https://github.com/specvital/web/commit/e3fff921cba4f1f61813821eb5ca2edef80835be))

## [worker/v1.2.4](https://github.com/specvital/worker/compare/v1.2.3...v1.2.4) (2026-02-02)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **fairness:** resolve user tier from DB instead of job args ([527c1ae](https://github.com/specvital/worker/commit/527c1aef75771ff026813e5e39faa076f42e23ee))

## [worker/v1.2.3](https://github.com/specvital/worker/compare/v1.2.2...v1.2.3) (2026-02-02)

### 🔧 Maintenance

#### ♻️ Refactoring

- **deploy:** separate Dockerfiles per service to remove buildArgs dependency ([a142f10](https://github.com/specvital/worker/commit/a142f108feeaf33991377133cb51e72a57064cbc))

## [worker/v1.2.2](https://github.com/specvital/worker/compare/v1.2.1...v1.2.2) (2026-02-02)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **deploy:** apply railway.json config during CLI deployment ([f0f5e79](https://github.com/specvital/worker/commit/f0f5e791ecd897ba692b49a03b8662f97a5be982))

#### ♻️ Refactoring

- **deploy:** move Railway config files to infra/railway/ ([bcd7a87](https://github.com/specvital/worker/commit/bcd7a87aeff8d968397a8115dcae313c3173d1d8))

## [worker/v1.2.1](https://github.com/specvital/worker/compare/v1.2.0...v1.2.1) (2026-02-02)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **deploy:** remove railway link command for Project Token usage ([50e4ac1](https://github.com/specvital/worker/commit/50e4ac1be248860f5a08ea88f48f4f12ef50f4df))

#### ♻️ Refactoring

- **deploy:** reorder workflow to run release after successful deployment ([404b956](https://github.com/specvital/worker/commit/404b95676af528a62449b9712958a261304c732e))

#### 🔨 Chore

- add railway cli ([98bf945](https://github.com/specvital/worker/commit/98bf945d2aee5952d7e7da41554f805869b3820b))

## [worker/v1.2.0](https://github.com/specvital/worker/compare/v1.1.1...v1.2.0) (2026-02-02)

### 🎯 Highlights

#### ✨ Features

- **analysis:** add usage_events recording on analysis completion ([0593683](https://github.com/specvital/worker/commit/05936838603fb478da028d87228001bab21c48a5))
- **analysis:** store parser_version in analysis records ([aa47dab](https://github.com/specvital/worker/commit/aa47dab02fc13401415ec0bd7485e0c4ba473b48))
- **autorefresh:** add parser_version comparison to auto-refresh conditions ([4cc6cc4](https://github.com/specvital/worker/commit/4cc6cc4323a87d1702448bb3c94a42af2a10cab2))
- **batch:** implement automatic chunking for large test sets ([cf7c691](https://github.com/specvital/worker/commit/cf7c6914ad71273c5e783f7a94b812be2edba3a9))
- **bootstrap:** add multi-queue subscription support ([dc4c254](https://github.com/specvital/worker/commit/dc4c254d2194a5faa9d710452ddc652d376e040d))
- **bootstrap:** register parser version to system_config on Worker startup ([290a5ef](https://github.com/specvital/worker/commit/290a5efa4e3edf928c68fb9e08d1e6ea41380f77))
- **config:** add feature flag infrastructure for Phase 1 V2 architecture ([9c1c4c9](https://github.com/specvital/worker/commit/9c1c4c9e6fa5aab1558bd490c3cc0f43263f83b4))
- **db:** add SQLc queries and domain model for test_files table ([ca530c4](https://github.com/specvital/worker/commit/ca530c4162a79b589b2945b7d58b4fdcf5781f72))
- **deploy:** add retention-cleanup service Railway configuration ([876c3f3](https://github.com/specvital/worker/commit/876c3f37b14cc0c7bf041ded6c197b0c619f70b7))
- **deploy:** migrate to Railway IaC with GitHub Actions deployment ([b842318](https://github.com/specvital/worker/commit/b842318fe5a3f6d7c56780dd93979299382030e4))
- **gemini:** add auto-recovery for missing test indices from AI output ([c3173e8](https://github.com/specvital/worker/commit/c3173e8a08fcbaca230463ec26f9767b88a5ee32))
- **gemini:** add domain description field to v3BatchResult ([62d5b8d](https://github.com/specvital/worker/commit/62d5b8d6894ef9ade164af83315b41dcfe8f054a))
- **gemini:** add Phase 1 V3 sequential batch architecture feature flag ([71fc07e](https://github.com/specvital/worker/commit/71fc07eb4e3f0fd4f8538167083e1325592f53fa))
- **gemini:** add Stage 1 Taxonomy cache ([7ccea47](https://github.com/specvital/worker/commit/7ccea4782b024e935caeb7bd4c7fbb7b10a4ec3e))
- **gemini:** add V3 prompt templates for order-based batch classification ([b186601](https://github.com/specvital/worker/commit/b1866013a2368ed12365e6d52c3907ad657e27d2))
- **gemini:** implement Batch API polling and result parsing ([66d7ddf](https://github.com/specvital/worker/commit/66d7ddf4e2bcb13eadfcc0e1bc0c8851d83da28a))
- **gemini:** implement Batch API Provider base structure ([7214cbd](https://github.com/specvital/worker/commit/7214cbd4e647f600739550708b0a0e660cd9b418))
- **gemini:** implement Phase 1 V2 Orchestrator ([f2f453b](https://github.com/specvital/worker/commit/f2f453b0fd8edd3720ca9fb44161705006676460))
- **gemini:** implement Phase1 quality metrics collector ([5aef0df](https://github.com/specvital/worker/commit/5aef0dfa62ea17769a63a330ed4f066d8b9e679e))
- **gemini:** implement Phase1PostProcessor for classification validation and normalization ([c7fee4e](https://github.com/specvital/worker/commit/c7fee4e241d5deebac4a0d51dad65bfc03985779))
- **gemini:** implement Stage 1 Taxonomy extraction logic ([be9bf50](https://github.com/specvital/worker/commit/be9bf503514e87783edd80a9016f2145bce96a52))
- **gemini:** implement Stage 2 Test Assignment logic ([531962b](https://github.com/specvital/worker/commit/531962bc38dd771af6b9d5dd4ec1bdcc16d759d3))
- **gemini:** implement V3 batch processing core logic with response validation ([c528346](https://github.com/specvital/worker/commit/c5283462c699e1c449d36bee9b689b15a0b19dcb))
- **gemini:** implement V3 batch retry logic with split and individual fallback ([0945156](https://github.com/specvital/worker/commit/094515633c8405ad8db358406561ce457231eefd))
- **gemini:** implement V3 orchestrator for sequential batch processing ([1df43b2](https://github.com/specvital/worker/commit/1df43b2773997d4d4e499c1ff4daaef6e008723d))
- **gemini:** integrate Phase 1 V2 feature flag router ([cf72f02](https://github.com/specvital/worker/commit/cf72f024894573d6462d13f1a82c18b7be3efe62))
- **gemini:** integrate quality components into V3 orchestrator ([513232b](https://github.com/specvital/worker/commit/513232b174e31e113d6537bc1267282f73074da2))
- **mapping:** add Core DomainHints to Domain model conversion ([570e07a](https://github.com/specvital/worker/commit/570e07ab788660352dcb2444cafb3942e606b046))
- **mock:** add response delay simulation to mock AI provider ([78aa6c4](https://github.com/specvital/worker/commit/78aa6c400478a72c1ca28cdd87ea0f2aa135cbb6))
- **prompt:** implement Stage 1 Taxonomy extraction prompt ([eb316c3](https://github.com/specvital/worker/commit/eb316c325e653a6df9b5f2c2c1c07e6969b2b318))
- **prompt:** implement Stage 2 Assignment prompt ([c9775d1](https://github.com/specvital/worker/commit/c9775d1b9c1d96383021e0115a6f6132884d6553))
- **queue:** add fairness middleware for per-user job limiting ([620849f](https://github.com/specvital/worker/commit/620849f7f212e02138f719fb061c9ec1e89454a2))
- **queue:** add multi-queue configuration support for River server ([71274b4](https://github.com/specvital/worker/commit/71274b4a414ee2dc68be6e0f62d3b41355629f82))
- **queue:** add per-user concurrent job limiting ([421904c](https://github.com/specvital/worker/commit/421904c4f4b0dd3732a18a990ae867da91deebad))
- **queue:** add priority queue name constants and environment config ([9a6df10](https://github.com/specvital/worker/commit/9a6df10743bdad5dde74aa394f39ace2fad0c0f6))
- **queue:** implement userID and tier extraction from River jobs ([b14baf5](https://github.com/specvital/worker/commit/b14baf5e91414cfd38a86b4800b51951111836d9))
- **queue:** integrate Batch mode into SpecView Worker ([5cb5155](https://github.com/specvital/worker/commit/5cb5155601f032d8cdc1c7ed1822c992147e31fd))
- **queue:** integrate fairness middleware into River server ([2206b19](https://github.com/specvital/worker/commit/2206b192c6ce0c008b7a0450c4fba6d7a67e63cc))
- **quota:** release quota reservation on job completion/failure in Worker ([8c69599](https://github.com/specvital/worker/commit/8c695993816871eac19ffde14961a02fed540ac7))
- **repository:** add system_config repository and parser version extraction ([72bd468](https://github.com/specvital/worker/commit/72bd4681216d5dffefba222d94f62c6cc60799b0))
- **retention:** add bootstrap and entry point for retention cleanup ([6e03a7f](https://github.com/specvital/worker/commit/6e03a7f81829b310d43a5c0c575f066fac20ec1f))
- **retention:** add domain and usecase layer for retention cleanup ([5e8c05a](https://github.com/specvital/worker/commit/5e8c05ae44ae5f263af7d0466c239959bffb1612))
- **retention:** add repository layer for retention-based cleanup ([792f106](https://github.com/specvital/worker/commit/792f106303968946ecd8ebcb22a72b8dccdb9fb5))
- **retention:** store retention_days_at_creation snapshot on record creation ([878d87c](https://github.com/specvital/worker/commit/878d87c9fde72ee44f2b0f61b7dd56e3b81dc3c0))
- **scheduler:** add dedicated queue support for auto-refresh jobs ([ba91bd9](https://github.com/specvital/worker/commit/ba91bd994ed3914a767eea0d092961263e2a9972))
- **specgen:** add MOCK_MODE env var for spec document generation without AI calls ([364dfda](https://github.com/specvital/worker/commit/364dfdab35f0d5b7ca67969791594340fadc8817))
- **specview:** add analysis_id and ETA to Phase 2 logs ([f568df6](https://github.com/specvital/worker/commit/f568df690b523b3c9953378c9c02b01145f225bf))
- **specview:** add Batch API routing logic to UseCase ([36560f3](https://github.com/specvital/worker/commit/36560f3bbb63cc217056760d291c9ecab168f3a9))
- **specview:** add behavior cache stats to SpecViewResult ([a70fa07](https://github.com/specvital/worker/commit/a70fa073ee72eea6be19bb7571f317ff9da870a6))
- **specview:** add behavior cache types and interfaces ([a871845](https://github.com/specvital/worker/commit/a871845be14de5e8908adafd22d7ce5879df4568))
- **specview:** add domain layer foundation for SPEC-VIEW Level 2 ([3644739](https://github.com/specvital/worker/commit/3644739c794512b912f221a3e0184080570e4fda))
- **specview:** add domain types for Phase 1 V2 two-stage architecture ([cc75bb5](https://github.com/specvital/worker/commit/cc75bb5a685a48969a95a8127bb7ffe720cd5faf))
- **specview:** add force regenerate option ([09821f9](https://github.com/specvital/worker/commit/09821f9f278507409880690579603406a855b913))
- **specview:** add Phase 1 chunking for large test repositories ([8831021](https://github.com/specvital/worker/commit/88310213b6b563e79c39d34b26571c7db55b8288))
- **specview:** add Phase 1 classification cache foundation types ([42a2a47](https://github.com/specvital/worker/commit/42a2a475480e3f9122c695011cc5154754899416))
- **specview:** add Phase 2 progress logging for job monitoring ([bfbb938](https://github.com/specvital/worker/commit/bfbb9384364fb2254de3d30c7e7e35401e9562ec))
- **specview:** add Phase 3 executive summary generation ([8ed1d84](https://github.com/specvital/worker/commit/8ed1d84f681be210d315d3cb6b3770a9d50c4222))
- **specview:** add phase timing logs for performance monitoring ([a04dc43](https://github.com/specvital/worker/commit/a04dc433b5f18d261d4ad8fda8c3b20755fe6d5d))
- **specview:** add placement AI adapter for incremental cache ([1cbefd6](https://github.com/specvital/worker/commit/1cbefd631b39218863b393c5a84ca6fdb425ed1d))
- **specview:** add real-time Gemini API token usage tracking ([109d572](https://github.com/specvital/worker/commit/109d57295ea501bd6a9a7179038a25584a64ff69))
- **specview:** add repository context (owner/repo) to logs ([ef57a38](https://github.com/specvital/worker/commit/ef57a38cd18507862ad4d4cc0e8830c0165f9d34))
- **specview:** add retry info and phase context to failure logs ([87a1993](https://github.com/specvital/worker/commit/87a1993a2aa275c7ab241e81e7dbc3210de9c72e))
- **specview:** add test diff calculation for incremental cache ([b61b78b](https://github.com/specvital/worker/commit/b61b78b26bd8c2353e448a3a0d3abff002b8b395))
- **specview:** add uncategorized handling for incremental cache ([b6b3ea5](https://github.com/specvital/worker/commit/b6b3ea5e78a5329148a21489b7bb9de3631d7b4a))
- **specview:** add user_specview_history recording ([abc148b](https://github.com/specvital/worker/commit/abc148b4c6e6fbe3a95f7317da1147b5ca2e23da))
- **specview:** add version management support for spec_documents ([a547803](https://github.com/specvital/worker/commit/a547803eddd022850d2eaf26ff26bd9c008e83e0))
- **specview:** implement behavior cache repository ([bea9a93](https://github.com/specvital/worker/commit/bea9a93dbd617b35f2b48bc77a45be57b674b36a))
- **specview:** implement classification cache repository ([ae3a8ba](https://github.com/specvital/worker/commit/ae3a8ba12f3d3c2bc71db772dcee878d3c0865f0))
- **specview:** implement Gemini-based AI Provider adapter ([fe5ea7d](https://github.com/specvital/worker/commit/fe5ea7dafa3619a6c6b195f5eaa8d5605d4178c0))
- **specview:** implement GenerateSpecViewUseCase for pipeline orchestration ([080978c](https://github.com/specvital/worker/commit/080978cf942463b5da339d49966d0234c20fe12e))
- **specview:** implement PostgreSQL repository for 4-table hierarchy ([e897bbf](https://github.com/specvital/worker/commit/e897bbf34c85c44ac9d0c0d6cd98804697552010))
- **specview:** implement SpecViewWorker and integrate with AnalyzerContainer ([f87e638](https://github.com/specvital/worker/commit/f87e638792c300207abbc696388fe621456b1507))
- **specview:** implement usage_events-based quota tracking ([94105c6](https://github.com/specvital/worker/commit/94105c61bb0e5fea6da9f4887348dbfdd2f0bae4))
- **specview:** integrate Phase 1 cache for incremental classification ([26a373b](https://github.com/specvital/worker/commit/26a373b9f180d8a630f41d6e1f5ab06fbfe425e2))
- **specview:** integrate Phase 2 behavior cache ([1799394](https://github.com/specvital/worker/commit/1799394d854dea2635fd173f0bb0b0beabbffe8c))
- **specview:** store user_id on spec_documents INSERT ([2ddac0d](https://github.com/specvital/worker/commit/2ddac0d487e3a62cbcdb1ee6f1baa1e6fb80bf6a))

#### 🐛 Bug Fixes

- **batch:** add debugging info for Batch API response parsing failures ([e9a6342](https://github.com/specvital/worker/commit/e9a6342b4cd6ed764163f5abf4d0b0b8e10c2f3e))
- **batch:** fix batch job repeatedly submitting instead of polling ([411e9bf](https://github.com/specvital/worker/commit/411e9bfc73212937bf4c402f46106da70c5e3426))
- **batch:** fix JSON corruption when parsing Batch API response ([93bbefb](https://github.com/specvital/worker/commit/93bbefb0d1fbf7c299cfada1dbe0acb784487fa6))
- **batch:** fix JSON parsing failure due to trailing commas in Gemini response ([cb6db9b](https://github.com/specvital/worker/commit/cb6db9b97a19b09722dae95d54548e5e9f0918f3))
- **batch:** fix parsing for Batch API split responses ([46e2302](https://github.com/specvital/worker/commit/46e2302b124133c73f5d9cae11db92d687272230))
- **gemini:** handle out-of-range test indices from AI without failing ([7d9382a](https://github.com/specvital/worker/commit/7d9382af1b5ff17c2d42562d0ac4455b48bb7d08))
- **gemini:** prevent AI hallucination causing invalid domain/feature pairs in Stage 2 assignment ([365f595](https://github.com/specvital/worker/commit/365f595ab425a046e8fbc1b82a5580d1051a6671))
- **gemini:** prevent context cancellation from truncating API responses in wave processing ([4dc421e](https://github.com/specvital/worker/commit/4dc421e38b83c8762789b731b75f6f142cd2da07))
- **gemini:** resolve taxonomy response truncation for large file sets ([0e71100](https://github.com/specvital/worker/commit/0e711005e12d159c4f6ae189db41fbe39ddf8e49))
- **prompt:** resolve Phase 2 output language mismatch with requested language ([ccea96e](https://github.com/specvital/worker/commit/ccea96e3a1ade5117220479c2d3d6bf7ce9a923b))
- **queue:** isolate dedicated queues per worker to resolve Unhandled job kind error ([3cfee6f](https://github.com/specvital/worker/commit/3cfee6fbaa5fcc1d31341d4cb9435e944cc5adfb))
- **queue:** replace colons with underscores in queue names ([0a43619](https://github.com/specvital/worker/commit/0a43619e869e3fd9ecf07566e5148d37ad8f6aa6))
- **specview:** charge quota only for AI-generated behaviors ([1cddcaa](https://github.com/specvital/worker/commit/1cddcaa5e27484c22a38327700fe846bb6457a13))
- **specview:** chunk cache not restoring on Phase 1 retry ([f4ef62d](https://github.com/specvital/worker/commit/f4ef62deb1e3418ddb0515b1ab1f2523dac0090a))
- **specview:** extend Phase 1 timeout for V3 sequential batch processing ([58bba48](https://github.com/specvital/worker/commit/58bba48ad4554924ace58d335a83c94f40d63cc6))
- **specview:** fix Batch API config not propagating to worker ([1c465c2](https://github.com/specvital/worker/commit/1c465c2cda2f229b4125ddd3e38045e2ee89e981))
- **specview:** isolate Phase 2 semaphore per job to prevent cross-job interference ([7668a98](https://github.com/specvital/worker/commit/7668a9855967367fb4bf5f0abf567efafef08945))
- **specview:** prevent behavior cache loss on Phase 2 failure ([0ec770b](https://github.com/specvital/worker/commit/0ec770bc05eca2b9207d6ba71307296f7fa5783d))
- **specview:** propagate Gemini env vars to AnalyzerContainer ([7485f08](https://github.com/specvital/worker/commit/7485f08ec2d9410b5d461cd765f545a9e5a241e4))
- **specview:** raise Phase/Job timeouts to allow large repo completion ([85089ad](https://github.com/specvital/worker/commit/85089ad6249ccc3b2c5a4ee4d826d3e7a59a3f85))
- **specview:** resolve Phase 1 response JSON truncation ([ac7d381](https://github.com/specvital/worker/commit/ac7d3819a76a6b570c2471cc2bd63199265ba3fd))
- **specview:** resolve Phase 2 output validation failures ([6c3ad64](https://github.com/specvital/worker/commit/6c3ad6403be48415acc1f9f5e22a54794bfcb8f7))

#### ⚡ Performance

- **gemini:** implement wave parallel processing to improve Phase 1 throughput ([cc33442](https://github.com/specvital/worker/commit/cc334420aafecfd15588b62f005f7e52cb37156e))
- **gemini:** reduce chunk size and enhance retry settings to minimize timeout risk ([e78ef82](https://github.com/specvital/worker/commit/e78ef82e9410036ef09a8bc3d41b176a05d04fe1))
- **gemini:** remove redundant inter-chunk delay to improve processing speed ([33a0d4d](https://github.com/specvital/worker/commit/33a0d4d8e6ffa60278cedf5ce9d6086f0c3d2926))
- **specview:** disable Gemini thinking mode to reduce Phase 1 timeout ([b47a4e9](https://github.com/specvital/worker/commit/b47a4e997abcb63923fbc25acf24c5951bf7dc50))
- **specview:** optimize prompts and timeouts to reduce Gemini 504 errors ([8857b90](https://github.com/specvital/worker/commit/8857b90e320e2a860e40f30480b26c0b4ea507bd))
- **specview:** reduce Phase 1 chunk size and add progress caching for large repositories ([7a08e01](https://github.com/specvital/worker/commit/7a08e01eaf61101e60e5329284c04ff7a2e207d2))
- **specview:** reduce Phase 1 chunk size to 500 and add JSON parse error retry ([c4edea8](https://github.com/specvital/worker/commit/c4edea8eb98ccbb22ebe27ca735e953a28abf958))

### 🔧 Maintenance

#### 📚 Documentation

- add build artifacts cleanup rule ([05e7b70](https://github.com/specvital/worker/commit/05e7b70cbeb8f3a07746bf481ce2b5641d62754c))
- add specvital-specialist agent ([9dded5d](https://github.com/specvital/worker/commit/9dded5d0e4c3610e766ee1e1bc8500eba55154c9))
- document Batch API environment variables and operation guide ([257e59b](https://github.com/specvital/worker/commit/257e59bfa797ff3b35e66f86fdca3396da868788))
- update spec-view docs ([e463e10](https://github.com/specvital/worker/commit/e463e1095a8f14e19aab40bb478dc33f1bdf8a40))

#### ♻️ Refactoring

- **app:** separate DI containers for analyzer and scheduler ([e3b1eae](https://github.com/specvital/worker/commit/e3b1eaed01306a349a75b5650886d9f1ec750e33))
- **gemini:** allow many-to-many file-feature relationship in Stage 1 taxonomy ([70bfa8f](https://github.com/specvital/worker/commit/70bfa8fa9892ee6a5e7e27deb87bacab7ae0df8a))
- **prompt:** redesign V3 classification prompt with principle-based approach ([ade4afc](https://github.com/specvital/worker/commit/ade4afc864c0c8104bcb1812e7c94fca0786c15f))
- **queue:** remove legacy queue support code ([e352f96](https://github.com/specvital/worker/commit/e352f96e083019e7a2e9eedef4f38a00e6b0310e))
- remove Scheduler service ([c163239](https://github.com/specvital/worker/commit/c163239f49c1602a285990a6f3a19b7f1c9459a3))
- **repository:** implement normalized storage using test_files table ([47cad88](https://github.com/specvital/worker/commit/47cad88accd953f2164ada321ba61f88ad791c19))
- separate queue adapter into analyze subpackage ([cf868bf](https://github.com/specvital/worker/commit/cf868bfb141a904117f6cc1ea6a722065dc3388a))
- separate worker binary into analyzer/spec-generator ([bd7f17f](https://github.com/specvital/worker/commit/bd7f17f8836d69e0c8b190f13aece43788f6413f))
- **specview:** improve Phase 2 prompt with Specification Notation style ([dbfc46d](https://github.com/specvital/worker/commit/dbfc46d30680fad6dd9f5916f24ac6422027250e))
- **specview:** improve Phase 3 prompt for user-friendly tone ([15f7e62](https://github.com/specvital/worker/commit/15f7e620753e588787d83ca9a049b33d77cbd5ec))
- **specview:** remove language constants and improve default handling ([31c1d40](https://github.com/specvital/worker/commit/31c1d409aace5dfee513dd3b7c4bac1a61644518))
- update module path and references for repository rename ([0629c97](https://github.com/specvital/worker/commit/0629c977b84525bf374ccadc14581967e8a71349))
- **worker:** separate AnalyzeWorker and SpecViewWorker into independent binaries ([f3fae45](https://github.com/specvital/worker/commit/f3fae45642a4459aed0df934891de39d0981cb1d))

#### ✅ Tests

- **gemini:** add integration tests for V3 sequential batch architecture ([30a23bc](https://github.com/specvital/worker/commit/30a23bcfbbffd484f4f7af5a3274cf4866e83d4a))
- **gemini:** implement V3 quality assurance integration tests ([79cd139](https://github.com/specvital/worker/commit/79cd139bc710321ddfa639796678e7bd54b06677))
- **gemini:** update V3 tests to expect path-based fallback instead of Uncategorized ([1964917](https://github.com/specvital/worker/commit/19649175cc0e99a550c250b5159501ce6ac445f2))
- **queue:** add integration tests for fairness middleware ([da41eea](https://github.com/specvital/worker/commit/da41eeaba9ac975ab5dca07098d4195b63271e05))
- **repository:** add integration tests for DomainHints and file_id FK relationship ([ddee5f6](https://github.com/specvital/worker/commit/ddee5f62df37ba23fc0b49138661bc0a2c99e398))
- **specview:** add E2E integration tests for SPEC-VIEW pipeline ([b803ca7](https://github.com/specvital/worker/commit/b803ca7bacc655850971501dbef7e1873070597a))

#### 🔨 Chore

- add air hot reload support for spec-generator ([c6697b2](https://github.com/specvital/worker/commit/c6697b2d8cd10ae9ffaf4017150afa731a92971d))
- add clean-containers command ([6a0d006](https://github.com/specvital/worker/commit/6a0d00608515bc5c56c57c16894420c20bf0c745))
- change license from MIT to Apache 2.0 ([393356c](https://github.com/specvital/worker/commit/393356cc417cfecd916e33aa13662ea55b5a1320))
- dumb schema ([9345080](https://github.com/specvital/worker/commit/9345080c67fd9667b65e3bd5b5b4f7594a43c588))
- dump schema ([641394d](https://github.com/specvital/worker/commit/641394d30eb4ef4f87cbd9e265f9f90199647445))
- dump schema ([ac3dd77](https://github.com/specvital/worker/commit/ac3dd77415cbadf5cd2f286e1c631bffaa27543e))
- dump schema ([5103aeb](https://github.com/specvital/worker/commit/5103aeb75f8ed56b64fdcf0f07de718358da0f57))
- dump schema ([33d8a46](https://github.com/specvital/worker/commit/33d8a46c38c8fdcfae4584b1d88705e4b658a64f))
- dump schema ([a61b1bb](https://github.com/specvital/worker/commit/a61b1bbd48b62f4f545cd84194ff8a828c509919))
- dump schema ([4450751](https://github.com/specvital/worker/commit/4450751cc76e29a0a753d10da559fa432f518013))
- dump schema ([cd3df33](https://github.com/specvital/worker/commit/cd3df33a6f74f855c1cda650516d015b88a86c3d))
- dump schema ([ea92b45](https://github.com/specvital/worker/commit/ea92b45e8c07c9d6cc13fbaa122d77fa9ae7a675))
- dump schema ([a5e6d73](https://github.com/specvital/worker/commit/a5e6d73cc218aa508d1b3d1c2c1f5655b3760d01))
- sync ai-config-toolkit ([ec52a6f](https://github.com/specvital/worker/commit/ec52a6f6eb3fb3da076414097437df03b72e53dc))
- sync docs ([8f9553b](https://github.com/specvital/worker/commit/8f9553b69ee56f7112fac89e0e42818f28ba7461))
- sync docs ([852cbfe](https://github.com/specvital/worker/commit/852cbfef6ce368fba59aeb4d6ba176cf3b17e931))
- sync-docs ([e4e8810](https://github.com/specvital/worker/commit/e4e881033d06e1d8b8a918b0d232bd83c80a3f79))
- update core ([8020c3b](https://github.com/specvital/worker/commit/8020c3b3046e8daf9f17b06f559ce2043474801f))
- update core ([dfd8e26](https://github.com/specvital/worker/commit/dfd8e26e2b9683bd180be4d624729f1742d7016e))
- update core ([a094e8d](https://github.com/specvital/worker/commit/a094e8deadd9466d8c244f8efc8e40e6d1b457b1))
- update core ([54ae38a](https://github.com/specvital/worker/commit/54ae38a02b6ec0ba1742b5ac5abd8141f26c7cad))
- update core ([ea3b03e](https://github.com/specvital/worker/commit/ea3b03e341e8538f8690b8da9ba53f327e89a3ca))
- update dev tool configs for worker→analyzer refactoring ([f1c4637](https://github.com/specvital/worker/commit/f1c4637b56e181246f5cf9e9a76fcc35d6da0f31))

## [infra/v1.3.0](https://github.com/specvital/infra/compare/v1.2.0...v1.3.0) (2026-02-02)

### 🎯 Highlights

#### ✨ Features

- add schema visualization and documentation tools ([f92cd2c](https://github.com/specvital/infra/commit/f92cd2c3415089175ecf7e420c1f4ecf9da81f02))
- **db:** add behavior_caches table for Phase 2 caching ([8917156](https://github.com/specvital/infra/commit/8917156eeb8d1a5fe8a2c3bc44cc19d5a87a5756))
- **db:** add classification_caches table for Phase 1 incremental caching ([87362bd](https://github.com/specvital/infra/commit/87362bd2de2415d9be663a8f4ee9517d0d2f5fc7))
- **db:** add monthly_price column to subscription_plans ([63badbc](https://github.com/specvital/infra/commit/63badbc5e0ca605fca825d89ded6ea4ec9294187))
- **db:** add parser version tracking for re-analysis support ([a681e0d](https://github.com/specvital/infra/commit/a681e0d08f1640b19b073e240b700b929b8feeca))
- **db:** add quota_reservations table for concurrent request handling ([d3f15bf](https://github.com/specvital/infra/commit/d3f15bf9108c8e52f4193a7e2b17bdc29f011c25))
- **db:** add retention_days_at_creation for creation-time retention policy ([7eb93aa](https://github.com/specvital/infra/commit/7eb93aac632ed03063451706da9beb07a6c8b26e))
- **db:** add schema support for spec_documents versioning ([1b79497](https://github.com/specvital/infra/commit/1b79497f887439db0e5d2233492505a10f8dffcc))
- **db:** add spec_view_cache table for AI conversion results ([9548a69](https://github.com/specvital/infra/commit/9548a69cb1d5aceaa344edb5e2d750d03afb8e03))
- **db:** add subscription plans schema for usage limits ([efda341](https://github.com/specvital/infra/commit/efda341fbd9358c9b0cb7ff04ca3959733bed5e3))
- **db:** add test_files table for schema normalization ([58b7c5b](https://github.com/specvital/infra/commit/58b7c5b9df1285d89f57bc41ad7b3cf5b7e30f7a))
- **db:** add usage_events table for quota tracking ([e0dc198](https://github.com/specvital/infra/commit/e0dc19852d0afcb88a2f98f2b2f0e4a94ccf8832))
- **db:** add user_id column to spec_documents table ([d861022](https://github.com/specvital/infra/commit/d861022e5aa8c13e66725a9ae84052d4f83a156c))
- **db:** add user_specview_history table for tracking SpecView generation ([a196258](https://github.com/specvital/infra/commit/a1962584afbd32f27d499e5874ee4217da04d1c8))
- **db:** replace spec_view_cache with hierarchical spec document schema ([38a33ad](https://github.com/specvital/infra/commit/38a33adaae2a8e28b27117a9cdf8a980be33dbc9))

#### 🐛 Bug Fixes

- **db:** add unique constraints for spec_documents concurrency issues ([ff17cb1](https://github.com/specvital/infra/commit/ff17cb1aeb2a4e61b0a0bc1f5a19880e9b4f5a1f))

#### ⚡ Performance

- **db:** add index for spec reuse query ([189f85c](https://github.com/specvital/infra/commit/189f85c150700876d472f1a4f1a1867c7f051e6c))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **ci:** apply lint before creating schema docs PR ([aa00927](https://github.com/specvital/infra/commit/aa00927d496ddb28d2065dbf8909c14c0c25866a))

#### 📚 Documentation

- add specvital-specialist agent ([aa1dd4b](https://github.com/specvital/infra/commit/aa1dd4be5ddce0e6ab2dc488b45a05d1896dadbd))

#### ♻️ Refactoring

- rename schema-doc command to gen-schema-docs ([25e0463](https://github.com/specvital/infra/commit/25e04637b03130c264a93ba9e712de38442400c9))

#### 🔨 Chore

- sync ai-config-toolkit ([692d774](https://github.com/specvital/infra/commit/692d774da4da3175db4a55830e625295981d7642))
- sync docs ([bade4cf](https://github.com/specvital/infra/commit/bade4cf2cb45e6d901a70dfdd36decce770aecc9))
- sync-docs ([4cc4b51](https://github.com/specvital/infra/commit/4cc4b511de533c7ae675b5f30d1439a0a9475db7))
- **vscode:** add schema tools to Quick Command Buttons ([e6eb0a6](https://github.com/specvital/infra/commit/e6eb0a63fea3fb4d1656cd66ec13c7f05c5767dd))

## [web/v1.3.1](https://github.com/specvital/web/compare/v1.3.0...v1.3.1) (2026-01-07)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **ui:** fix AI Analysis button overflow on mobile ([e530add](https://github.com/specvital/web/commit/e530add96bde80f7ea3cb8bd1ae2e600f461fcbd))

### 🔧 Maintenance

#### ✅ Tests

- **spec-view:** skip handler tests during feature suspension ([df074c8](https://github.com/specvital/web/commit/df074c88bdbd2a626422a4fbf10ad091c7b9ad7a))

## [web/v1.3.0](https://github.com/specvital/web/compare/v1.2.1...v1.3.0) (2026-01-07)

### 🎯 Highlights

#### ✨ Features

- **ai-notice:** Implement AI feature suspension notice with Coming Soon modal ([8e33313](https://github.com/specvital/web/commit/8e33313d87902ab5482c077a5c00b5f4d84b162b))
- **api:** add OpenAPI spec for Spec View feature ([04c3f9c](https://github.com/specvital/web/commit/04c3f9c5836e2871619f82be422fccf475217383))
- **spec-view:** add cache freshness indicator and manual regeneration ([0f89732](https://github.com/specvital/web/commit/0f89732e449cc6f645ed6ef6f557af736f5114ef))
- **spec-view:** add spec view mode toggle with language selection dialog ([66151cf](https://github.com/specvital/web/commit/66151cfbd9f652a4da1d5b2cd869f2c9e2643fea))
- **spec-view:** implement cache repository for AI conversion results ([fc879a0](https://github.com/specvital/web/commit/fc879a0d0ae53916a65117cd1c7230a213aeaa9c))
- **spec-view:** implement domain layer for Spec View feature ([d55edf9](https://github.com/specvital/web/commit/d55edf9cf73667aee54886eebadea85731fd1d50))
- **spec-view:** implement frontend components and API integration ([eae6c12](https://github.com/specvital/web/commit/eae6c129bc4ea339be369855096be1dd7fe2b579))
- **spec-view:** implement Gemini AI Provider adapter ([6fdae64](https://github.com/specvital/web/commit/6fdae64e991e6d63c6abe099dab2976bda55e0ec))
- **spec-view:** implement usecase and handler for spec conversion API ([0f6fe4a](https://github.com/specvital/web/commit/0f6fe4a2bbb5b7419c4bbddc7f0a9681bc468263))
- **spec-view:** improve accessibility with keyboard navigation and ARIA support ([f6f7db4](https://github.com/specvital/web/commit/f6f7db4d332f72c5ff292620ff8ad44b08cacdc1))

#### 🐛 Bug Fixes

- **auth:** fix home to dashboard redirect failure on cold start ([a517b45](https://github.com/specvital/web/commit/a517b4596102c63c0c5205c168c04a8504c69c2d))
- **spec-view:** fix converted test names being mapped to wrong tests ([477ad26](https://github.com/specvital/web/commit/477ad265e50b8f700fae89e8e9e49e4017d63e2f))
- **spec-view:** fix spec view cache save failure ([d7f8f8b](https://github.com/specvital/web/commit/d7f8f8b75b46eb6549a78336b351cd0257d42965))
- **spec-view:** fix test name conversion failing for files after the first ([ce7d4e3](https://github.com/specvital/web/commit/ce7d4e356de72a428b8eef1e00a34d1138214159))

#### ⚡ Performance

- **spec-view:** fix AI conversion timeout for large repositories ([d507c91](https://github.com/specvital/web/commit/d507c91d14a1ce5248156f7d1c6bcbbb4f29e526))
- **spec-view:** improve AI prompt for better conversion quality ([ca9695a](https://github.com/specvital/web/commit/ca9695a60ca38b75d6a0c970b8d988039f68e879))
- **spec-view:** improve AI test name conversion quality (2nd iteration) ([cb02ad0](https://github.com/specvital/web/commit/cb02ad0fd2c7090ed0d4b0d062bc8e8a2271f775))

### 🔧 Maintenance

#### ♻️ Refactoring

- **spec-view:** extract AI prompt logic to separate file ([2f50536](https://github.com/specvital/web/commit/2f505361bf44c19c459f4e7ded9c14599837a21b))

#### 🔨 Chore

- sync docs ([4811f84](https://github.com/specvital/web/commit/4811f8461874ca92f189f450921c636ff4714d77))

## [web/v1.2.1](https://github.com/specvital/web/compare/v1.2.0...v1.2.1) (2026-01-04)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **auth:** fix dashboard redirect bug after time passes while logged in ([ef604fd](https://github.com/specvital/web/commit/ef604fdc41724ae42e62c9e57566712fb191d399))
- **routing:** fix 404 error for paths containing dots in repository names ([cfb36fc](https://github.com/specvital/web/commit/cfb36fc25ace5fff6303edc73ce72e6c4f3e8811))

### 🔧 Maintenance

#### 💄 Styles

- replace Korean comments with English ([796db0c](https://github.com/specvital/web/commit/796db0c5df3cd223a43a688c899e7e6fa9827e26))

## [web/v1.2.0](https://github.com/specvital/web/compare/v1.1.2...v1.2.0) (2026-01-04)

### 🎯 Highlights

#### ✨ Features

- **analysis:** add markdown export for analysis results ([9e63661](https://github.com/specvital/web/commit/9e63661695e0ca1c3b4a2b7e91e04e1d0f6eec3c))
- **analysis:** add status and framework filter functionality ([92a90e8](https://github.com/specvital/web/commit/92a90e83140b2dd8d7c446fffa8dcf139b03eb7a))
- **analysis:** add status mini-bar to test suite headers ([fa826d2](https://github.com/specvital/web/commit/fa826d2b8543a723183cfeb62570f60fd26bf4f8))
- **analysis:** add test search functionality ([688e383](https://github.com/specvital/web/commit/688e3833e24497753ea1d88ec7fb87353ec9ca32))
- **analysis:** add tree structure utilities for test suite navigation ([4f21ddf](https://github.com/specvital/web/commit/4f21ddfa8830375999d741f62a04f775efe4d1c4))
- **analysis:** add tree view UI with list/tree toggle ([2d1b3e5](https://github.com/specvital/web/commit/2d1b3e57f403fd93a1e6d4ad138dc763c83f4678))
- **analysis:** auto-track history on analysis page view ([ec23a75](https://github.com/specvital/web/commit/ec23a75585e296e43e23abbebbb612f73b99140c))
- **analysis:** improve accessibility for analysis page (WCAG compliance) ([c27d520](https://github.com/specvital/web/commit/c27d5203053a42d045e559702005b5707f0f909f))
- **analyze-dialog:** widen analyze repository modal ([1763ccb](https://github.com/specvital/web/commit/1763ccbf88549476b7d30e77585e77eb497b3577))
- **analyzer:** add cursor-based pagination to ListRepositoryCards usecase ([8f3f48d](https://github.com/specvital/web/commit/8f3f48d1f40e42e1d3c79f8e3f866fde432b5cc8))
- **analyzer:** add domain models and port interface for pagination ([9e8396a](https://github.com/specvital/web/commit/9e8396a97a9d400a535f07863ce9e721ebc81a9e))
- **analyzer:** add pagination support to GetRecentRepositories API ([eabf53b](https://github.com/specvital/web/commit/eabf53bd73911033e6857b848817056bc3cd9810))
- **analyzer:** add Repository feature for Dashboard ([f2b63b0](https://github.com/specvital/web/commit/f2b63b064464303752487c61fc23ba05d9b423e4))
- **analyzer:** implement keyset pagination SQL queries and adapter ([3887f69](https://github.com/specvital/web/commit/3887f6909afaffef22c4816b8b295ebd2835c50d))
- **api:** add domain-specific Dashboard API schema and types ([487f37b](https://github.com/specvital/web/commit/487f37b26c038a30748926707d79e705f8095c03))
- **api:** add others option to ownership filter ([f60f71f](https://github.com/specvital/web/commit/f60f71ffc6c4ba7610f1553e4ee7a05f91586cec))
- **api:** add testSummary field to RepositoryCard response ([f4791f3](https://github.com/specvital/web/commit/f4791f3a8f392086cf94ba36740a759ff7250c09))
- **api:** allow unauthenticated users to view Community tab data ([1e6e176](https://github.com/specvital/web/commit/1e6e1765bfc61da1bf41d91c2c5c2a14deeecd30))
- **api:** implement add analysis to dashboard API ([661a906](https://github.com/specvital/web/commit/661a90677f88de2e24732d457975d700f6d5b66d))
- **auth:** add bookmark CRUD for Dashboard feature ([d7a4611](https://github.com/specvital/web/commit/d7a46111123e2317238d6e97a55728c3f46c8cb2))
- **auth:** add dev-only test login endpoint ([5fc5a75](https://github.com/specvital/web/commit/5fc5a75264062c01f9e870fb961e9e239b39d80a))
- **auth:** add premium CTA styling to LoginButton ([a70a958](https://github.com/specvital/web/commit/a70a95898f5eb8b7e3279c1cae0b99ec87e18a43))
- **auth:** allow unauthenticated users to access Explore page ([66cea9c](https://github.com/specvital/web/commit/66cea9c4c4fb767143f51eb851f7359277b85caa))
- **auth:** implement Access + Refresh token pair handling in AuthHandler/Middleware ([12625bc](https://github.com/specvital/web/commit/12625bc9c5673081bcd7274e6ec33637b3f511e7))
- **auth:** implement access/refresh token generation logic ([9f27e0d](https://github.com/specvital/web/commit/9f27e0d178c9685e67f7dbfe07b53fbebe1b5251))
- **auth:** implement refresh token infrastructure ([5ddcced](https://github.com/specvital/web/commit/5ddcced38da76373f3a4d569bbf7201f9b9e004f))
- **auth:** introduce login modal for provider selection ([e69282d](https://github.com/specvital/web/commit/e69282da4aa85c56517877fbb877421c8b62f9b8))
- **backend:** implement GitHub App installation status and install URL APIs ([1e37193](https://github.com/specvital/web/commit/1e37193f05fa09ebabed1291c5d13de3b73bc3c1))
- **backend:** implement GitHub App module backend foundation ([4496014](https://github.com/specvital/web/commit/44960143523e2ec21d30dbbd72eb08cc4c82766c))
- **backend:** implement GitHub App webhook handler for installation events ([944f5ec](https://github.com/specvital/web/commit/944f5ec4838f27c031ce77d7a82adfcde129f150))
- **backend:** integrate GitHub App token for organization repository access ([008e58a](https://github.com/specvital/web/commit/008e58a1368df74be4e61559bdde36583a71cc57))
- **button:** add cta variant for premium call-to-action buttons ([5552958](https://github.com/specvital/web/commit/555295818c80b0c62a726a2f6b411134bdff6740))
- **button:** add header icon button variant and size ([1df89f3](https://github.com/specvital/web/commit/1df89f3719d088dcf77c319552a71860b4b76100))
- **dashboard:** add Attention Zone highlighting repos needing action ([e8502f2](https://github.com/specvital/web/commit/e8502f267cc0931cccdef8f48d719be8820b7faf))
- **dashboard:** add backend API support for ownership filter ([e30c460](https://github.com/specvital/web/commit/e30c46099441c00bdd3f3f82cf46eda5b62ec5a9))
- **dashboard:** add card grid improvements and empty state variants ([e6a040d](https://github.com/specvital/web/commit/e6a040d995c62fe5fc25633a44e0a6ab87174c0b))
- **dashboard:** add collapsible Discovery section at bottom ([4b581ad](https://github.com/specvital/web/commit/4b581add6467687d000978f09ac0eb50ab6fe093))
- **dashboard:** add dashboard route with auth-based redirects ([e9b36f9](https://github.com/specvital/web/commit/e9b36f9c5900ea18732fd74bbc85851f1f52cc4d))
- **dashboard:** add DashboardHeader, DashboardContent components ([ca85b87](https://github.com/specvital/web/commit/ca85b8763fd20c57fdedf47a01217451cf1abc10))
- **dashboard:** add discovery section for unanalyzed GitHub repositories ([f527e5a](https://github.com/specvital/web/commit/f527e5a37288b6f4e65c0b1a55bdc15cc1a32d6d))
- **dashboard:** add error handling, toast notifications, and keyboard accessibility ([35c868e](https://github.com/specvital/web/commit/35c868e0b46d73a1aeff27889ef1f0fe011ca533))
- **dashboard:** add frontend API client and types ([c155458](https://github.com/specvital/web/commit/c155458058ce8ab12bc2a4206274b8a4ac8b98fb))
- **dashboard:** add i18n messages for dashboard components ([365cec6](https://github.com/specvital/web/commit/365cec6805f1772322e449d725b785e823882a7b))
- **dashboard:** add immediate visual feedback on Update button click ([cfbc663](https://github.com/specvital/web/commit/cfbc66366d53f88f57f28de44ea40ff0467ceb95))
- **dashboard:** add mobile responsive polish with bottom nav and filter drawer ([2bbc8f4](https://github.com/specvital/web/commit/2bbc8f48b47b2c7bc4b7416a97e555c099b988b8))
- **dashboard:** add modal-based new analysis from dashboard ([b6dc302](https://github.com/specvital/web/commit/b6dc302435032afb6caa26127f9eefd996500190))
- **dashboard:** add my analyses tab with ownership filter ([62d5b03](https://github.com/specvital/web/commit/62d5b0337ff045193048c6db0013e1dfd15adfc8))
- **dashboard:** add organization repository discovery flow ([cbf37b5](https://github.com/specvital/web/commit/cbf37b569c6e59323be6940eb7090dcdbce0b068))
- **dashboard:** add others filter and set default view to my ([3e6a08c](https://github.com/specvital/web/commit/3e6a08c4c3ea1c3af224eceeaa4a3b1cd8df03cd))
- **dashboard:** add React Query hooks for Dashboard feature ([ec9507f](https://github.com/specvital/web/commit/ec9507fc85f78bef782308be191bbc26dde6de86))
- **dashboard:** add real UpdateStatus calculation for repository cards ([26671fb](https://github.com/specvital/web/commit/26671fb5e9c2b3c6647bee1630cb4b88a16a35ed))
- **dashboard:** add RepositoryCard and RepositorySkeleton components ([93f0189](https://github.com/specvital/web/commit/93f0189c7120541f710932010b4d1243b4fdd13a))
- **dashboard:** add RepositoryList, BookmarkedSection, EmptyState components ([dc21547](https://github.com/specvital/web/commit/dc21547bc802db2f5397b5d691fc2ba3847c5747))
- **dashboard:** add scroll-area component for horizontal scroll UI ([bd2bc42](https://github.com/specvital/web/commit/bd2bc42e1e1d53c6648287279b74db05ecd60b2e))
- **dashboard:** add starred filter toggle for bookmarked repositories ([3219362](https://github.com/specvital/web/commit/3219362686fecae7b85154bae96233bf6e87ca1f))
- **dashboard:** add Summary Section with animated stats cards ([1cc2b31](https://github.com/specvital/web/commit/1cc2b316826056b27c3d5f1730d69ce007d7140c))
- **dashboard:** add tab infrastructure with nuqs URL state management ([93c512a](https://github.com/specvital/web/commit/93c512a9114726ca198f166ca02ccc78f0a4e792))
- **dashboard:** add TestDeltaBadge and UpdateStatusBadge components ([6f99147](https://github.com/specvital/web/commit/6f99147cfb00b1f3402c9c1ec552343962aa49c9))
- **dashboard:** add unified View filter to restore my/community analysis distinction ([3d3d95c](https://github.com/specvital/web/commit/3d3d95cdde799057f83e7640a67741d78c7b6707))
- **dashboard:** add useInfiniteQuery-based pagination hooks ([fb96ae7](https://github.com/specvital/web/commit/fb96ae7ce050adceb919adb2b934c58ce92a257e))
- **dashboard:** align repository card header vertically ([16f1f19](https://github.com/specvital/web/commit/16f1f192249c15f7218b079c2cd246c994d7aa31))
- **dashboard:** apply gradient icons and depth styling to StatCard ([2506892](https://github.com/specvital/web/commit/25068928294daee502221038940677620fc4b1ef))
- **dashboard:** consolidate tabs into single repository list ([8402f84](https://github.com/specvital/web/commit/8402f84ae0953f5efc8e6a13435a2bf9ce4ea9e0))
- **dashboard:** improve AttentionZone card styling and dark mode accessibility ([25c7722](https://github.com/specvital/web/commit/25c7722a75a4693252b06a278335ee8558b369f4))
- **dashboard:** improve mobile responsive layout for search/sort controls ([48c452d](https://github.com/specvital/web/commit/48c452d212dd40026d0bf5a964760edab026b101))
- **dashboard:** integrate Load More button and pagination UI ([3756946](https://github.com/specvital/web/commit/37569463f6f20df048765a96f9ac2b122cb0478c))
- **dashboard:** replace Load More button with infinite scroll ([0a6b394](https://github.com/specvital/web/commit/0a6b3944ba0411bfe2f6c2b72abff4bbb8867f0a))
- **db:** add SQL queries for Dashboard feature ([10fe653](https://github.com/specvital/web/commit/10fe653906e65a363fafee5da1fe163092b2b70b))
- display describe block names (suiteName) in test suites ([9e6a5ab](https://github.com/specvital/web/commit/9e6a5ab11fd90f698f75e2a6528494211c411432))
- **dropdown:** add glassmorphism effect to dropdown menu ([e51b4d8](https://github.com/specvital/web/commit/e51b4d89d146f0a83d59e26e36809b6c3217f2b9))
- **explore:** add Explore page and GNB navigation tabs ([d1a7b7a](https://github.com/specvital/web/commit/d1a7b7a678826b91710b4d31097b3a51e1148f97))
- **explore:** add login prompt UI for auth-required tabs ([a5b8dc2](https://github.com/specvital/web/commit/a5b8dc24601d5c552530bf8e5425da9bf75fdca8))
- **explore:** filter private repositories from Community tab ([f25e80e](https://github.com/specvital/web/commit/f25e80ef7263be5f194bbe4123ac9b6ac2200270))
- **explore:** implement "Add to Dashboard" feature in Community tab ([5e5e6cf](https://github.com/specvital/web/commit/5e5e6cf54bd47fb10382808ed2d53c687d2daa9a))
- **explore:** improve Add to Dashboard button UX ([f75a53d](https://github.com/specvital/web/commit/f75a53dc121dea40564dfe330b4c14d40a73be3a))
- **explore:** integrate Discovery tabs for repo exploration ([ec36872](https://github.com/specvital/web/commit/ec368723535162383a8d462b63376aaa389e0dec))
- **explore:** show login modal when unauthenticated users click bookmark ([403f616](https://github.com/specvital/web/commit/403f6167b91d23f902f8d6b29a725f3b641721a8))
- **frontend:** add context-aware input feedback and help tooltip for URL input ([05f61ee](https://github.com/specvital/web/commit/05f61ee98550b2d26bd19da295cd087f3df42763))
- **frontend:** add GitHub icon prefix and improve mobile touch targets ([0cdb871](https://github.com/specvital/web/commit/0cdb8711fa2a04cd66ed3874c3092c9f09bf960d))
- **frontend:** add metadataBase for OG image support ([9ed0417](https://github.com/specvital/web/commit/9ed04175b6f85337a6447a770371fa419296dd0c))
- **frontend:** add real-time URL validation and improve error UI ([3612ec5](https://github.com/specvital/web/commit/3612ec51facd8b5e92b2629e52b097a87bf2287f))
- **frontend:** add tagline to Home page for clearer value proposition ([042f542](https://github.com/specvital/web/commit/042f5423e090f1182b235bbb4f229b643b327c71))
- **frontend:** implement auto token refresh with request retry on 401 ([a19e50a](https://github.com/specvital/web/commit/a19e50acd378747667dcf4fdcc451bece1f806ee))
- **frontend:** implement GitHub App organization connection UI ([140eb16](https://github.com/specvital/web/commit/140eb1664e1a6ef0e25a91d7342e675ac07e6a90))
- **frontend:** support various GitHub URL input formats ([f276b15](https://github.com/specvital/web/commit/f276b15395a750171ad4b04c0781046bf3857387))
- **frontend:** wrap form in Card component and add Trust Badges ([eea8a5b](https://github.com/specvital/web/commit/eea8a5bb60eec3e100c9bc457c867f6a94cdd013))
- **github:** add user GitHub repositories and organizations API ([a4db76e](https://github.com/specvital/web/commit/a4db76e8cf2dca28166370ccad75045b362d7415))
- **header:** add accessibility improvements for WCAG compliance ([3200e23](https://github.com/specvital/web/commit/3200e231233f9145ae612ec0b4f22f8072362dd1))
- **header:** add icon button grouping with visual hierarchy ([66c372b](https://github.com/specvital/web/commit/66c372b7fb5620364472b8d68fbb58b79b3a60ab))
- **header:** add logo image to header ([9a2bb37](https://github.com/specvital/web/commit/9a2bb37d05be6aca30be9e9a935c3a63c3a9885f))
- **header:** add rotation animation to theme toggle button ([6004d1e](https://github.com/specvital/web/commit/6004d1e1a0f0d6ae00679e19560bb5c9f85745d4))
- **header:** add scroll shadow effect for visual separation ([393ef33](https://github.com/specvital/web/commit/393ef330b94f4f8e39a0640f4c629baee9b26d47))
- **header:** add tooltips to icon buttons and improve animation ([eb2a9b7](https://github.com/specvital/web/commit/eb2a9b7a69908ae79fb23ae04de813925d45451a))
- **header:** enhance glassmorphism effect and add gradient borders ([aebb087](https://github.com/specvital/web/commit/aebb087c35b19bfce8e6399802d1146ab18afa3f))
- **header:** improve header border styling ([7a78ba2](https://github.com/specvital/web/commit/7a78ba23bfc6b2a3f675a067442b9c1ec48ee90c))
- **home:** add page load animation using Motion library ([e6d4350](https://github.com/specvital/web/commit/e6d4350e87d5ab4598f32426d207cfb4a4eb2604))
- **home:** restyle trust badges and add supported frameworks dialog ([02bb7bb](https://github.com/specvital/web/commit/02bb7bb50912a00845927366c4c616bf7d91dd2c))
- **mobile-nav:** hide Dashboard link for unauthenticated users ([4262f60](https://github.com/specvital/web/commit/4262f6060432bcc91e7d7c761d1b28a43a821289))
- **nav:** show navigation tabs on homepage with auth-based filtering ([176c812](https://github.com/specvital/web/commit/176c8123aa6741b1fdb93871dd3e21336e67489d))
- **overlay:** add glassmorphism effect to dialog, popover, and tooltip ([50aa14c](https://github.com/specvital/web/commit/50aa14c4f9be397d378648d9b031b5d8d427f822))
- **query:** auto-refresh dashboard cache on analysis completion ([bb833ee](https://github.com/specvital/web/commit/bb833ee38d4fc484846e7a8c5738e675ce1115cf))
- **style:** simplify homepage messaging and layout ([b97fffd](https://github.com/specvital/web/commit/b97fffd39bccec536645e6fed7746ad004897fba))
- **ui:** add AST-powered trust badge and improve layout ([8ccdd4d](https://github.com/specvital/web/commit/8ccdd4da069a4afcd37e56d2751fe0f8c51bc3c9))
- **ui:** add branch name and commit time display to analysis page ([9d0fe50](https://github.com/specvital/web/commit/9d0fe5099ec3b941277e97258bead72efd9c2d9d))
- **ui:** add Card depth and hero gradient for premium look ([8d95215](https://github.com/specvital/web/commit/8d9521507f9c7e12dadc64fec2c84588aa1cb3e5))
- **ui:** add Card depth and micro-interaction feedback ([8660e66](https://github.com/specvital/web/commit/8660e66ea0e44f0317f1e9ca3112f32d9b0173cf))
- **ui:** add mobile bottom navigation bar ([7c8a6c5](https://github.com/specvital/web/commit/7c8a6c519930deefae0eaeb8099f28fb5ead21bc))
- **ui:** add Motion animations to analysis page ([392bfb2](https://github.com/specvital/web/commit/392bfb22f5b09b3bfc5712e7c2232570a492ef8c))
- **ui:** add ResponsiveTooltip component for mobile touch support ([f1f35d4](https://github.com/specvital/web/commit/f1f35d427d76668325dedc407ef1439828d1e247))
- **ui:** add slide-up animation and polish to mobile bottom bar ([c48d79b](https://github.com/specvital/web/commit/c48d79bae6211136943c216da6f281a451c64421))
- **ui:** apply Cool Slate color palette ([04ab37d](https://github.com/specvital/web/commit/04ab37d208e3c37a70c23b5c0dbfa256e92b6fd1))
- **ui:** apply Pretendard font for Korean typography ([6617845](https://github.com/specvital/web/commit/6617845d2c822e02f73ee15a5ff7adf72cbae58c))
- **ui:** fix ViewModeToggle to right edge on mobile ([f6faf34](https://github.com/specvital/web/commit/f6faf34952670112391fc257aa52b04a9af7e1b0))
- **ui:** improve mobile bottom bar visual distinction ([1969f77](https://github.com/specvital/web/commit/1969f7755f707cc6184977ac2ca14b8714e1dbe5))
- **ui:** improve stats card visual hierarchy and design ([20348d8](https://github.com/specvital/web/commit/20348d8226007bc18d5fe4ff8a791d0804a011e3))
- **ui:** improve typography and border visual refinement ([06ee2c7](https://github.com/specvital/web/commit/06ee2c70fc2f66e8cb2742b860677de84a5fd6b5))
- **ui:** optimize mobile filter bar touch targets and scroll ([c14584f](https://github.com/specvital/web/commit/c14584fe405c8ccbbc58df12631909f57fcbc67d))
- **ui:** reduce Card elevated hover lift effect ([b306bff](https://github.com/specvital/web/commit/b306bff84ec9f89917ec2db822c5922f4bc3d538))
- **ui:** unify analysis page header button styles and show Share text on mobile ([ab3bcf0](https://github.com/specvital/web/commit/ab3bcf09bf28d1ac1de9494917069aa690d42c86))
- **ui:** unify analysis page header button styling and improve mobile layout ([6c27048](https://github.com/specvital/web/commit/6c270483fc9319bb700bd866ab5d9cbdcd46d7f0))
- **ui:** unify filter button heights to h-9 ([9576fa2](https://github.com/specvital/web/commit/9576fa26f14662273db92129e86e5c118095e09f))
- **user:** add user analyzed repositories API ([acc56a7](https://github.com/specvital/web/commit/acc56a767e53549f2321d75dea33bef0d125e7df))

#### 🐛 Bug Fixes

- **analyze-dialog:** close modal when navigating to analysis page ([4515a65](https://github.com/specvital/web/commit/4515a65d827ad461aec4faa1352af031d6f5d9b5))
- **auth:** add missing refresh_token cookie in OAuth callback and middleware token refresh ([4f2e6d5](https://github.com/specvital/web/commit/4f2e6d502c753876dc6c81a08b7a7d796ba3f30f))
- **auth:** fix auto-logout after access token expiry due to skipped refresh ([159ea1d](https://github.com/specvital/web/commit/159ea1d0bf87d1d7625aa88906250eea9be1adcd))
- **auth:** fix incorrect redirect to homepage after token refresh ([d60b0a6](https://github.com/specvital/web/commit/d60b0a6c7240fff16fd310e5495c57a2950322f2))
- **auth:** fix intermittent home redirect on dashboard refresh and unauthenticated 401 errors ([6994a6b](https://github.com/specvital/web/commit/6994a6b48c0021adeff3ae57af267841aba7c7a0))
- **dashboard:** bookmark toggle not reflecting UI state changes ([6e37766](https://github.com/specvital/web/commit/6e3776625c39fe727da3ae941f617a4b599ba171))
- **dashboard:** integrate Tab UI and fix ownership filter SQL bug ([74444c1](https://github.com/specvital/web/commit/74444c1db2556de56415e5593b03370ccc0c759c))
- **dashboard:** maintain consistent repository card height regardless of Update button ([62d5ef0](https://github.com/specvital/web/commit/62d5ef0a8564642a2d97766cb2e8c78d9a242091))
- **dashboard:** scope summary stats to user's analyzed repositories ([18b18b4](https://github.com/specvital/web/commit/18b18b435cb3b8cb79e2f0bfb0d662817c2e9451))
- **explore:** display repository fullName in Organizations tab ([caf2b91](https://github.com/specvital/web/commit/caf2b910a6b295ca465e8d5510f1a9b1047f11fd))
- **frontend:** add missing padding to discovery sheet content area ([ca8a963](https://github.com/specvital/web/commit/ca8a9630a86c6bfea84d411be8c882bdd0c7c6c0))
- **frontend:** improve card visual consistency in drawer ([ce617b0](https://github.com/specvital/web/commit/ce617b05e36cfb676443288ae543c1c104383c5e))
- **frontend:** resolve infinite redirect loop caused by invalid token after DB reset ([d44975f](https://github.com/specvital/web/commit/d44975f89f1c5db07ef3d1f09d64fb65e83db1ce))
- **frontend:** resolve relativeTime hydration mismatch warning ([3220325](https://github.com/specvital/web/commit/3220325943aca57714167a15880247eb17884dcd))
- **frontend:** resolve relativeTime hydration mismatch warning ([24bdcd5](https://github.com/specvital/web/commit/24bdcd5701b86f74fcef9cfdaff327fada4af893))
- **i18n:** add missing dashboard.title translation key ([d80b571](https://github.com/specvital/web/commit/d80b571ee8ce8e016887ce46cc0e7dddb3488546))
- organization repos not showing for OAuth-restricted orgs ([14e3b8a](https://github.com/specvital/web/commit/14e3b8ae8175af648f13a78c25eb9bc64c43111e))
- **ui:** add cursor-pointer consistency to interactive components ([3b840d6](https://github.com/specvital/web/commit/3b840d65cea76c1cbab2c04fba47e183b46d6204))
- **ui:** align Input and Analyze button heights ([554a3d7](https://github.com/specvital/web/commit/554a3d7bd992ec863fc7f29ad893959649b625bf))
- **ui:** align mobile bottom navigation styling with design system ([1941ccc](https://github.com/specvital/web/commit/1941cccc934ac2a0be78264a6b364561dc15641b))
- **ui:** fix homepage unwanted scroll and content centering issue ([6f6807b](https://github.com/specvital/web/commit/6f6807b2e57eb329e6831b6403370fc49825a97e))
- **ui:** improve URL input form layout and placeholder on mobile ([b535cd4](https://github.com/specvital/web/commit/b535cd41240b21080d2536328f4398b73da9a8d5))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **ci:** pin oapi-codegen version to fix CI build failure ([9ffdedb](https://github.com/specvital/web/commit/9ffdedb6c9105d352407d9258a87fbb4481a425e))
- **devcontainer:** fix network creation failure in Codespaces ([1d33990](https://github.com/specvital/web/commit/1d3399098e836b092912e7a67b92db23b752e0c9))
- update missing lock files ([e53e325](https://github.com/specvital/web/commit/e53e3258f13742251ac072c5374b2e8df63d4444))

#### 📚 Documentation

- **backend:** update CLAUDE.md to reflect Clean Architecture migration ([2dbfbd5](https://github.com/specvital/web/commit/2dbfbd57ef5aa198b12a5405214ad0ca183ce8fe))
- document GitHub App environment variables ([26911af](https://github.com/specvital/web/commit/26911af494eff00302efcd263516bf91234508b5))
- fix feedback links in README ([63f84b6](https://github.com/specvital/web/commit/63f84b674dc1ed6ffd7381d9192fe09d8ca9f129))
- update and add specific CLAUDE.md ([1769def](https://github.com/specvital/web/commit/1769def3e4f625c1e7e110b3a515cb2fb37f9433))
- update README.md ([b58d172](https://github.com/specvital/web/commit/b58d1722df901e2eb6cdf424b3fdbf2e8492dc76))

#### ♻️ Refactoring

- **analysis:** unify framework color system with name-based lookup ([63b7430](https://github.com/specvital/web/commit/63b7430e7aafeb6ae3eefd8f69c439e82cefaa52))
- **auth:** improve protected route management with PROTECTED_ROUTES array ([4125d9c](https://github.com/specvital/web/commit/4125d9c9a37e0b116956cc10fc995f73a197bdcb))
- **auth:** remove unused TokenVersion infrastructure and update auth docs ([7a853df](https://github.com/specvital/web/commit/7a853dfcf9ff273104020cf02611ed92a638a139))
- **auth:** simplify LoginButton to use cta variant and add reduced motion support ([0d256bf](https://github.com/specvital/web/commit/0d256bf53ee50a466d3ec3c31bd009f7f6c2638b))
- **backend:** add compile-time interface checks to github adapter ([f98e9f0](https://github.com/specvital/web/commit/f98e9f0da05546419ace6b9481953e8ead5d6b0f))
- **backend:** align auth module structure with Clean Architecture pattern ([568e944](https://github.com/specvital/web/commit/568e944f35f3760cacce5be537a940f87eeee3c1))
- **backend:** apply Clean Architecture adapter layer to analyzer module ([1bb5eee](https://github.com/specvital/web/commit/1bb5eee826dfbc9f2070e518b1def80843367486))
- **backend:** apply Clean Architecture domain layer structure to analyzer module ([66f325f](https://github.com/specvital/web/commit/66f325fcfa963ed305ec41819bf6b79d3242d9e6))
- **backend:** apply Clean Architecture handler layer to analyzer module ([d07cadc](https://github.com/specvital/web/commit/d07cadc7f18ce42855e95022b03548fa0b5c69a6))
- **backend:** apply Clean Architecture usecase layer to analyzer module ([496459f](https://github.com/specvital/web/commit/496459f91ec1257ce9f2669c7ff888b3164a8863))
- **backend:** apply full Clean Architecture to github module ([d5ba9ad](https://github.com/specvital/web/commit/d5ba9ad66f0661e231f5d215412c602f0f01c6dc))
- **backend:** apply full Clean Architecture to user module ([8c7c0e7](https://github.com/specvital/web/commit/8c7c0e7e3286a6b63835dfd4e47d4fb482cd024e))
- **backend:** fully apply Clean Architecture to auth module (Handler → Usecase direct connection) ([0872a32](https://github.com/specvital/web/commit/0872a3212e1a104f39e5ed9213975cc4106cab1e))
- **backend:** move bookmark module from auth to user ([ddcb32f](https://github.com/specvital/web/commit/ddcb32fcd58e387b2bd390080d046c560d43e484))
- **backend:** simplify code by making GitHub App config required ([b8ce528](https://github.com/specvital/web/commit/b8ce5284b0c7ba10c23094dfb920589749e48864))
- **backend:** unify auth module mapper location under adapter ([12d176d](https://github.com/specvital/web/commit/12d176d018e89750f6e49fa56f81b861b5097118))
- **backend:** unify user module errors.go location under domain package ([5eec2c6](https://github.com/specvital/web/commit/5eec2c64912f937eb755364e62ee81a3a1bb0982))
- **dashboard:** clean up dead code after dashboard simplification ([b99cfe8](https://github.com/specvital/web/commit/b99cfe84b03371e7e7da34ee45d3ed633a2ecef0))
- **dashboard:** clean up pagination legacy code and add optimizations ([f717ce2](https://github.com/specvital/web/commit/f717ce2810bf5ef13484977d41aba264a9bb57e5))
- **dashboard:** improve dashboard information hierarchy ([a87a699](https://github.com/specvital/web/commit/a87a69992a676ba6bf49c72a08aaa9f6152f705f))
- **dashboard:** redesign filter bar with 2-axis filter system ([c57496b](https://github.com/specvital/web/commit/c57496b2ddf0373b124af6650cd090d28dc4d75a))
- **dashboard:** remove AttentionZone section ([c942e89](https://github.com/specvital/web/commit/c942e89fe130c46e74bf8bd884c940fbcc1ed0e2))
- **dashboard:** remove community filter and add Explore CTA ([281a943](https://github.com/specvital/web/commit/281a9430ba1434c01d5be8e2213bf21ce152ef22))
- **dashboard:** remove duplicate bookmark icon from repository card ([8a11fc0](https://github.com/specvital/web/commit/8a11fc0a5997baca89921c13098bbb743244a33f))
- **dashboard:** replace BookmarkedSection with inline bookmark indicators ([c505bb1](https://github.com/specvital/web/commit/c505bb1a81dfda63d9e24558a76c873619886cb9))
- **docs:** change API docs path to /api/docs ([7ff1f42](https://github.com/specvital/web/commit/7ff1f42db5fcbba79f639217cf6e5285a25dfe9d))
- **explore:** remove redundant "Add to Dashboard" button from cards ([9a869c2](https://github.com/specvital/web/commit/9a869c230a20bc10711b0f9c6a2097eb55c3fc04))
- **ui:** apply cta variant to all primary action buttons ([2978b17](https://github.com/specvital/web/commit/2978b176adb77c988de85043135d64f1ef9d86b2))
- **ui:** centralize glassmorphism and shadow system with CSS variables ([5e6caaf](https://github.com/specvital/web/commit/5e6caafa377da4d7d1e8349ef0a45d89658a93b9))
- **ui:** extract AnalysisHeader component and add debounce cleanup ([68d4223](https://github.com/specvital/web/commit/68d4223728aefca8a08a7bc9a581832448c3c296))
- **ui:** improve AnalysisHeader information grouping structure ([bb183f0](https://github.com/specvital/web/commit/bb183f0186f028d26cd5bf3a83f6e86174b31d77))

#### 🔨 Chore

- \*.tsbuildinfo gitignore ([5236b43](https://github.com/specvital/web/commit/5236b43095dd9c311b6b1b407b2cdbab3df430ff))
- add development environment setup for GitHub App integration ([9dc482e](https://github.com/specvital/web/commit/9dc482e37eb96c54bbccab0cef52b712635ec6cb))
- add just smee command ([68b9b73](https://github.com/specvital/web/commit/68b9b73bab870c74f64f7db84384486c1e320445))
- add migrate-local action command ([9cf475e](https://github.com/specvital/web/commit/9cf475e817c18001589ca34b5e24badbfad99aa9))
- **api:** remove deprecated marker from ViewFilterParam ([5a96793](https://github.com/specvital/web/commit/5a967930578a2598386f9734b7ac2f08b4d3cbd3))
- dead code cleanup ([12eac95](https://github.com/specvital/web/commit/12eac95a70fb3b1d66586ea9be13e0032083144c))
- delete commit message markdown ([a034af5](https://github.com/specvital/web/commit/a034af58bc9750bec42fd8f02aaee9f1ce375b4c))
- **deps-dev:** bump @semantic-release/release-notes-generator ([b17f20b](https://github.com/specvital/web/commit/b17f20bb21d09a323bb90a65ffe4e1da4af19b39))
- **deps:** bump actions/setup-go from 5 to 6 ([990bab2](https://github.com/specvital/web/commit/990bab24df366c2afe7fec87decce404ada1aba0))
- dump schema ([0012dd0](https://github.com/specvital/web/commit/0012dd08d3a6ffd0559ab0af01098665cbc709bc))
- dump schema ([d92dc7b](https://github.com/specvital/web/commit/d92dc7b502cb99182dd458f4f5fb79b2f57a4264))
- dump schema ([ff21aa3](https://github.com/specvital/web/commit/ff21aa3f403b4d9795cfaa477a4d3fbdbeccb399))
- dump schema ([90d1816](https://github.com/specvital/web/commit/90d181663ecf0b6969df8f17dee00ad538bafe23))
- fix vscode import area not automatically collapsing ([428d3cf](https://github.com/specvital/web/commit/428d3cf9c41f875fd4f720a8d175fb92d24cdaf4))
- **frontend:** add ESLint flat config and improve lint commands ([7e096cf](https://github.com/specvital/web/commit/7e096cf5a23cfa4a6cba4b4401cc89cb9b93c622))
- **frontend:** add shadcn/ui components for Home UI/UX improvements ([75c659c](https://github.com/specvital/web/commit/75c659c0c211d93458987b241b9a1e0ed1824744))
- improved the claude code status line to display the correct context window size. ([1441d70](https://github.com/specvital/web/commit/1441d70b3d90b90faba2fe070a69b16ab50bf500))
- **motion:** install motion library and setup animation utilities ([31ed5c6](https://github.com/specvital/web/commit/31ed5c603422f8bd90b38dbbb5430ee0dfcaace3))
- sync docs ([3ce9578](https://github.com/specvital/web/commit/3ce9578105731593400143a299f6d2a402a2b85e))
- sync docs ([71f79a8](https://github.com/specvital/web/commit/71f79a83e6bc2b7b5c4100e7d32f02c0b26190cc))
- **ui:** add safe area infrastructure for mobile bottom navigation ([ed8ebb8](https://github.com/specvital/web/commit/ed8ebb8bf2b6a79464e284ad84a306cc36f63269))

## [worker/v1.1.1](https://github.com/specvital/worker/compare/v1.1.0...v1.1.1) (2026-01-04)

### 🔧 Maintenance

#### 🔨 Chore

- update core ([1f224fd](https://github.com/specvital/worker/commit/1f224fddf200ebb72e9fe11acfc114d988e21fba))

## [worker/v1.1.0](https://github.com/specvital/worker/compare/v1.0.6...v1.1.0) (2026-01-04)

### 🎯 Highlights

#### ✨ Features

- add Clone-Rename race condition detection ([ebbf443](https://github.com/specvital/worker/commit/ebbf443e3b9271c4dba82d38a567d3efdc0236a9))
- add codebase lookup queries based on external_repo_id ([d6e0b79](https://github.com/specvital/worker/commit/d6e0b797ec46aab29eb897c6c1d6a8cdbb6b47c6))
- add codebase stale handling queries and repository methods ([939e078](https://github.com/specvital/worker/commit/939e07886a8a26c385ad3206c71fe8205a6bd001))
- add GitHub API client for repository ID lookup ([46e40b8](https://github.com/specvital/worker/commit/46e40b806843f7e87ec98269baca7bce136064bd))
- determine repository visibility via reversed git ls-remote order ([0bc988e](https://github.com/specvital/worker/commit/0bc988e839f795678281ac6431b41199e4f95f95))
- integrate codebase resolution case branching into AnalyzeUseCase ([0f58440](https://github.com/specvital/worker/commit/0f58440f0012c15fd215f57a58917370ff93b2a9))
- record user analysis history on analysis completion ([e2b2095](https://github.com/specvital/worker/commit/e2b2095c47dffbd51fa57d3d24c550c19cfed851))
- store commit timestamp on analysis completion ([24bdbd7](https://github.com/specvital/worker/commit/24bdbd7050a40fbcf41965e6abfb728dc9460870))

#### 🐛 Bug Fixes

- add missing focused and xfail TestStatus types ([b24ee33](https://github.com/specvital/worker/commit/b24ee333e5c5e0f71098326f934c09853976fee6))
- add missing is_private column to test schema ([5744b95](https://github.com/specvital/worker/commit/5744b956a7b408f3dbc7583456eaedbf7fa1f4f6))
- ensure transaction atomicity for multi-step DB operations ([16834ef](https://github.com/specvital/worker/commit/16834ef0837917df4c30d31135e4f97a8a07eb3b))
- exclude stale codebases from Scheduler auto-refresh ([933c417](https://github.com/specvital/worker/commit/933c41711f375979d96cf5401ba93e6171891b49))
- fix visibility not being updated on reanalysis ([2424a5f](https://github.com/specvital/worker/commit/2424a5fadaebfd0fed1aba07045f2c86ddd5c585))
- prevent duplicate analysis job enqueue for same commit ([1a996ea](https://github.com/specvital/worker/commit/1a996ea38ad6742647317932e7acbb24939146e1))
- prevent unnecessary job retries on duplicate key error ([40eda32](https://github.com/specvital/worker/commit/40eda32b890206f3f3ef5913ce8ed4f9afdc0cdb))
- resolve stray containers from failed testcontainers cleanup ([1ef5124](https://github.com/specvital/worker/commit/1ef5124a617fcbc1ddd434b6a74baa6dd5ab390a))

#### ⚡ Performance

- improve DB save performance for large repository analysis ([200a527](https://github.com/specvital/worker/commit/200a5275cf639a2c0c65d955e79dbe65ad4f7068))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **devcontainer:** fix network creation failure in Codespaces ([2054227](https://github.com/specvital/worker/commit/2054227927b13127fb2c770323dcc17e6bba4d0a))
- isolate git ls-remote environment to fix private repo misclassification ([7d15fb8](https://github.com/specvital/worker/commit/7d15fb82534cb2c4c34ea368173265c185abf543))

#### 📚 Documentation

- add CLAUDE.md ([5194d71](https://github.com/specvital/worker/commit/5194d713b2f07fd2d4d2a66df62f861520b027bc))
- add missing version headers and improve CHANGELOG hierarchy ([d6436ab](https://github.com/specvital/worker/commit/d6436ab60b12e4bf551c23d59009fa66782e6eb4))
- rename infra repo in docs ([1bdb806](https://github.com/specvital/worker/commit/1bdb806dabc7fd082cb114e93f349aaa619d5315))

#### 💄 Styles

- format code ([8616fbd](https://github.com/specvital/worker/commit/8616fbdae4105860c87569093f302ba6a877c6c7))

#### ♻️ Refactoring

- remove unused deprecated Stop method ([c034ecc](https://github.com/specvital/worker/commit/c034ecc56660bda965a297072e7d23400e8b8e61))
- **test:** auto-sync test schema with production schema ([77668e0](https://github.com/specvital/worker/commit/77668e0e946003dc4f0d3b9e9c086c85b70f8fab))

#### 🔨 Chore

- changing the environment variable name for accessing GitHub MCP ([553c63d](https://github.com/specvital/worker/commit/553c63d358a5b1fd1c607843d41b90544d86330e))
- dump schema ([ba3fc16](https://github.com/specvital/worker/commit/ba3fc165a074f0827417ee6212002e79c9d5340e))
- dump schema ([425b609](https://github.com/specvital/worker/commit/425b6098dc1ee104189a4a33dc635f5e0b9f0352))
- dump schema ([52575e5](https://github.com/specvital/worker/commit/52575e5701088de44401abb227080800250094d8))
- dump schema ([abdaa2e](https://github.com/specvital/worker/commit/abdaa2eda93763d793b2a8a67f6fe2f3b4e14166))
- fix vscode import area not automatically collapsing ([ac92e87](https://github.com/specvital/worker/commit/ac92e87ee1be68a886e4df8b5ed006d0fc8ba0dd))
- improved the claude code status line to display the correct context window size. ([e1fa775](https://github.com/specvital/worker/commit/e1fa775b9dfd49ed57ec5d66aaf0eab4ec0e34b8))
- modified container structure to support codespaces ([0d1fec6](https://github.com/specvital/worker/commit/0d1fec6ec9af2bd3fb1df5a292242e240e13a36e))
- modify local db migration to always initialize the database ([7709a5b](https://github.com/specvital/worker/commit/7709a5b8af0a8fd7bee795ebd533dd5d3944d243))
- sync ai-config-toolkit ([0d00d4a](https://github.com/specvital/worker/commit/0d00d4a615fa3b1c162e8976b0f86b87948f0eaf))
- sync docs ([86772da](https://github.com/specvital/worker/commit/86772da7cb514400b7f7c89ea0defde95241195e))
- update core ([9092761](https://github.com/specvital/worker/commit/9092761f54e28b114b70a7dfbab14e8b82e27bdc))
- update core ([e6613c3](https://github.com/specvital/worker/commit/e6613c3a8e85189621056981ae0e3d91ff266e41))
- update core ([c163ae9](https://github.com/specvital/worker/commit/c163ae92f08d30046712de8c4b86b3162eaae758))

## [core/v1.5.0](https://github.com/specvital/core/compare/v1.4.0...v1.5.0) (2026-01-04)

### 🎯 Highlights

#### ✨ Features

- **cargotest:** detect test macros by analyzing same-file macro_rules! definitions ([4f3d697](https://github.com/specvital/core/commit/4f3d6975418aa9b381d8ebbdb59d429c119cbbcb))
- **junit4:** add JUnit 4 framework support ([7b96c63](https://github.com/specvital/core/commit/7b96c631594e7170273ff83358c4e32c16426854))
- **junit5:** add Java 21+ implicit class test detection ([d7c1218](https://github.com/specvital/core/commit/d7c1218c4fa8d8e6828b7940f4d6ed2507483d3a))
- **kotest:** detect tests defined in init blocks ([4a5a2d8](https://github.com/specvital/core/commit/4a5a2d81f294332cbf6b2dc30566923d8c338a11))
- **source:** add commit timestamp retrieval to GitSource ([afbe437](https://github.com/specvital/core/commit/afbe4372532ffeaff10c9a6f0622c2df70a881bc))
- **swift-testing:** add Apple Swift Testing framework support ([161b650](https://github.com/specvital/core/commit/161b650a186bdf2fc6a7e85d8b4303d7c3f84fc4))
- **vitest:** add test.for/it.for API support (Vitest 4.0+) ([5c7c8fa](https://github.com/specvital/core/commit/5c7c8fa6df22d18dda3257171617f68d1792a17d))

#### 🐛 Bug Fixes

- **cargo-test:** add macro-based test detection for Rust ([caa4d1b](https://github.com/specvital/core/commit/caa4d1bd76f2d42bfee1a48a3b01bbc21155ee83))
- **detection:** detect .cy.ts files as Cypress even within Playwright scope ([8ee2526](https://github.com/specvital/core/commit/8ee2526c4b1af920ff61410dbb5d9e31dbc0f96f))
- **dotnet:** detect tests inside C# preprocessor directives ([295f836](https://github.com/specvital/core/commit/295f836a750e4b1c8fb7813da03de58d9a1dc0e7))
- **dotnet:** support individual counting for parameterized test attributes ([e1d0d5f](https://github.com/specvital/core/commit/e1d0d5fff3b46490d4feec3b76498fdffeaa03b9))
- **gotesting:** support Test_xxx pattern test function detection ([fb7aeaf](https://github.com/specvital/core/commit/fb7aeafcbb5458cdef8ded2d2dba92a80d985c8a))
- **gtest:** add missing TYPED_TEST and TYPED_TEST_P macro detection ([cbb3914](https://github.com/specvital/core/commit/cbb391430c34fcd1117b7e263088eb553c068066))
- **gtest:** detect nested tests within tree-sitter ERROR nodes ([0ade3c7](https://github.com/specvital/core/commit/0ade3c7d443947918fa73861cc2910d0f998a5ea))
- **integration:** fix missing validation for multi-framework repositories ([5abc0a4](https://github.com/specvital/core/commit/5abc0a4c18bbc633d95a9c968c62a97e5eee7e3e))
- **jest:** support multiple root directories via Jest config roots array ([7e5bfea](https://github.com/specvital/core/commit/7e5bfeaee7940477d41f06265ee27a6400cd9347))
- **jstest:** add missing test detection in variable declarations ([d17f77d](https://github.com/specvital/core/commit/d17f77df86463b35c20ee13a5fbe4c6cbe22d5f8))
- **jstest:** count it.each/describe.each as single test ([fdbd484](https://github.com/specvital/core/commit/fdbd484070a159f24e988185cb3265358f75d5bb))
- **jstest:** detect jscodeshift defineTest calls as dynamic tests ([92bfb56](https://github.com/specvital/core/commit/92bfb5674afff7e89c3161fe1dc09f815b943c6e))
- **jstest:** detect tests inside IIFE conditional describe/it patterns ([1635945](https://github.com/specvital/core/commit/1635945c34026b2fadf4687f031cdb1f46790e6b))
- **jstest:** detect tests inside loop statements ([6c5b066](https://github.com/specvital/core/commit/6c5b066793971629fb7ed35f999b71e54c4833ec))
- **jstest:** detect tests using member expression as test name ([1570397](https://github.com/specvital/core/commit/157039798be41f72f258386c0026e642f4b53747))
- **jstest:** detect tests with variable names inside forEach callbacks ([efe2ec5](https://github.com/specvital/core/commit/efe2ec5b7dd9d2c3366519a24a6e698498770353))
- **jstest:** filter out Vitest conditional skip API from test detection ([3280998](https://github.com/specvital/core/commit/3280998a3fad6e900f2e2a3a0c85c95be07bb2ee))
- **jstest:** fix test detection failure in TSX files ([3d57940](https://github.com/specvital/core/commit/3d57940f70c76a3ed53904ae656088a2887950ff))
- **jstest:** fix test detection inside custom wrapper functions ([9c91958](https://github.com/specvital/core/commit/9c919581e4b8a3894a886f43737b7c0c32c7a572))
- **jstest:** support dynamic test detection in forEach/map callbacks and object arrays ([bc51894](https://github.com/specvital/core/commit/bc51894086b4d2730ea7f1dd78203e9e64f83ef9))
- **jstest:** support ESLint RuleTester.run() pattern detection ([b5e18f9](https://github.com/specvital/core/commit/b5e18f9f5dc975b25b12f11f1970ee3eb0f40903))
- **jstest:** support include/exclude pattern parsing in Jest/Vitest configs ([801e455](https://github.com/specvital/core/commit/801e45595aa0ec8d1a920bc396efc36a2aa9b716))
- **junit4:** detect tests inside nested static classes ([5673d83](https://github.com/specvital/core/commit/5673d83cee6c18c758adc5149b49bf411ea2ed83))
- **junit5:** add @TestFactory and @TestTemplate annotation support ([e868ef5](https://github.com/specvital/core/commit/e868ef56ddabd0ac85af83046c6e6b6c969c6eed))
- **junit5:** add custom @TestTemplate-based annotation detection ([242e320](https://github.com/specvital/core/commit/242e320505a7dfa7db664187babcc2e22f6d2b6f))
- **junit5:** detect Kotlin test files ([9090b51](https://github.com/specvital/core/commit/9090b51147cdd6ee4947ae8c02ed46fb666293a6))
- **junit5:** exclude JUnit4 test files from JUnit5 detection ([02aaed1](https://github.com/specvital/core/commit/02aaed191af6b47097f501f35d0db9421d53d79d))
- **junit5:** exclude src/main path from test file detection ([7e5ce26](https://github.com/specvital/core/commit/7e5ce26df4fdd5976c3874dd37130e6f4363497e))
- **kotest:** add missing WordSpec, FreeSpec, ShouldSpec style parsing ([424ab3a](https://github.com/specvital/core/commit/424ab3aa159fe13b9c132fa9117324cd0e5050d1))
- **kotest:** detect tests inside forEach/map chained calls ([cbd1fb1](https://github.com/specvital/core/commit/cbd1fb1bd85d39a1c0447a59dea58c05fb2624f5))
- **minitest:** resolve Minitest files being misdetected as RSpec ([93305d9](https://github.com/specvital/core/commit/93305d997c3cfdbd9dc945b59e022f9d611a682b))
- **mstest:** support C# file discovery under test/ directory ([95bcc31](https://github.com/specvital/core/commit/95bcc31e4f11aeee70b7bf00f52a91775f696618))
- **parser:** handle NULL bytes in source files that caused test detection failure ([d9f959c](https://github.com/specvital/core/commit/d9f959cf101cef4bdd11603650a5018a34894217))
- **phpunit:** add missing indirect inheritance detection for \*Test suffix ([106b73d](https://github.com/specvital/core/commit/106b73d5aa1c8cc4354d915d6bbe481044f655c8))
- **playwright:** detect conditional skip API calls as non-test ([129f0c0](https://github.com/specvital/core/commit/129f0c09ad8bc63c6b942c67ba8154874884aa3f))
- **playwright:** detect indirect import tests even with import type present ([e353cac](https://github.com/specvital/core/commit/e353cac3091a344f452e797c507a9cfd7483adfc))
- **playwright:** detect indirect imports and test.extend() patterns ([aa22e18](https://github.com/specvital/core/commit/aa22e18b15173d801805856dd4364bbdca55bbcb))
- **playwright:** fix config scope-based framework detection bug ([4983492](https://github.com/specvital/core/commit/4983492e1e4ac0908892f9764c8efde5f49e8d61))
- **playwright:** support test function detection with import aliases ([4ba46b7](https://github.com/specvital/core/commit/4ba46b752406b8439e15ea2aba2f04d07841c60f))
- **pytest:** fix unittest.mock imports being misclassified as unittest ([4ad41ed](https://github.com/specvital/core/commit/4ad41ed450747445ed3928619d6a2a9bf90e2352))
- **rspec:** detect tests inside loop blocks (times, each, etc.) ([a68b270](https://github.com/specvital/core/commit/a68b2705d793f145bfd4fff88419f464e5ad615a))
- **rspec:** resolve RSpec files being misdetected as Minitest ([b21dda9](https://github.com/specvital/core/commit/b21dda9a43b7d3bf5d70c221543c832bc3373aab))
- **scanner:** exclude **fixtures**, **mocks** directories from scan ([881f360](https://github.com/specvital/core/commit/881f36037d25cb6583b2e1fa30d3e715435801a6))
- **scanner:** fix symlink duplicate counting and coverage pattern bug ([aba78d1](https://github.com/specvital/core/commit/aba78d169a1f7909307687c67c45d67c19a84b99))
- **scanner:** use relative path instead of absolute path for test file detection ([e1937d2](https://github.com/specvital/core/commit/e1937d281ab27136b0df8847b8fec84a8de4baa3))
- **testng:** add missing class-level @Test annotation detection ([7912275](https://github.com/specvital/core/commit/79122754e9a3d413ebb69836434f9bad7bc5ab78))
- **testng:** detect @Test inside nested classes ([aaad38e](https://github.com/specvital/core/commit/aaad38ed4ff83d1f48e130935acea6f9bc934282))
- **xunit:** support custom xUnit test attributes (*Fact, *Theory) ([4845628](https://github.com/specvital/core/commit/48456284a5d5df720a05a1b7debfaaa22e1d977c))

### 🔧 Maintenance

#### 📚 Documentation

- **dotnetast:** document tree-sitter-c-sharp preprocessor limitation ([104ec57](https://github.com/specvital/core/commit/104ec5728101e4fed8fcb3593b61dbfc245bcbae))
- **validate-parser:** add ADR policy review step ([06e46ec](https://github.com/specvital/core/commit/06e46ec548e4c320f303cc004cade92e14261472))
- **validate-parser:** allow repos.yaml repos on explicit request and enforce Korean report ([a9c1852](https://github.com/specvital/core/commit/a9c18529a9e1f5a396cef46a0bae3622357219b1))

#### ✅ Tests

- add integration test case - kubrickcode/baedal ([b81e4ac](https://github.com/specvital/core/commit/b81e4ac432735abc4afeb1dcba1e7ee84b77b038))
- add integration test case - specvital/collector ([ba5e703](https://github.com/specvital/core/commit/ba5e703d7087fe3502e5368317ff054658087176))
- add integration test case - specvital/core ([8523567](https://github.com/specvital/core/commit/85235670ae9407086433b056c83226e23a05c7bb))
- add integration test case - specvital/web ([17b455f](https://github.com/specvital/core/commit/17b455fc05076c1e1e18ce080ca57156680c5a90))
- add test case - kubrickcode/quick-command-buttons ([14c93c6](https://github.com/specvital/core/commit/14c93c680b85b581e9d987d43313da2a4b908c01))
- **junit5:** update integration snapshots for Kotlin support ([3e9aa54](https://github.com/specvital/core/commit/3e9aa548ae730819731dbecebba1bf2b982fbfe9))

#### 🔨 Chore

- add custom commands for parser validation ([c8bd024](https://github.com/specvital/core/commit/c8bd0246c1b466000b1049805254e9954212b6c4))
- add flask as integration test ([cf63e7c](https://github.com/specvital/core/commit/cf63e7c879c01440f458fcf0876da7d7b62d7690))
- fix vscode import area not automatically collapsing ([218fb9e](https://github.com/specvital/core/commit/218fb9e5f9e93ceef583237b10854ac3fb6d546e))
- **integration:** update cypress test repository to v15.8.1 ([36d6040](https://github.com/specvital/core/commit/36d60408103d0310c74a62d85dc7c473e2275a18))
- setting up devcontainers for testing various languages ​​and frameworks ([0f3b08e](https://github.com/specvital/core/commit/0f3b08ec5ddb18655b70ccdf56f081cdddc71a5f))
- snapshot update ([601530e](https://github.com/specvital/core/commit/601530e22796d92101eeb7deada55c351108dc2b))
- sync docs ([d8ec48c](https://github.com/specvital/core/commit/d8ec48c4b2757bda51c8637d1afb420390541577))
- sync docs ([4716fb2](https://github.com/specvital/core/commit/4716fb2f67a5ce77d87caa71734054dabee57070))
- sync docs ([bd09a40](https://github.com/specvital/core/commit/bd09a40a6b72689e7b3579f10467c85f35d163b1))
- sync docs ([ba18e47](https://github.com/specvital/core/commit/ba18e4780f7ca7c1da8370d6f43f98b6e4c8bd97))
- sync docs ([ae6a331](https://github.com/specvital/core/commit/ae6a331bbf8cf291acab6c6b6ba44960766c2367))
- sync docs ([4f00b6c](https://github.com/specvital/core/commit/4f00b6c4004d7fb23cbbd5e75b0f507a1294c4a1))

## [infra/v1.2.0](https://github.com/specvital/infra/compare/v1.1.0...v1.2.0) (2026-01-04)

### 🎯 Highlights

#### ✨ Features

- **db:** add committed_at column to analyses table ([66a993d](https://github.com/specvital/infra/commit/66a993dcd00dc5ef891806c105ef6880cc106d2d))
- **db:** add external_repo_id column and integrity indexes ([848036b](https://github.com/specvital/infra/commit/848036b7a074c6e1f5549d436ae0db0ea9f502cb))
- **db:** add GitHub App Installation table ([cd33ecb](https://github.com/specvital/infra/commit/cd33ecb5c2f20d76355c91a826f4da6f7a0c5278))
- **db:** add GitHub cache tables for repository and organization data ([1605686](https://github.com/specvital/infra/commit/16056864c865991a87858815592b10db94b202f4))
- **db:** add is_private column to codebases table ([b688ba8](https://github.com/specvital/infra/commit/b688ba89c88eebeb2599a83a64a8324a9304bb04))
- **db:** add refresh token table for hybrid authentication ([0db7539](https://github.com/specvital/infra/commit/0db75399ddf1a326ba59c14e77a91fca05a32efa))
- **db:** add user_analysis_history table for dashboard personalization ([1044f38](https://github.com/specvital/infra/commit/1044f38993ce2629630fd9321de60ab64fd93a15))
- **db:** add user_bookmarks table for dashboard favorites ([7866748](https://github.com/specvital/infra/commit/78667485c8d51845dbb3c484adc0f40e57af78f6))

#### ⚡ Performance

- **db:** optimize index for cursor pagination ([d358516](https://github.com/specvital/infra/commit/d358516dd0eb603fcef8a59a998aa62578d4d484))

### 🔧 Maintenance

#### 📚 Documentation

- add CLAUDE.md ([5ef6ab0](https://github.com/specvital/infra/commit/5ef6ab0a933e3b0995acb08537f36f830dbf6589))
- add missing version headers and improve CHANGELOG hierarchy ([34c3614](https://github.com/specvital/infra/commit/34c3614a190afb5d31ab26bc27b70cfc6fe763fb))
- update README.md ([82b6396](https://github.com/specvital/infra/commit/82b6396cf7d276f81f15893e4883e226f58eb4ea))

#### ♻️ Refactoring

- **db:** change composite PK to surrogate PK for consistency ([dad65f8](https://github.com/specvital/infra/commit/dad65f846501a04ff648fe76c0b24a84efd041f8))

#### 🔨 Chore

- add sync-docs action command ([a8b519f](https://github.com/specvital/infra/commit/a8b519f03c8e1b46dcd73a31402cbfe387a754e6))
- auto-remove River DROP statements from makemigration ([53eb9ec](https://github.com/specvital/infra/commit/53eb9ece7b0359bbc7aa633e8a217620e6259c07))
- changing the environment variable name for accessing GitHub MCP ([3b74e68](https://github.com/specvital/infra/commit/3b74e68e41d19a0c44fc9b779e9f75c085eb2ef5))
- delete unused claude skills ([5c01ef8](https://github.com/specvital/infra/commit/5c01ef828ada131952325868c0ea5287eeb273ee))
- **deps-dev:** bump @semantic-release/release-notes-generator ([5197985](https://github.com/specvital/infra/commit/51979859d9a9b5796899874d81f476c29ab9315b))
- **deps:** bump actions/checkout from 4 to 6 ([8d1f8a4](https://github.com/specvital/infra/commit/8d1f8a4c99f42b378d889c452a24d250ee35b040))
- **deps:** bump actions/setup-node from 4 to 6 ([45ca48d](https://github.com/specvital/infra/commit/45ca48de2d3a9266eb23498d57bb82d6f320abb8))
- improved the claude code status line to display the correct context window size. ([928558e](https://github.com/specvital/infra/commit/928558e4d0f2070989d1cf475b2f855e9e9620a5))
- modified container structure to support codespaces ([558ee28](https://github.com/specvital/infra/commit/558ee28996f145f9f0b3a6d87f6892c91c0b081f))
- sync ai-config-toolkit ([bb51262](https://github.com/specvital/infra/commit/bb512622768223293c922300b3eb00d24423f2bd))
- sync docs ([34ab8a2](https://github.com/specvital/infra/commit/34ab8a24eed1824c3b3e3d9c5c1dfda948d9b254))
- sync docs ([9d595ac](https://github.com/specvital/infra/commit/9d595ac7477d30d956c18cd8d4cc689a6f6a02f6))

## [web/v1.1.2](https://github.com/specvital/web/compare/v1.1.1...v1.1.2) (2025-12-20)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix semantic-release CHANGELOG version header not rendering ([372a8ef](https://github.com/specvital/web/commit/372a8efc305de5267f727dfda02a3b285cbab22b))

## [web/v1.1.1](https://github.com/specvital/web/compare/v1.1.0...v1.1.1) (2025-12-20)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix semantic-release CHANGELOG version header not rendering ([5bff792](https://github.com/specvital/web/commit/5bff792d9dcca96eb0d08c7d2347a7940a784506))

## [web/v1.1.0](https://github.com/specvital/web/compare/v1.0.4...v1.1.0) (2025-12-20)

### 🎯 Highlights

#### ✨ Features

- add favicon and migrate to Next.js 16 proxy ([63ef413](https://github.com/specvital/web/commit/63ef413dc6d3a577e6a9a55447374ec7c86b181c))
- **ui:** apply Cloud Dancer theme color palette ([ac9039f](https://github.com/specvital/web/commit/ac9039fe8e5e3b221e918dffe9b28001daa67714))
- **ui:** enhance Stats Card visual hierarchy and unify CSS variables ([8690aa0](https://github.com/specvital/web/commit/8690aa0faffcc71ad394f226b17b0c88b4b7cf13))
- **ui:** improve accordion expand/collapse visual feedback ([beb2a1e](https://github.com/specvital/web/commit/beb2a1e977ff200a7a881aac6769fd474341ef76))
- **ui:** improve analysis page loading UX with skeleton and status banner ([ff30530](https://github.com/specvital/web/commit/ff305309f78632ff3775c809932a1d3ba856c11a))

#### 🐛 Bug Fixes

- fixed an error that occurred when a user was deleted and the user was not found. ([b73c723](https://github.com/specvital/web/commit/b73c7235326bed4d2d902d140c68c70af68df51b))
- **queue:** prevent duplicate analysis requests for same repository ([6acf9c3](https://github.com/specvital/web/commit/6acf9c350e7816efa739b7032c1f04ff9ee7c408))
- **ui:** input field indistinguishable from background color ([54cf29b](https://github.com/specvital/web/commit/54cf29b2e7d10cc88de6614ba24240d15538af61))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- claude settings workspace name fix ([21a8af8](https://github.com/specvital/web/commit/21a8af8b5ec13c4d0a57db10136ec1db6794cc26))

#### 📚 Documentation

- add missing version headers and improve CHANGELOG hierarchy ([4bb4427](https://github.com/specvital/web/commit/4bb44278b69753e72b01871e3be93fe7b2d546c3))

#### 🔨 Chore

- changing the environment variable name for accessing GitHub MCP ([e224631](https://github.com/specvital/web/commit/e224631e290786c81c7b9559d6d29a4f796acff0))
- collector -> worker structure and command changes ([2149fd8](https://github.com/specvital/web/commit/2149fd86bf7403e477e39f0ccc0f60ecc5bfd4ea))
- delete unused mcp ([c6b6124](https://github.com/specvital/web/commit/c6b6124899ae13ebeece23b89389426232ae2941))
- modified container structure to support codespaces ([ddca957](https://github.com/specvital/web/commit/ddca957df7a2619403fdde48699a08c0ec95c655))
- modify local db migration to always initialize the database ([e0de29d](https://github.com/specvital/web/commit/e0de29d329e177d85a090f99711f4f0d130b329c))
- sync ai-config-toolkit ([012bf02](https://github.com/specvital/web/commit/012bf02dc67f2fc37a3c4c168d8030ea04dabe94))

## [core/v1.4.0](https://github.com/specvital/core/compare/v1.3.0...v1.4.0) (2025-12-20)

### 🎯 Highlights

#### ✨ Features

- **crypto:** add NaCl SecretBox encryption package ([2bab1b3](https://github.com/specvital/core/commit/2bab1b313d720e7dcea1a148db6516202b25c035))
- **gotesting:** add Benchmark/Example/Fuzz function support ([76296d5](https://github.com/specvital/core/commit/76296d5d019cd91ee0cdd23f838d5a6a20c86494))
- **jstest:** add Jest/Vitest concurrent modifier support ([704b25c](https://github.com/specvital/core/commit/704b25c5e11972d86a6502aab918d873b8b7ec03))
- **mocha:** add Mocha TDD interface support ([2348b66](https://github.com/specvital/core/commit/2348b66975275675f88be5e96938ee73007774fc))
- **vitest:** add bench() function support ([b1f8949](https://github.com/specvital/core/commit/b1f89495c596f5ea6ae176200622a4678506ac0a))

#### 🐛 Bug Fixes

- disable implicit credential helper in git operations ([08f42b2](https://github.com/specvital/core/commit/08f42b24ef3fa0a889025f9ee08364f1d1eaa380))

### 🔧 Maintenance

#### 📚 Documentation

- add missing version headers and improve CHANGELOG hierarchy ([f38e681](https://github.com/specvital/core/commit/f38e6815e707bc2e0f91813cc5e29496c1349a3d))
- edit CLAUDE.md ([9ab5494](https://github.com/specvital/core/commit/9ab54945ab2d368e1f3a28b178152757728200ed))
- update README.md ([e76202e](https://github.com/specvital/core/commit/e76202ea4ecfae71cc2d077a85b88841b58e5f34))

#### 💄 Styles

- sort justfile ([78bacfd](https://github.com/specvital/core/commit/78bacfdc08f5b94ec6d1e5ef2aa253b7e52f2d6c))

#### 🔨 Chore

- ai-config-toolkit sync ([7d219c3](https://github.com/specvital/core/commit/7d219c3f04439b04609740456ac3e329c023e41c))
- changing the environment variable name for accessing GitHub MCP ([2079973](https://github.com/specvital/core/commit/20799738fab4dff39e53cd8b333b41d576afc465))
- delete unused mcp ([c3d1551](https://github.com/specvital/core/commit/c3d15516e01d0ef63770a0ef76eb55d641ab02af))
- **deps-dev:** bump @semantic-release/commit-analyzer ([9595ec8](https://github.com/specvital/core/commit/9595ec87709f40f1773cd1f90615295bde5b6baf))
- **deps-dev:** bump @semantic-release/github from 11.0.1 to 12.0.2 ([e05ee45](https://github.com/specvital/core/commit/e05ee45879c7d73ab619fe91bb7f064629aacfca))
- **deps-dev:** bump conventional-changelog-conventionalcommits ([9a13ed8](https://github.com/specvital/core/commit/9a13ed8541d955d77774762daab4d9834c1e863c))
- **deps:** bump actions/cache from 4 to 5 ([51d8d2b](https://github.com/specvital/core/commit/51d8d2b6c2ab69154f3ed645851bb7bb758d7465))
- **deps:** bump actions/checkout from 4 to 6 ([9bd5b6c](https://github.com/specvital/core/commit/9bd5b6c1feeadd1f81031f90457b0e8395053fea))
- **deps:** bump actions/setup-go from 5 to 6 ([12a3121](https://github.com/specvital/core/commit/12a31213d48b800004d9ec8d461c552a912eb94b))
- **deps:** bump actions/setup-node from 4 to 6 ([361f566](https://github.com/specvital/core/commit/361f5662b355e37b7e1be264ccf1c3059a13e0a5))
- **deps:** bump extractions/setup-just from 2 to 3 ([0fedea4](https://github.com/specvital/core/commit/0fedea40fa85390bc1c5a38873084dfa431a4790))
- **deps:** bump github.com/bmatcuk/doublestar/v4 from 4.8.1 to 4.9.1 ([d9f058f](https://github.com/specvital/core/commit/d9f058f8fef22cdb71de3db0e7d8191ef631834a))
- **deps:** bump github.com/stretchr/testify from 1.9.0 to 1.11.1 ([ab0bc83](https://github.com/specvital/core/commit/ab0bc8337c549d1e6fb44cbb98339947b0d99511))
- **deps:** bump golang.org/x/sync from 0.18.0 to 0.19.0 ([2b0070b](https://github.com/specvital/core/commit/2b0070bad52dfdf97102ef125c68eb0cdbf43417))
- Global document synchronization ([06c079c](https://github.com/specvital/core/commit/06c079c1487222eff2d771fc6db4d259a5bf273d))
- improved the claude code status line to display the correct context window size. ([365a13e](https://github.com/specvital/core/commit/365a13e86cd78593f63c51caea6f9df97dbbb200))
- modified container structure to support codespaces ([9e02cd4](https://github.com/specvital/core/commit/9e02cd44033a17317ef40f7e438eac1f0f013dcd))
- snapshot update ([78579ac](https://github.com/specvital/core/commit/78579accc74ec0596c3e1fae971fe33cb1da3e1e))
- snapshot update ([053ce8a](https://github.com/specvital/core/commit/053ce8a74203585873f6b578706cc7593e16511f))
- snapshot-update ([7c2fb1c](https://github.com/specvital/core/commit/7c2fb1c052e0af70880bde453622f25af0fb2410))
- sync ai-config-toolkit ([b7b852a](https://github.com/specvital/core/commit/b7b852ae34a6c46f0eb24471cb47fad528f13c77))
- sync ai-config-toolkit ([0b95b2d](https://github.com/specvital/core/commit/0b95b2d2bbc600fb23db92e7ddaabf41fbcb2957))

## [web/v1.0.4](https://github.com/specvital/web/compare/v1.0.3...v1.0.4) (2025-12-19)

### 🔧 Maintenance

#### ♻️ Refactoring

- migrate job queue from asynq to river ([72fce89](https://github.com/specvital/web/commit/72fce895b4cff07bef68244d7be08be59348b660))

## [worker/v1.0.6](https://github.com/specvital/worker/compare/v1.0.5...v1.0.6) (2025-12-19)

### 🎯 Highlights

#### 🐛 Bug Fixes

- resolve 60-second timeout failure during large analysis jobs ([ed18bc3](https://github.com/specvital/worker/commit/ed18bc3f587c0a446ea41b09431126ab1d22bba5))

## [worker/v1.0.5](https://github.com/specvital/worker/compare/v1.0.4...v1.0.5) (2025-12-19)

### 🎯 Highlights

#### 🐛 Bug Fixes

- resolve 60-second timeout issue during bulk INSERT operations on NeonDB ([0b6bc9b](https://github.com/specvital/worker/commit/0b6bc9bbef0c14190ec953a14caa5b29da0422d5))

## [worker/v1.0.4](https://github.com/specvital/worker/compare/v1.0.3...v1.0.4) (2025-12-19)

### 🔧 Maintenance

#### ♻️ Refactoring

- migrate queue system from asynq(Redis) to river(PostgreSQL) ([9664002](https://github.com/specvital/worker/commit/9664002057ef1f801dd8313e9f081760c3e0af21))

#### 🔨 Chore

- missing changes ([de9c0ec](https://github.com/specvital/worker/commit/de9c0ecaa7e136fe71d70a4a231dbd194ed7a33d))

## [infra/v1.1.0](https://github.com/specvital/infra/compare/v1.0.0...v1.1.0) (2025-12-19)

### 🎯 Highlights

#### ✨ Features

- **db:** add River job queue migration ([86b6157](https://github.com/specvital/infra/commit/86b61576794e3df0a097f151e67afad9f38c2abc))

### 🔧 Maintenance

#### 🔨 Chore

- adding a go environment to a container for riverqueue use ([ee45552](https://github.com/specvital/infra/commit/ee45552c4d80fd457c61df5f31c110534d4a0f7f))
- remove Redis dependency ([916c622](https://github.com/specvital/infra/commit/916c6227d3646e6d8baad8efe8e663e3f50b525b))

## [web/v1.0.3](https://github.com/specvital/web/compare/v1.0.2...v1.0.3) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- cookie not being set after GitHub login ([f4fccee](https://github.com/specvital/web/commit/f4fccee642db5089681837f66af66ee3b92a8e68))

## [web/v1.0.2](https://github.com/specvital/web/compare/v1.0.1...v1.0.2) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- "failed to get latest commit" error during repository analysis ([0de5c39](https://github.com/specvital/web/commit/0de5c399abe3d02435c81640c50d43d1a5bfa37f))

## [web/v1.0.1](https://github.com/specvital/web/compare/v1.0.0...v1.0.1) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- page not working in production environment ([21a60f7](https://github.com/specvital/web/commit/21a60f7700180cbe01faef41458cc5b73be645d0))

## [web/v1.0.0](https://github.com/specvital/web/commits/v1.0.0) (2025-12-18)

### 🎯 Highlights

#### ✨ Features

- add asynq queue client and DB repository infrastructure ([9b3136f](https://github.com/specvital/web/commit/9b3136f51a682aecccb13886542f023574fe8e7e))
- add C# xUnit test framework analysis support ([09878f5](https://github.com/specvital/web/commit/09878f57e339a7d1096de93315bb8409571b607c))
- add commit-based cache validation for analysis results ([a2170e5](https://github.com/specvital/web/commit/a2170e5cf8c2af056e98e8e75ad6dbe2ef548a0c))
- add context-aware Logger wrapper ([85a6c66](https://github.com/specvital/web/commit/85a6c669ee5d8ce1c6c4480bd1e63de59e8fa44e))
- add JUnit5 (Java) test framework analysis support ([3acaf97](https://github.com/specvital/web/commit/3acaf972d329acac3feff1095002bffdcd67d6b7))
- add local Redis/PostgreSQL services to devcontainer ([6c6281b](https://github.com/specvital/web/commit/6c6281bd9ad669e4aefbb0fcbd7fdbcee40adb21))
- add pytest framework analysis support ([46c2853](https://github.com/specvital/web/commit/46c2853f73ace027a28881002c7e77dd3f82119c))
- add run-collector command ([a06a2e7](https://github.com/specvital/web/commit/a06a2e7184757e74cf6837160bb4008fc74b46ad))
- add Scalar-based API documentation page ([402c7dd](https://github.com/specvital/web/commit/402c7dddf4ea268b7b0c5fea1220118dd950d979))
- add share button to analysis page ([9b85f5d](https://github.com/specvital/web/commit/9b85f5d03e154fff4fc4cc4897148b958e24de2f))
- **analyzer:** add user_id field to job payload ([40be02a](https://github.com/specvital/web/commit/40be02a8938156a309e944abffd71c3e1342db53))
- **analyzer:** update last_viewed_at on repository view ([07808a9](https://github.com/specvital/web/commit/07808a9d44272c9394abe0f5648005b2b83f5767))
- **auth:** add GitHub OAuth client and token encryption module ([32b4a0f](https://github.com/specvital/web/commit/32b4a0f9d674568621f5c74ffab9847e23452941))
- **auth:** add HTTP handler and JWT middleware for OAuth authentication ([ecc65d0](https://github.com/specvital/web/commit/ecc65d0d006032c33325c4226bd7c9d7c35ab63a))
- **auth:** add OAuth authentication foundation ([5ba0e4b](https://github.com/specvital/web/commit/5ba0e4b5c16a32e42d20f8fb70c80fb9d9592cad))
- **auth:** add private repo support and security hardening ([8c17231](https://github.com/specvital/web/commit/8c17231cbda537b3c291b2b5a35c03178b1139b1))
- **auth:** add repository and service layer for OAuth authentication ([9b98347](https://github.com/specvital/web/commit/9b98347fa3b2a374ac6aea937d877fbe9297a374))
- **auth:** complete frontend OAuth authentication integration ([62195e0](https://github.com/specvital/web/commit/62195e03eff4540da31ffc2f44e21feba565dffe))
- **backend:** implement GitHub API client package ([e710b85](https://github.com/specvital/web/commit/e710b850671f0133ca083f0fe3d345de056267d3))
- **backend:** implement real test file analysis with GitHub API integration ([fe33145](https://github.com/specvital/web/commit/fe331458e9a4486e68376dd14566cf822cc51bf0))
- extend schema to support multiple test framework statuses ([0c29d77](https://github.com/specvital/web/commit/0c29d773a6c4676a43e44b21232ed775e8e732ac))
- **i18n:** add Korean/English internationalization support ([b7da3df](https://github.com/specvital/web/commit/b7da3dfe0d29c347ddd3cb319df6bd4893800a5d))
- implement loading skeleton and error state handling ([55fb993](https://github.com/specvital/web/commit/55fb9938e61527cb86a8255a5d5e74b3ebfb73a5))
- implement mock analyzer API endpoint ([a60f0de](https://github.com/specvital/web/commit/a60f0de286833a1d7ac8c5817471a8140c39575f))
- implement Next.js 16 frontend basic structure ([0f5e6b1](https://github.com/specvital/web/commit/0f5e6b1cfef67cbd2f8f710cf2edc0d111d4269c))
- implement test dashboard UI with mock API integration ([a7c1733](https://github.com/specvital/web/commit/a7c17339c49d06350803ec4b5a98652cdd79ccb8))
- introduce TanStack Query for polling-based analysis status management ([86db089](https://github.com/specvital/web/commit/86db0890cfb2259d98c3cf2a646b411b846b5796))
- migrate analyzer module to queue-based architecture ([fddda57](https://github.com/specvital/web/commit/fddda579d952b4094c12cdecb8ba1ce97de83f46))
- set local mode as default execution environment ([85eef82](https://github.com/specvital/web/commit/85eef82a4254209bfb0710280b084a25b6cb3eb7))
- set up Go backend skeleton ([d32ebb5](https://github.com/specvital/web/commit/d32ebb5cee080711d8320b682f720e48738ee1b4))
- setup OpenAPI-based type generation pipeline ([cd60bb0](https://github.com/specvital/web/commit/cd60bb0ba92ef198962fc7a1cec4e9668cff5927))
- support for python unittest framework ([6ce91d6](https://github.com/specvital/web/commit/6ce91d6a7910686dcdcd7de15cae6c117c3a86fd))
- **ui:** add dark mode support ([28bd403](https://github.com/specvital/web/commit/28bd403277f4684d90f1f90ac7510015d2a1b279))
- **ui:** add empty state component for repositories without tests ([a38187f](https://github.com/specvital/web/commit/a38187fb93d964fa0a346fac9ec7d375aa46657b))
- **ui:** add framework breakdown to test statistics card ([eb2bf5d](https://github.com/specvital/web/commit/eb2bf5d8060b21cbd43701b78eb68b1c6eb948cf))
- **ui:** add global header with navigation ([d080965](https://github.com/specvital/web/commit/d0809655cd4da5fdcabb279fe3e783511d393105))

#### 🐛 Bug Fixes

- **analyzer:** allow retry for failed analysis requests ([4e101fc](https://github.com/specvital/web/commit/4e101fc9194478d2eb74edb7ed4beddf92d2c158))
- **analyzer:** Jest projects incorrectly detected as Vitest ([04cc006](https://github.com/specvital/web/commit/04cc0066c491e343a3ba85d997c424ccbfcb9b59))
- **analyzer:** Playwright tests incorrectly detected as Jest ([416c944](https://github.com/specvital/web/commit/416c9449e931192ea05a9017dc55454e8b84f8f7))
- **analyzer:** return empty slice instead of nil for suites ([4cf6aa3](https://github.com/specvital/web/commit/4cf6aa342b428235af19987f0ffeb8fc90634215))
- **analyzer:** vitest globals mode incorrectly detected as Jest ([2fa095d](https://github.com/specvital/web/commit/2fa095defe0421076e5431c380024e565d6ee276))
- **auth:** allow authenticated users to access private repositories ([f8babcb](https://github.com/specvital/web/commit/f8babcb74d6f66ad31346ec33da05f8b618a6ccd))
- **client:** block unauthenticated access to private repositories ([b89bacf](https://github.com/specvital/web/commit/b89bacfe47c1a5ecf59141142b3b4517ad385da4))
- **github:** fix double encoding of slashes in file paths ([8aa8d44](https://github.com/specvital/web/commit/8aa8d4452a60c031358b8c2a74b98537f0b7e80e))
- resolve hydration mismatch in LanguageSelector component ([dc6994a](https://github.com/specvital/web/commit/dc6994a00ede2eeedaf13e72ec1d75a24a21ffc3))
- test suites not displaying after analysis completion ([f9aa9b8](https://github.com/specvital/web/commit/f9aa9b848073ce37149c70aafe35cb135edab5c0))
- **test-list:** migrate from container scroll to page scroll virtualization ([2c45796](https://github.com/specvital/web/commit/2c45796059c80d1bf9a9880e17c743a1af7794b2))

#### ⚡ Performance

- **web:** enhance error handling and optimize large test list performance ([7f115c3](https://github.com/specvital/web/commit/7f115c3d58bea9d1724010147e19bd038b14f9b5))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **cors:** fix default CORS origin using backend port instead of frontend ([3729fb9](https://github.com/specvital/web/commit/3729fb9ccc7a1d405c03a5bc00248fe4deedb8bc))
- fix pnpm install failure in devcontainer ([7258e39](https://github.com/specvital/web/commit/7258e39e324900b3daadef56b1af5e24d3372c53))

#### 📚 Documentation

- add project documentation (README, API docs, CLAUDE.md) ([c08c730](https://github.com/specvital/web/commit/c08c730659987230e26eefda89083a85bc79a248))
- update CLAUDE.md ([01c22ea](https://github.com/specvital/web/commit/01c22eacff1da7dd748e6ee6a54a3cbf5dbbc300))

#### ♻️ Refactoring

- **analyzer:** abstract Service dependency with GitHost interface ([8baa79e](https://github.com/specvital/web/commit/8baa79ec9317fa38f4b3aacb99ae5c5b7baf4370))
- **analyzer:** delegate analysis record creation to collector-centric architecture ([06a3bf6](https://github.com/specvital/web/commit/06a3bf66a29d96d288ac7f3ad4baf873bb806d2a))
- **auth:** clarify GitHub OAuth environment variable naming ([be907c6](https://github.com/specvital/web/commit/be907c6c74b63976d4455b0c2e1b6402e4c983ee))
- **auth:** migrate to crypto package from specvital/core ([e3d974a](https://github.com/specvital/web/commit/e3d974aed9a0ffdf3c240efb492845efad00ec02))
- **backend:** consolidate duplicate port config into common package ([712508d](https://github.com/specvital/web/commit/712508da81c401fe06201d6baaf1d415f256a66c))
- **backend:** migrate from singleton to dependency injection pattern ([2848504](https://github.com/specvital/web/commit/284850426064d99507877b187f1bb0b3f42a8f93))
- **backend:** remove unused mock data code ([45bc693](https://github.com/specvital/web/commit/45bc693da125f45bf3e326738fdeb64ddecfa7a9))
- **backend:** reorganize package structure for clarity ([eb24d8b](https://github.com/specvital/web/commit/eb24d8b45af0a3e6652489e5eada6103482908b6))
- decouple HTTP status codes from service layer ([8ec957f](https://github.com/specvital/web/commit/8ec957fc2b487e259fc94015839ebc83a36a3b80))
- **frontend:** modernize data fetching with React 19 use() hook ([0c94f38](https://github.com/specvital/web/commit/0c94f38c7f65f086fd0f9284d9e837885b047a0d))
- **frontend:** reorganize to feature-based folder structure ([d0b471e](https://github.com/specvital/web/commit/d0b471e1a49986487477a0d07431ab91b05347bb))
- introduce APIHandlers composition pattern ([0df28f4](https://github.com/specvital/web/commit/0df28f40d12cf04348e00f67ed0a0e95828da74f))
- introduce Backend Service layer and enable Strict Server Interface ([31179da](https://github.com/specvital/web/commit/31179da5bb8282a8070f9f510c9237319dc5a47b))
- introduce domain layer to analyzer module ([8184bf2](https://github.com/specvital/web/commit/8184bf249ab30188b54258f7cff96dc31a471c6e))
- make framework type dynamic for better extensibility ([9654a6e](https://github.com/specvital/web/commit/9654a6e297d40b89404de71e2e7be693d118e15d))
- migrate frontend types to OpenAPI generated types ([5ce7203](https://github.com/specvital/web/commit/5ce7203e1bb0bf40f1b1317149b33a65006b8063))
- remove ~800 lines of unused frontend code ([b3d588a](https://github.com/specvital/web/commit/b3d588a2437d9fad6751bbf8bb56bd3a7231a025))
- remove httplog dependency and migrate to slog-based logging ([a31aac7](https://github.com/specvital/web/commit/a31aac7181275daad8b00d4e606adb1f9b4691ef))
- simplify framework imports with unified package ([5b56788](https://github.com/specvital/web/commit/5b56788d71d421fb9bc39c2976589652d63b4a47))

#### 🔧 CI/CD

- add OpenAPI type sync verification CI and update documentation ([f5a03b2](https://github.com/specvital/web/commit/f5a03b247ec163947123f9b91560aee6b59efc64))
- add Railway deployment infrastructure and semantic-release setup ([acca511](https://github.com/specvital/web/commit/acca5115ab7cfd33c86936e4df55470e9c4b3c6c))

#### 🔨 Chore

- add a port shutdown command ([0206fab](https://github.com/specvital/web/commit/0206fabbfd6dc7380999d335fb17218d20e65ed4))
- Add an item to gitignore ([6673219](https://github.com/specvital/web/commit/6673219d6a69395709da325424b460d72d3912c0))
- add bootstrap command ([68a82ef](https://github.com/specvital/web/commit/68a82efddbdabb8eca503c4044ca087045d49c17))
- add install-oapi-codegen in justfile ([58b68ee](https://github.com/specvital/web/commit/58b68eeb876b07ac8bb44a4e0553edc568677fcf))
- add integration run buttons ([0bca0b8](https://github.com/specvital/web/commit/0bca0b8a22d2dcb163ec333aa2bc982e9e09f18c))
- add next-env.d.ts ignore ([102610b](https://github.com/specvital/web/commit/102610b67e8f6b0fa650fdb8150ba032aa01fa3d))
- add run collector command ([c2e1bc1](https://github.com/specvital/web/commit/c2e1bc17d3fcad6548c27e7d2f064724bd7a7f07))
- add specvital-network connection to devcontainer ([564eab4](https://github.com/specvital/web/commit/564eab47b8e732a27e1f4e5b5373c5fe75aeedae))
- add specvital/core package update command ([653ad1f](https://github.com/specvital/web/commit/653ad1fd9f79c3102e1f7f0b3c2a60392fb2245f))
- add sqlc and PostgreSQL infrastructure setup ([e44cc89](https://github.com/specvital/web/commit/e44cc891f987c064ad539939c8058c43aba10361))
- add useful action button ([e77796d](https://github.com/specvital/web/commit/e77796d98d67e5b9ea4a002dcc84d0eb70c9ce37))
- add useful action buttons ([908d637](https://github.com/specvital/web/commit/908d637f9e8d5731ba144a374e63a616600c6d7d))
- add useful action buttons ([6fac36a](https://github.com/specvital/web/commit/6fac36a85f7b18e99db2148bc6d5ee9f71387c7b))
- added a playwright infrastructure installation command ([81e9283](https://github.com/specvital/web/commit/81e92831c5e91a7a80221971c375cf91ba918497))
- adding recommended extensions ([9ad2922](https://github.com/specvital/web/commit/9ad292289124f05915d988582af3c869be758781))
- ai-config-toolkit sync ([a407dde](https://github.com/specvital/web/commit/a407dde3feefc3ba61fd9fb6d8279ccf7e9ed561))
- ai-config-toolkit sync ([c4b91e6](https://github.com/specvital/web/commit/c4b91e61a29a4eb36daea40a965592366b597e33))
- ai-config-toolkit sync ([272a8fb](https://github.com/specvital/web/commit/272a8fb66a0cf9cb98ad6091b1dc993c9fad75db))
- change backend port ([40020e6](https://github.com/specvital/web/commit/40020e671ef2729f6ed31b446e92a426e982c4c2))
- chore action buttons ([1e84ff4](https://github.com/specvital/web/commit/1e84ff4be4462e199e902f93cdb7968e6edd0ab2))
- **deps-dev:** Bump prettier from 3.6.2 to 3.7.4 ([45c5f72](https://github.com/specvital/web/commit/45c5f720a21c06a6eaac72d679e69b5bddf233b3))
- dump schema ([ebd30b8](https://github.com/specvital/web/commit/ebd30b8fe7007bbb9d15e0d32c00154026dc1be2))
- dump schema ([0fb59d4](https://github.com/specvital/web/commit/0fb59d4ea500e4789141cc5be2859f3accf68ae1))
- dump schema ([c60bf06](https://github.com/specvital/web/commit/c60bf0671467f74af2c530e5bb5c43996054e4fd))
- Global document synchronization ([1a6dfe8](https://github.com/specvital/web/commit/1a6dfe88c4245e6279cf5bea05891e0251daf2bb))
- sync ai-config-toolkit ([013160f](https://github.com/specvital/web/commit/013160f9cfb07ab7d9f2149e0900c1f9f3bc0667))
- sync core ([fe0effc](https://github.com/specvital/web/commit/fe0effcdb798e4ce63fb5641556c6bcea1daf0df))
- syncing documents from ai-config-toolkit ([3fb53a0](https://github.com/specvital/web/commit/3fb53a09fdc7409301c622278882ad64b4f1c9d4))
- update package.json ([28eab4e](https://github.com/specvital/web/commit/28eab4e8052f17aacc9ecd6e60bf5184cc319328))
- upgrade zod from 3.25.28 to 4.1.13 ([dc7e48b](https://github.com/specvital/web/commit/dc7e48b0c5cd205871ac81f81e2a4893f515c53c))

## [worker/v1.0.3](https://github.com/specvital/worker/compare/v1.0.2...v1.0.3) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- fix git clone failure in runtime container ([f11dfa3](https://github.com/specvital/worker/commit/f11dfa3e4090a6412af71b58a7eca6f081e49d4d))

### 🔧 Maintenance

#### ♻️ Refactoring

- remove unused dead code ([95cee17](https://github.com/specvital/worker/commit/95cee17e1307fe3e0cc23ba7d549292b05c19744))

#### 🔨 Chore

- sync docs ([9007e97](https://github.com/specvital/worker/commit/9007e97bbff365afcae69cfcba2f501732a20c8b))

## [worker/v1.0.2](https://github.com/specvital/worker/compare/v1.0.1...v1.0.2) (2025-12-17)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix asynq logs incorrectly classified as error in Railway ([d2180cc](https://github.com/specvital/worker/commit/d2180cc1182a0f1187f7dd63b982fc7816e3be47))

## [worker/v1.0.1](https://github.com/specvital/worker/compare/v1.0.0...v1.0.1) (2025-12-17)

### 🎯 Highlights

#### 🐛 Bug Fixes

- enable CGO for go-tree-sitter build ([50b1fea](https://github.com/specvital/worker/commit/50b1fea3c7834bd585c3a23615d8acb5cbae8a5f))

## [worker/v1.0.0](https://github.com/specvital/worker/releases/tag/v1.0.0) (2025-12-17)

### 🎯 Highlights

#### ✨ Features

- add adaptive decay logic for auto-refresh scheduling ([8a85854](https://github.com/specvital/worker/commit/8a858547c7b8a5190176253830b862588cda8042))
- add enqueue CLI tool ([5697cb9](https://github.com/specvital/worker/commit/5697cb9533b8dcfd8c7c90fa34d5419b181fa287))
- add focused/xfail to test_status and support modifier column ([cd60233](https://github.com/specvital/worker/commit/cd602333aa9eba34e8a3b21bbcd91040bfe59936))
- add job timeout to prevent long-running analysis jobs ([392b43e](https://github.com/specvital/worker/commit/392b43e7395aa923fd05a96872e5f0c8911a8845))
- add local development mode support to justfile ([2ca2d51](https://github.com/specvital/worker/commit/2ca2d51337cd8c425e54a5f3ffd257bde28e403c))
- add local development services to devcontainer ([a30ca9e](https://github.com/specvital/worker/commit/a30ca9eca3df2aa4a418d94b1525f8ac170ee6c1))
- add OAuth token parameter to VCS Clone interface ([de518d0](https://github.com/specvital/worker/commit/de518d05eb534c97cf5e647168ff9d8a686ae00e))
- add semaphore to limit concurrent git clones ([9ddbc06](https://github.com/specvital/worker/commit/9ddbc06291c8fb41e0e39def7f32edaa5830f1ba))
- add UserRepository for OAuth token lookup ([9a16ec1](https://github.com/specvital/worker/commit/9a16ec1f37031d15612198b27f5316d4ab066225))
- implement analysis pipeline (git clone → parse → DB save) ([66dd262](https://github.com/specvital/worker/commit/66dd2627b6e7d5b46ab7bc4358a1f6178d77cfee))
- implement asynq-based worker basic structure ([4dd16ad](https://github.com/specvital/worker/commit/4dd16ad22b355229f6c9db12e076427f7ff0c2ea))
- initialize collector service project ([1d3c8cf](https://github.com/specvital/worker/commit/1d3c8cf3a570a719301fa953c9cccb9c53a0358a))
- integrate OAuth token lookup logic into UseCase ([cb3f911](https://github.com/specvital/worker/commit/cb3f9114f89b9435f037379b6eac07c51dc53d96))
- integrate scheduler for automatic codebase refresh ([e0a1a15](https://github.com/specvital/worker/commit/e0a1a15dac6ef26b62debf1d3b52d172cf8d8ed6))
- record failure status in DB when analysis fails ([6485ac3](https://github.com/specvital/worker/commit/6485ac345c9562404e435ea32962a1b90a13fd5b))
- support external analysis ID for record creation ([5448202](https://github.com/specvital/worker/commit/54482021c84b6faa466f0ce006de06ddcd79d22d))
- support OAuth token decryption for private repo analysis ([8d0ad30](https://github.com/specvital/worker/commit/8d0ad307ce07aa7562daa6a38230e19ce2cc1644))

#### 🐛 Bug Fixes

- handle missing error logging and DB status update on analysis task failure ([64ae8d9](https://github.com/specvital/worker/commit/64ae8d9fcbc620e7470e2c7cc21ade39d7327f8d))
- parser scan failing due to unexported method type assertion ([6256673](https://github.com/specvital/worker/commit/6256673dc23d48998851f977fcd034498c591642))
- remove unnecessary wait and potential deadlock in graceful shutdown ([b78c981](https://github.com/specvital/worker/commit/b78c981c662d32fde66c0890da0e226e9b4a4d3e))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- go mod tidy ([c58f73b](https://github.com/specvital/worker/commit/c58f73b40f2de49c8c69b1d67efd45b1487c0359))

#### 💄 Styles

- format code ([5e994e2](https://github.com/specvital/worker/commit/5e994e2ab90f6ae6a8cd64d392b946c9bde0bd1d))

#### ♻️ Refactoring

- centralize dependency wiring with DI container ([c1b8215](https://github.com/specvital/worker/commit/c1b82151bdba8b62e194100e6f04271fd3f4e026))
- extract domain layer with zero infrastructure dependencies ([7ba9e51](https://github.com/specvital/worker/commit/7ba9e51a0ffa76327736e91c828fe5949cbfbcb6))
- extract repository layer from AnalyzeHandler ([464ecfa](https://github.com/specvital/worker/commit/464ecfa6087d91d4d399a7fee032ed2a9109a151))
- extract service layer from AnalyzeHandler ([d9faf20](https://github.com/specvital/worker/commit/d9faf200da77b3097417cd1671a3a1c5fbc5fe06))
- implement handler layer and clean up legacy packages ([23a093f](https://github.com/specvital/worker/commit/23a093f6d5558d10decb21272289c1d99e583101))
- implement repository adapter layer (Clean Architecture Commit 3) ([8b0e433](https://github.com/specvital/worker/commit/8b0e43372fec1b15fa7e7d76794be37d42b6988e))
- implement use case layer with dependency injection ([b2be6ff](https://github.com/specvital/worker/commit/b2be6ff3d3a1c3662f816958d614f4a89215aba6))
- implement VCS and parser adapter layer (Clean Architecture Commit 4) ([1b2e34f](https://github.com/specvital/worker/commit/1b2e34f61665d6513c2654712390b67524dfd731))
- move infrastructure packages to internal/infra ([6cc1d1c](https://github.com/specvital/worker/commit/6cc1d1caf8722317967eae5c2de6c40c71467ce2))
- separate Scheduler from Worker into independent service ([9481141](https://github.com/specvital/worker/commit/9481141e99f0adcc225e93d05c8104846f836c17))
- split entry points for Railway separate deployments ([d899192](https://github.com/specvital/worker/commit/d899192cb1772fe9ed16d426d460e016c1bbf2ee))

#### ✅ Tests

- add AnalyzeHandler unit tests ([0286e7c](https://github.com/specvital/worker/commit/0286e7cd687d65be853e503a883cd74010f8dede))
- remove unnecessary skipped tests ([f8c0eb4](https://github.com/specvital/worker/commit/f8c0eb40a4f99252650bbf9e0f9ca93a378223fb))

#### 🔧 CI/CD

- configure semantic-release automated deployment pipeline ([37f128f](https://github.com/specvital/worker/commit/37f128f2c9d113144d2af530e88f84d3209f235c))

#### 🔨 Chore

- add bootstrap command ([c8371f0](https://github.com/specvital/worker/commit/c8371f0d5f47a19353c260ecf83c5033b4e5ba53))
- add Dockerfile for collector service ([6e3b0e4](https://github.com/specvital/worker/commit/6e3b0e4225b4b0875e2ee0bae909a594b1b9f87c))
- add example env file ([64a24a4](https://github.com/specvital/worker/commit/64a24a4de88aa5e4954a59b049b915c2012da79e))
- add gitignore item ([8fc64a6](https://github.com/specvital/worker/commit/8fc64a6ab0a3abca1cf1d73f458adf17bb752ced))
- add migrate local command ([baabcfe](https://github.com/specvital/worker/commit/baabcfe97f3905122026a752ea2ba7f7ed07917b))
- add PostgreSQL connection and sqlc configuration ([eecc4a6](https://github.com/specvital/worker/commit/eecc4a69a8b8c6e5c67a8f652a97ae784ecca1c1))
- add useful action buttons ([02fa778](https://github.com/specvital/worker/commit/02fa7785ac4ba1505e06ac3add60621cf01d1be9))
- adding recommended extensions ([30d5d0b](https://github.com/specvital/worker/commit/30d5d0b0fccc3190313456433e24e1342c18d641))
- ai-config-toolkit sync ([3091cf4](https://github.com/specvital/worker/commit/3091cf46ca2e6a24f5c299fe4f8008659fe1b8c8))
- ai-config-toolkit sync ([decf96b](https://github.com/specvital/worker/commit/decf96b2c47b278ef56a3e76c6174c9688f883c3))
- delete file ([f48005c](https://github.com/specvital/worker/commit/f48005cad322fa2586e9bf315e2bce3c608dcd8b))
- dump schema ([b90bab0](https://github.com/specvital/worker/commit/b90bab0f0c33d52683b6ff6a1f132702eb54a077))
- dump schema ([370409c](https://github.com/specvital/worker/commit/370409cee67512bc3f21ac3f5835357303db9b57))
- dump schema ([d704305](https://github.com/specvital/worker/commit/d7043054a0a59ac755bc23a85a4fd39f5ce97a0a))
- Global document synchronization ([cead255](https://github.com/specvital/worker/commit/cead255f25f48397848d55cf1417f21466dae67c))
- sync ai-config-toolkit ([e559889](https://github.com/specvital/worker/commit/e55988903526ade4630d2d6516e67ad1354ff67e))
- update core ([d358131](https://github.com/specvital/worker/commit/d358131e3e6197ee8958655b3cc1cfa7d0ed9ca6))
- update core ([b47592e](https://github.com/specvital/worker/commit/b47592e2a6668c25585d0338099e83e7b72bf1d5))
- update core ([395930a](https://github.com/specvital/worker/commit/395930a21bb48b8283cac037cde1999e44ae69c6))
- update schema.sql path in justfile ([0bcbe79](https://github.com/specvital/worker/commit/0bcbe794cdbccff58e2babe75a6308aacc6ad5d0))
- update-core version ([cc65b03](https://github.com/specvital/worker/commit/cc65b0325a1e828e24270753d76fa91ff01eeb45))

## [infra/v1.0.0](https://github.com/specvital/infra/releases/tag/v1.0.0) (2025-12-17)

### 🎯 Highlights

#### ✨ Features

- add Atlas-based database schema management ([da9fb70](https://github.com/specvital/infra/commit/da9fb70f603c5cbc686b1f0412350f29d18969fa))
- add PostgreSQL/Redis infrastructure for local development ([a86dc00](https://github.com/specvital/infra/commit/a86dc0074e954c85b5cf94e0225eeec4fcaddf9f))
- **db:** add last_viewed_at column for auto-refresh service ([7f2b1cf](https://github.com/specvital/infra/commit/7f2b1cf1fa24462df960827620529c2c474d04bc))
- **db:** add users and oauth_accounts tables for GitHub OAuth ([3295843](https://github.com/specvital/infra/commit/3295843b40edafe4cffe2c37917f4a2c807aec4a))
- extend schema for multi-framework test status support ([cc2531e](https://github.com/specvital/infra/commit/cc2531e9e62b7aa567c0497023ece0e6e8d8e87a))

#### 🐛 Bug Fixes

- **ci:** add revisions_schema config and allow-dirty flag for atlas migration ([5a71d60](https://github.com/specvital/infra/commit/5a71d608ac406eeb344feb28c9404a50f484d0fd))
- **db:** test case save failure when name exceeds 500 characters ([9598962](https://github.com/specvital/infra/commit/9598962b24aeb60bf8ce579441e26bb4d722b5a8))
- **db:** unique constraint violation on analysis retry ([bb10f8a](https://github.com/specvital/infra/commit/bb10f8ae749d5ed64190e0ba0bd7f2ead1012a16))

### 🔧 Maintenance

#### 💄 Styles

- format code ([b8b1d36](https://github.com/specvital/infra/commit/b8b1d36e93a49886faccd52b282d7c6879d8f2b2))

#### 🔧 CI/CD

- add release workflow for semantic-release ([817f077](https://github.com/specvital/infra/commit/817f0776175cf311f9cbcd098603fb6a9a4145f3))
- setup NeonDB migration and release automation pipeline ([fd3a039](https://github.com/specvital/infra/commit/fd3a03936691fd5d80917d9d592914d3a97fffcb))

#### 🔨 Chore

- add "hashicorp.hcl" extension in recommended ([6b063b1](https://github.com/specvital/infra/commit/6b063b184832c2de0258178a34635a7f379a49d1))
- add claude session volume ([5d2f745](https://github.com/specvital/infra/commit/5d2f745177332acccd940dd4d65f3895a080560f))
- add neon db extension ([1324222](https://github.com/specvital/infra/commit/1324222f4c31b9edc7c95bc5462d69f69ed41cc1))
- add Redis reset capability to reset command ([4840861](https://github.com/specvital/infra/commit/48408613f85a53b658d9ece231ba7928459a2e08))
- add release command ([bb79d68](https://github.com/specvital/infra/commit/bb79d68f25b91a30fc593cd63d56497b93992299))
- add useful action buttons ([219fb7f](https://github.com/specvital/infra/commit/219fb7ff45e79f97f573a612cf512bfaf664f75d))
- adding recommended extensions ([0d4b17a](https://github.com/specvital/infra/commit/0d4b17a924b956b709cfcbaf715c9f3bb02427b2))
- ai-config-toolkit sync ([0a2fa86](https://github.com/specvital/infra/commit/0a2fa868a46e3c040ae8d730221ace3f6b032775))
- ai-config-toolkit sync ([c78e010](https://github.com/specvital/infra/commit/c78e010b6caaf97a4f5274db4482e19841399bf5))
- Global document synchronization ([15dc7da](https://github.com/specvital/infra/commit/15dc7dad10632e8c505efeb0459eea5feee2a0f7))
- sync ai-config-toolkit ([d4dc1d6](https://github.com/specvital/infra/commit/d4dc1d68dc85fab03d5467ac0f7d4359da52f162))

## [core/v1.3.0](https://github.com/specvital/core/compare/v1.2.2...v1.3.0) (2025-12-11)

### 🎯 Highlights

#### ✨ Features

- **cypress:** add Cypress E2E testing framework support ([b87f92b](https://github.com/specvital/core/commit/b87f92b68e0e0263f9571a238693e6a96390c232))
- **gtest:** add C++ Google Test framework support ([3821565](https://github.com/specvital/core/commit/382156588fe30534679fb24e7005cd45be1315e7))
- **kotest:** add Kotlin Kotest test framework support ([374b696](https://github.com/specvital/core/commit/374b6963d5d9a666297c450ec2d8aeecb8adfca9))
- **minitest:** add Ruby Minitest test framework support ([2cf2c22](https://github.com/specvital/core/commit/2cf2c22eb6fb3964b3611ee72f12d50a96cb4566))
- **mocha:** add Mocha JavaScript test framework support ([fdbb49e](https://github.com/specvital/core/commit/fdbb49e7d7d048238da0e1af53e1f4ed3d24b627))
- **mstest:** add C# MSTest test framework support ([9ab6565](https://github.com/specvital/core/commit/9ab6565c8141a0b8fe4a80eb197da76432306bb0))
- **phpunit:** add PHP PHPUnit framework support ([d395cd5](https://github.com/specvital/core/commit/d395cd5627a94e54eeaee989c0e48817810f79ca))
- **testng:** add Java TestNG test framework support ([3c50d31](https://github.com/specvital/core/commit/3c50d318b58031b0c833c6872de8c248aa97cacb))
- **xctest:** add Swift XCTest test framework support ([7c62c95](https://github.com/specvital/core/commit/7c62c95a2f4b7120b938daf7f6e09466479f3d23))

### 🔧 Maintenance

#### ♻️ Refactoring

- delete deprecated code ([c55a9ac](https://github.com/specvital/core/commit/c55a9ac53695ae0e4bfac1f5c70d9c20b9a20239))
- remove dead code from MVP development ([3a08c21](https://github.com/specvital/core/commit/3a08c21fe0e15ed11a793bb2a7943b1c53736191))

#### 🔨 Chore

- add missing framework constants ([324bad0](https://github.com/specvital/core/commit/324bad03d66c25fcaeb4ee6be0859e2dfa5e2607))

## [core/v1.2.2](https://github.com/specvital/core/compare/v1.2.1...v1.2.2) (2025-12-10)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **release:** fix 404 error on commit links in release notes ([3bcff5e](https://github.com/specvital/core/commit/3bcff5e9498bec9aa56edbb9797d51263888088b))

## [core/v1.2.1](https://github.com/specvital/core/compare/v1.2.0...v1.2.1) (2025-12-10)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **release:** fix broken commit links and long hash display in release notes ([fe38507](https://github.com/specvital/core/commit/fe3850790f60df701af655b4e7177899bfcb80ff))

#### 🔨 Chore

- adding recommended extensions ([328447f](https://github.com/specvital/core/commit/328447f811601b35b6ca2e71c3bf83183a77af35))

## [core/v1.2.0](https://github.com/specvital/core/compare/v1.1.2...v1.2.0) (2025-12-10)

### 🎯 Highlights

#### ✨ Features

- add all package for bulk parser strategy registration ([96ffbe6](https://github.com/specvital/core/commit/96ffbe688e750a18df1556b9e41157f4a0d4306e))
- add C# language and xUnit test framework support ([3b3c685](https://github.com/specvital/core/commit/3b3c685c6fb0f26cbf4b6865a0dd64f2231fac55))
- add GitSource implementation for remote repository access ([c8743a5](https://github.com/specvital/core/commit/c8743a5872ea641f2916c035c57f88960122da77))
- add Java language and JUnit 5 test framework support ([cc1a6ba](https://github.com/specvital/core/commit/cc1a6ba153bc5e9ea4c5af9e5c2672a2cd9020a7))
- add NUnit test framework support for C# ([b62c420](https://github.com/specvital/core/commit/b62c4208777cbc43f8e63af370ef0c9c01636f39))
- add Python pytest framework support ([b153129](https://github.com/specvital/core/commit/b153129a6dc81e3c56f90c33a03731d45cee5b1c))
- add Python unittest framework support ([bcac628](https://github.com/specvital/core/commit/bcac628e882152a610c2d898c93ac9e5824c642e))
- add Ruby language and RSpec test framework support ([3e28c47](https://github.com/specvital/core/commit/3e28c476c8f338b8ab25e08774feed4ca272fd5e))
- add Source interface and LocalSource implementation ([af0e2ed](https://github.com/specvital/core/commit/af0e2ed2e49c51bfc5d788f7e96397a862a8850a))
- **domain:** add Modifier field to Test/TestSuite ([a1b9363](https://github.com/specvital/core/commit/a1b93633275941ddd3b4fee7db46327a656eab21))
- **parser:** add Rust cargo test framework support ([30feca7](https://github.com/specvital/core/commit/30feca749a0d9c64b4eccb7ea6ed5a66c3ab4516))
- **source:** add Branch method to GitSource ([8d6f10d](https://github.com/specvital/core/commit/8d6f10d556506732d3b00750352889f992e9520e))
- **source:** add CommitSHA method to GitSource ([97256ec](https://github.com/specvital/core/commit/97256ec27766fdba1fb67942fc1d58fe4252f36f))
- **vitest:** add VitestContentMatcher for vi.\* pattern detection ([9d2c72e](https://github.com/specvital/core/commit/9d2c72e8fdf71bfca654e7835be902c05e698862))

#### 🐛 Bug Fixes

- **detection:** fix Go test files not being detected ([8487f71](https://github.com/specvital/core/commit/8487f71642be502c3a6ba66ba29398bae273d42b))
- **detection:** fix scope-based framework detection bugs ([3589928](https://github.com/specvital/core/commit/35899280fbf45a7ec8dae2987a51fb48143adef2))
- **parser:** prevent slice bounds panic in tree-sitter node text extraction ([465e9bc](https://github.com/specvital/core/commit/465e9bc0d0aeee688c29ab786ddeafbb88d76d87))
- **tspool:** fix flaky tests caused by tree-sitter parser reuse ([256c9aa](https://github.com/specvital/core/commit/256c9aa1780471334ee0d28ede877b050a5cc2d6))

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix nondeterministic integration test results ([41e3d38](https://github.com/specvital/core/commit/41e3d3831892ca52c59e621d75172651ca0ecbdc))

#### 💄 Styles

- format code ([71d8f66](https://github.com/specvital/core/commit/71d8f66631e6fb29e55e9d3ea934806e1a1b806f))

#### ♻️ Refactoring

- change Scanner to read files through Source interface ([11507ac](https://github.com/specvital/core/commit/11507accf9d0a9f34e18cb8bdaf80e62f6333c5e))
- **detection:** redesign with unified framework definition system ([9ba32af](https://github.com/specvital/core/commit/9ba32af300f73bf08746ac24e3fcb4ea48d5291b))
- **detection:** replace score accumulation with early-return approach ([ab30e72](https://github.com/specvital/core/commit/ab30e72e4d2a2bcb4d45baed9eac8cc422286ba5))
- **domain:** align TestStatus constants with DB schema ([babec36](https://github.com/specvital/core/commit/babec3602a02ece88b8a22b8729f335a96163555))

#### ✅ Tests

- add 8 complex case repositories for edge case coverage ([619f361](https://github.com/specvital/core/commit/619f361801a76059e7e4f7e8206a1486e67de420))
- add golden snapshot comparison to integration tests ([1cffd01](https://github.com/specvital/core/commit/1cffd019a34302a9f4d253cda11d6868d6fe61f9))
- add integration test infrastructure with real GitHub repos ([476b3eb](https://github.com/specvital/core/commit/476b3eb16953add6a64023f64bb68aa4de8e841f))
- add unittest integration test repositories ([7d31dcf](https://github.com/specvital/core/commit/7d31dcfa256d9106bf831e595831e35722b5e72e))

#### 🔧 CI/CD

- add integration test CI workflow and documentation ([d9368e1](https://github.com/specvital/core/commit/d9368e181da2745c02b007652c00694dc88b0d7d))

#### 🔨 Chore

- add snapshot-update command and refresh golden snapshots ([c3e47e8](https://github.com/specvital/core/commit/c3e47e8bf274eef00f0088a4a13d7a66a24c072b))
- add useful action buttons ([ef1a60c](https://github.com/specvital/core/commit/ef1a60cd9f88ca7457e366bc1978c03750019316))
- ai-config-toolkit sync ([e631a30](https://github.com/specvital/core/commit/e631a30fde776b9ba023ec00989cf2a8605e39d6))
- ai-config-toolkit sync ([42eeba3](https://github.com/specvital/core/commit/42eeba3426c41231ebefa9fc431fd3884f954b2d))
- snapshot update ([f4c171d](https://github.com/specvital/core/commit/f4c171dbf86a02cfa471806d1da98f014899c161))
- sync integration repos ([02c6a8d](https://github.com/specvital/core/commit/02c6a8d4311bcae40bca218e8a2081b8392a4755))
- sync snapshot ([6c086e9](https://github.com/specvital/core/commit/6c086e9c4297b074b7184678c3fde40a5bbdc00f))

## [core/v1.1.2](https://github.com/specvital/core/compare/v1.1.1...v1.1.2) (2025-12-05)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **detection:** fix glob patterns being incorrectly treated as comments ([85fd875](https://github.com/specvital/core/commit/85fd875d706cd1330fd0b8a27f3d1514f36e4013))

## [core/v1.1.1](https://github.com/specvital/core/compare/v1.1.0...v1.1.1) (2025-12-05)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **detection:** add ProjectContext for source-agnostic framework detection ([708f70a](https://github.com/specvital/core/commit/708f70aac041918ea7ff41d698fca45e43d6809d))

## [core/v1.1.0](https://github.com/specvital/core/compare/v1.0.3...v1.1.0) (2025-12-05)

### 🎯 Highlights

#### ✨ Features

- **parser:** add hierarchical test framework detection system ([7655868](https://github.com/specvital/core/commit/76558682788612995f762de422d965f4fa2836ad))

### 🔧 Maintenance

#### 🔨 Chore

- add useful action buttons ([eb2b93b](https://github.com/specvital/core/commit/eb2b93b8e163c2e538a025cff0e35abad891a87b))
- delete unused file ([d6f2203](https://github.com/specvital/core/commit/d6f220316bd8e66423366f61a153627aa0daa7bd))
- syncing documents from ai-config-toolkit ([1faaf43](https://github.com/specvital/core/commit/1faaf4364d1782493008d8abbae66283d35861af))

## [core/v1.0.3](https://github.com/specvital/core/compare/v1.0.2...v1.0.3) (2025-12-04)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **gotesting:** fix Go test parser incorrectly returning pending status ([14f1336](https://github.com/specvital/core/commit/14f133635410d9ced0d747d7245238e84f6014c9))

### 🔧 Maintenance

#### 📚 Documentation

- sync CLAUDE.md ([167df5b](https://github.com/specvital/core/commit/167df5b587fbbaa8b6ade0dbb4c0ecc0ea41fb98))

#### 🔨 Chore

- add auto-formatting to semantic-release pipeline ([f185576](https://github.com/specvital/core/commit/f185576d2247234c46ec1c0027c8898a775ef5cd))

## [core/v1.0.2](https://github.com/specvital/core/compare/v1.0.1...v1.0.2) (2025-12-04)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix Go module zip creation failure ([3ceb7d6](https://github.com/specvital/core/commit/3ceb7d626ead57835083b0c45d2c7091cb62757f))

## [core/v1.0.1](https://github.com/specvital/core/compare/v1.0.0...v1.0.1) (2025-12-04)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- exclude unnecessary files from Go module zip ([0e3f8fa](https://github.com/specvital/core/commit/0e3f8fa9598ce226632139c2b18dd4d710ad79af))

## [core/v1.0.0](https://github.com/specvital/core/releases/tag/v1.0.0) (2025-12-04)

### 🎯 Highlights

#### ✨ Features

- add Go test parser support ([3e147a5](https://github.com/specvital/core/commit/3e147a59b2ec6799db588702a648fd25bb3d44c0))
- add parallel test file scanner with worker pool ([d8dbe13](https://github.com/specvital/core/commit/d8dbe13cc5095a4c2385add15c320c1f9148f76d))
- add Playwright test parser support ([c779d70](https://github.com/specvital/core/commit/c779d7063fdc58e60b085e9daf21b2a8453db7b0))
- add test file detector for automatic discovery ([a71bec4](https://github.com/specvital/core/commit/a71bec4e61c6a05406b9021e6ebd929dce4fff05))
- add Vitest test parser support ([d4226f5](https://github.com/specvital/core/commit/d4226f5238edb8131074fe22cd54d492eae70a94))
- implement Jest test parser core module ([caffaab](https://github.com/specvital/core/commit/caffaab77d810283a126266ed806f4bb1bdc2a0a))

#### ⚡ Performance

- add parser pooling and query caching for concurrent parsing ([e8ff8f4](https://github.com/specvital/core/commit/e8ff8f40ddecd3143d56ecc97c78075e112806cd))

### 🔧 Maintenance

#### 📚 Documentation

- add GoDoc comments and library usage guide ([72f5220](https://github.com/specvital/core/commit/72f5220e7ab96b1497fedfe4f59230774cefe369))

#### ♻️ Refactoring

- move go.mod to root to enable external imports ([8976869](https://github.com/specvital/core/commit/89768699151849542582997faa26ec9d6557e923))
- move src/pkg to pkg for standard Go layout ([3ed1d78](https://github.com/specvital/core/commit/3ed1d782be5b6e3d2fcfbeb45aa003ad48e2eb10))

#### 🔧 CI/CD

- configure semantic-release based automated release pipeline ([3e85cee](https://github.com/specvital/core/commit/3e85ceeb26ca91009c7c76dc71108ef985ea9538))

#### 🔨 Chore

- **deps-dev:** bump lint-staged from 15.2.11 to 16.2.7 ([94b8012](https://github.com/specvital/core/commit/94b801204734591d3a8aaece07562ae0423354b7))
