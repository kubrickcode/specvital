# Changelog

## [1.6.3](https://github.com/specvital/web/compare/v1.6.2...v1.6.3) (2026-02-09)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- **infra:** set Railway deploy region to US East ([db548d2](https://github.com/specvital/web/commit/db548d204eb6f32bc7f316c019e6dc8bf38bcdf6))

## [1.6.2](https://github.com/specvital/web/compare/v1.6.1...v1.6.2) (2026-02-05)

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

## [1.6.1](https://github.com/specvital/web/compare/v1.6.0...v1.6.1) (2026-02-04)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **i18n:** apply missing i18n to toast messages ([9be4e2f](https://github.com/specvital/web/commit/9be4e2f88e203a3e80c25adbf9bc0b694e2f6590))
- **ui:** fix loading spinner appearing below visual center ([516cf00](https://github.com/specvital/web/commit/516cf005070cf3ca92144d1d04b2923d49e965e7))

## [1.6.0](https://github.com/specvital/web/compare/v1.5.0...v1.6.0) (2026-02-03)

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

## [1.5.0](https://github.com/specvital/web/compare/v1.4.1...v1.5.0) (2026-02-02)

### 🎯 Highlights

#### ✨ Features

- **spec-view:** group version history by commit SHA ([42b852b](https://github.com/specvital/web/commit/42b852b0662972ccd61459cc601e2d778aa0538d))

#### 🐛 Bug Fixes

- **spec-view:** fix multiple documents showing as "latest" in version history ([b217af4](https://github.com/specvital/web/commit/b217af4486ccf86224c0bcacc6293dd3d241f095))

## [1.4.1](https://github.com/specvital/web/compare/v1.4.0...v1.4.1) (2026-02-02)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix release failure due to semantic-release plugin version mismatch ([0426f2d](https://github.com/specvital/web/commit/0426f2dee52237c2a8063ebd9d32379722cc8b39))

## [1.4.0](https://github.com/specvital/web/compare/v1.3.1...v1.4.0) (2026-02-02)

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

## [1.3.1](https://github.com/specvital/web/compare/v1.3.0...v1.3.1) (2026-01-07)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **ui:** fix AI Analysis button overflow on mobile ([e530add](https://github.com/specvital/web/commit/e530add96bde80f7ea3cb8bd1ae2e600f461fcbd))

### 🔧 Maintenance

#### ✅ Tests

- **spec-view:** skip handler tests during feature suspension ([df074c8](https://github.com/specvital/web/commit/df074c88bdbd2a626422a4fbf10ad091c7b9ad7a))

## [1.3.0](https://github.com/specvital/web/compare/v1.2.1...v1.3.0) (2026-01-07)

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

## [1.2.1](https://github.com/specvital/web/compare/v1.2.0...v1.2.1) (2026-01-04)

### 🎯 Highlights

#### 🐛 Bug Fixes

- **auth:** fix dashboard redirect bug after time passes while logged in ([ef604fd](https://github.com/specvital/web/commit/ef604fdc41724ae42e62c9e57566712fb191d399))
- **routing:** fix 404 error for paths containing dots in repository names ([cfb36fc](https://github.com/specvital/web/commit/cfb36fc25ace5fff6303edc73ce72e6c4f3e8811))

### 🔧 Maintenance

#### 💄 Styles

- replace Korean comments with English ([796db0c](https://github.com/specvital/web/commit/796db0c5df3cd223a43a688c899e7e6fa9827e26))

## [1.2.0](https://github.com/specvital/web/compare/v1.1.2...v1.2.0) (2026-01-04)

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

## [1.1.2](https://github.com/specvital/web/compare/v1.1.1...v1.1.2) (2025-12-20)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix semantic-release CHANGELOG version header not rendering ([372a8ef](https://github.com/specvital/web/commit/372a8efc305de5267f727dfda02a3b285cbab22b))

## [1.1.1](https://github.com/specvital/web/compare/v1.1.0...v1.1.1) (2025-12-20)

### 🔧 Maintenance

#### 🔧 Internal Fixes

- fix semantic-release CHANGELOG version header not rendering ([5bff792](https://github.com/specvital/web/commit/5bff792d9dcca96eb0d08c7d2347a7940a784506))

## [1.1.0](https://github.com/specvital/web/compare/v1.0.4...v1.1.0) (2025-12-20)

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

## [1.0.4](https://github.com/specvital/web/compare/v1.0.3...v1.0.4) (2025-12-19)

### 🔧 Maintenance

#### ♻️ Refactoring

- migrate job queue from asynq to river ([72fce89](https://github.com/specvital/web/commit/72fce895b4cff07bef68244d7be08be59348b660))

## [1.0.3](https://github.com/specvital/web/compare/v1.0.2...v1.0.3) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- cookie not being set after GitHub login ([f4fccee](https://github.com/specvital/web/commit/f4fccee642db5089681837f66af66ee3b92a8e68))

## [1.0.2](https://github.com/specvital/web/compare/v1.0.1...v1.0.2) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- "failed to get latest commit" error during repository analysis ([0de5c39](https://github.com/specvital/web/commit/0de5c399abe3d02435c81640c50d43d1a5bfa37f))

## [1.0.1](https://github.com/specvital/web/compare/v1.0.0...v1.0.1) (2025-12-18)

### 🎯 Highlights

#### 🐛 Bug Fixes

- page not working in production environment ([21a60f7](https://github.com/specvital/web/commit/21a60f7700180cbe01faef41458cc5b73be645d0))

## [1.0.0](https://github.com/specvital/web/commits/v1.0.0) (2025-12-18)

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
