# Changelog

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
