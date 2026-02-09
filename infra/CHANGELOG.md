# Changelog

## [1.3.0](https://github.com/specvital/infra/compare/v1.2.0...v1.3.0) (2026-02-02)

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

## [1.2.0](https://github.com/specvital/infra/compare/v1.1.0...v1.2.0) (2026-01-04)

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

## [1.1.0](https://github.com/specvital/infra/compare/v1.0.0...v1.1.0) (2025-12-19)

### 🎯 Highlights

#### ✨ Features

- **db:** add River job queue migration ([86b6157](https://github.com/specvital/infra/commit/86b61576794e3df0a097f151e67afad9f38c2abc))

### 🔧 Maintenance

#### 🔨 Chore

- adding a go environment to a container for riverqueue use ([ee45552](https://github.com/specvital/infra/commit/ee45552c4d80fd457c61df5f31c110534d4a0f7f))
- remove Redis dependency ([916c622](https://github.com/specvital/infra/commit/916c6227d3646e6d8baad8efe8e663e3f50b525b))

## [1.0.0](https://github.com/specvital/infra/releases/tag/v1.0.0) (2025-12-17)

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
