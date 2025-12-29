import { defineConfig } from "vitepress";

export default defineConfig({
  title: "Specvital",
  description: "Open-source test coverage insights for code review",

  head: [
    ["link", { rel: "icon", href: "/favicon.ico" }],
    ["meta", { property: "og:image", content: "https://specvital.github.io/og.png" }],
  ],

  vite: {
    publicDir: ".vitepress/public",
  },

  // Clean URLs without .html
  cleanUrls: true,

  // Ignore dead links in template files
  ignoreDeadLinks: [/XX-decision-title/],

  // Multi-language configuration
  locales: {
    en: {
      label: "English",
      lang: "en",
      link: "/en/",
      themeConfig: {
        nav: [
          { text: "Home", link: "/en/" },
          { text: "PRD", link: "/en/prd/" },
          { text: "ADR", link: "/en/adr/" },
        ],
        sidebar: {
          "/en/": [
            {
              text: "Overview",
              items: [
                { text: "Home", link: "/en/" },
                { text: "Architecture", link: "/en/architecture" },
                { text: "Tech Radar", link: "/en/tech-radar" },
                { text: "Releases", link: "/en/releases" },
                { text: "Glossary", link: "/en/glossary" },
              ],
            },
            {
              text: "PRD",
              collapsed: false,
              items: [
                { text: "Overview", link: "/en/prd/" },
                { text: "Product Overview", link: "/en/prd/00-overview" },
                { text: "Architecture", link: "/en/prd/01-architecture" },
                { text: "Core Engine", link: "/en/prd/02-core-engine" },
                { text: "Web Platform", link: "/en/prd/03-web-platform" },
                { text: "Collector Service", link: "/en/prd/04-collector-service" },
                { text: "Database Design", link: "/en/prd/05-database-design" },
                { text: "Tech Stack", link: "/en/prd/06-tech-stack" },
              ],
            },
            {
              text: "ADR",
              collapsed: false,
              items: [
                { text: "Overview", link: "/en/adr/" },
                { text: "Static Analysis Approach", link: "/en/adr/01-static-analysis-approach" },
                {
                  text: "Competitive Differentiation",
                  link: "/en/adr/02-competitive-differentiation",
                },
                { text: "API Worker Separation", link: "/en/adr/03-api-worker-service-separation" },
                { text: "Queue-Based Processing", link: "/en/adr/04-queue-based-async-processing" },
                { text: "Repository Strategy", link: "/en/adr/05-repository-strategy" },
                { text: "PaaS-First Infrastructure", link: "/en/adr/06-paas-first-infrastructure" },
                { text: "Shared Infrastructure", link: "/en/adr/07-shared-infrastructure" },
                {
                  text: "External Repo ID Integrity",
                  link: "/en/adr/08-external-repo-id-integrity",
                },
                {
                  text: "Core",
                  collapsed: true,
                  items: [
                    { text: "Overview", link: "/en/adr/core/" },
                    { text: "Core Library Separation", link: "/en/adr/core/01-core-library-separation" },
                    { text: "Dynamic Test Counting Policy", link: "/en/adr/core/02-dynamic-test-counting-policy" },
                    { text: "Tree-sitter AST Parsing", link: "/en/adr/core/03-tree-sitter-ast-parsing-engine" },
                    { text: "Early-Return Detection", link: "/en/adr/core/04-early-return-framework-detection" },
                    { text: "Parser Pooling Disabled", link: "/en/adr/core/05-parser-pooling-disabled" },
                    { text: "Unified Framework Definition", link: "/en/adr/core/06-unified-framework-definition" },
                    { text: "Source Abstraction Interface", link: "/en/adr/core/07-source-abstraction-interface" },
                    { text: "Shared Parser Modules", link: "/en/adr/core/08-shared-parser-modules" },
                    { text: "Config Scope Resolution", link: "/en/adr/core/09-config-scope-resolution" },
                    { text: "Standard Go Project Layout", link: "/en/adr/core/10-standard-go-project-layout" },
                    { text: "Golden Snapshot Testing", link: "/en/adr/core/11-integration-testing-golden-snapshots" },
                    { text: "Parallel Scanning", link: "/en/adr/core/12-parallel-scanning-worker-pool" },
                    { text: "NaCl SecretBox Encryption", link: "/en/adr/core/13-nacl-secretbox-encryption" },
                    { text: "Indirect Import Unsupported", link: "/en/adr/core/14-indirect-import-unsupported" },
                  ],
                },
                {
                  text: "Collector",
                  collapsed: true,
                  items: [
                    { text: "Overview", link: "/en/adr/collector/" },
                    {
                      text: "Scheduled Re-collection",
                      link: "/en/adr/collector/01-scheduled-recollection",
                    },
                    {
                      text: "Clean Architecture Layers",
                      link: "/en/adr/collector/02-clean-architecture-layers",
                    },
                    { text: "Graceful Shutdown", link: "/en/adr/collector/03-graceful-shutdown" },
                    {
                      text: "OAuth Token Degradation",
                      link: "/en/adr/collector/04-oauth-token-graceful-degradation",
                    },
                    {
                      text: "Worker-Scheduler Separation",
                      link: "/en/adr/collector/05-worker-scheduler-separation",
                    },
                    {
                      text: "Semaphore Clone Concurrency",
                      link: "/en/adr/collector/06-semaphore-clone-concurrency",
                    },
                    { text: "Repository Pattern", link: "/en/adr/collector/07-repository-pattern" },
                  ],
                },
                {
                  text: "Web",
                  collapsed: true,
                  items: [
                    { text: "Overview", link: "/en/adr/web/" },
                    { text: "Go Backend Language", link: "/en/adr/web/01-go-backend-language" },
                  ],
                },
              ],
            },
          ],
        },
      },
    },
    ko: {
      label: "한국어",
      lang: "ko",
      link: "/ko/",
      themeConfig: {
        nav: [
          { text: "홈", link: "/ko/" },
          { text: "PRD", link: "/ko/prd/" },
          { text: "ADR", link: "/ko/adr/" },
        ],
        sidebar: {
          "/ko/": [
            {
              text: "개요",
              items: [
                { text: "홈", link: "/ko/" },
                { text: "아키텍처", link: "/ko/architecture" },
                { text: "기술 레이더", link: "/ko/tech-radar" },
                { text: "릴리즈 노트", link: "/ko/releases" },
                { text: "용어집", link: "/ko/glossary" },
              ],
            },
            {
              text: "PRD",
              collapsed: false,
              items: [
                { text: "개요", link: "/ko/prd/" },
                { text: "제품 개요", link: "/ko/prd/00-overview" },
                { text: "아키텍처", link: "/ko/prd/01-architecture" },
                { text: "코어 엔진", link: "/ko/prd/02-core-engine" },
                { text: "웹 플랫폼", link: "/ko/prd/03-web-platform" },
                { text: "컬렉터 서비스", link: "/ko/prd/04-collector-service" },
                { text: "데이터베이스 설계", link: "/ko/prd/05-database-design" },
                { text: "기술 스택", link: "/ko/prd/06-tech-stack" },
              ],
            },
            {
              text: "ADR",
              collapsed: false,
              items: [
                { text: "개요", link: "/ko/adr/" },
                { text: "정적 분석 접근법", link: "/ko/adr/01-static-analysis-approach" },
                { text: "경쟁 차별화", link: "/ko/adr/02-competitive-differentiation" },
                { text: "API/Worker 분리", link: "/ko/adr/03-api-worker-service-separation" },
                { text: "큐 기반 비동기 처리", link: "/ko/adr/04-queue-based-async-processing" },
                { text: "리포지토리 전략", link: "/ko/adr/05-repository-strategy" },
                { text: "PaaS 우선 인프라", link: "/ko/adr/06-paas-first-infrastructure" },
                { text: "공유 인프라", link: "/ko/adr/07-shared-infrastructure" },
                { text: "External Repo ID 무결성", link: "/ko/adr/08-external-repo-id-integrity" },
                {
                  text: "Core",
                  collapsed: true,
                  items: [
                    { text: "개요", link: "/ko/adr/core/" },
                    { text: "코어 라이브러리 분리", link: "/ko/adr/core/01-core-library-separation" },
                    { text: "동적 테스트 카운팅 정책", link: "/ko/adr/core/02-dynamic-test-counting-policy" },
                    { text: "Tree-sitter AST 파싱", link: "/ko/adr/core/03-tree-sitter-ast-parsing-engine" },
                    { text: "Early-Return 탐지", link: "/ko/adr/core/04-early-return-framework-detection" },
                    { text: "파서 풀링 비활성화", link: "/ko/adr/core/05-parser-pooling-disabled" },
                    { text: "통합 Framework Definition", link: "/ko/adr/core/06-unified-framework-definition" },
                    { text: "Source 추상화 인터페이스", link: "/ko/adr/core/07-source-abstraction-interface" },
                    { text: "공유 파서 모듈", link: "/ko/adr/core/08-shared-parser-modules" },
                    { text: "Config 스코프 해석", link: "/ko/adr/core/09-config-scope-resolution" },
                    { text: "표준 Go 프로젝트 레이아웃", link: "/ko/adr/core/10-standard-go-project-layout" },
                    { text: "골든 스냅샷 테스트", link: "/ko/adr/core/11-integration-testing-golden-snapshots" },
                    { text: "병렬 스캐닝", link: "/ko/adr/core/12-parallel-scanning-worker-pool" },
                    { text: "NaCl SecretBox 암호화", link: "/ko/adr/core/13-nacl-secretbox-encryption" },
                    { text: "간접 Import 미지원", link: "/ko/adr/core/14-indirect-import-unsupported" },
                  ],
                },
                {
                  text: "Collector",
                  collapsed: true,
                  items: [
                    { text: "개요", link: "/ko/adr/collector/" },
                    {
                      text: "스케줄 기반 재수집",
                      link: "/ko/adr/collector/01-scheduled-recollection",
                    },
                    {
                      text: "Clean Architecture 레이어",
                      link: "/ko/adr/collector/02-clean-architecture-layers",
                    },
                    { text: "Graceful Shutdown", link: "/ko/adr/collector/03-graceful-shutdown" },
                    {
                      text: "OAuth 토큰 Degradation",
                      link: "/ko/adr/collector/04-oauth-token-graceful-degradation",
                    },
                    {
                      text: "Worker-Scheduler 분리",
                      link: "/ko/adr/collector/05-worker-scheduler-separation",
                    },
                    {
                      text: "Semaphore Clone 동시성",
                      link: "/ko/adr/collector/06-semaphore-clone-concurrency",
                    },
                    { text: "Repository 패턴", link: "/ko/adr/collector/07-repository-pattern" },
                  ],
                },
                {
                  text: "Web",
                  collapsed: true,
                  items: [
                    { text: "개요", link: "/ko/adr/web/" },
                    { text: "Go 백엔드 언어", link: "/ko/adr/web/01-go-backend-language" },
                  ],
                },
              ],
            },
          ],
        },
      },
    },
  },

  themeConfig: {
    logo: "/logo.png",

    socialLinks: [{ icon: "github", link: "https://github.com/specvital" }],

    search: {
      provider: "local",
    },

    footer: {
      message: "Open-source test coverage insights",
      copyright: "Copyright © 2024 Specvital",
    },

    editLink: {
      pattern: "https://github.com/specvital/specvital.github.io/edit/main/docs/:path",
      text: "Edit this page on GitHub",
    },
  },
});
