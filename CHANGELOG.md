# Changelog

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
