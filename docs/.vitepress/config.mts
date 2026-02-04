import { withMermaid } from "vitepress-plugin-mermaid";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

function getSchemaSidebar() {
  const schemaDir = path.resolve(__dirname, "../schema");
  const files = fs
    .readdirSync(schemaDir)
    .filter((f) => f.endsWith(".md") && f !== "README.md")
    .sort();

  const publicTables = files
    .filter((f) => f.startsWith("public."))
    .map((f) => ({
      text: f.replace("public.", "").replace(".md", ""),
      link: `/schema/${f.replace(".md", "")}`,
    }));

  const internalTables = files
    .filter((f) => !f.startsWith("public."))
    .map((f) => ({
      text: f.replace(".md", ""),
      link: `/schema/${f.replace(".md", "")}`,
    }));

  return [
    {
      text: "Database Schema",
      items: [
        { text: "Overview", link: "/schema/" },
        { text: "Tables", collapsed: false, items: publicTables },
        ...(internalTables.length > 0
          ? [{ text: "Internal", collapsed: true, items: internalTables }]
          : []),
      ],
    },
  ];
}

export default withMermaid({
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

  // Map README.md to index for tbls compatibility
  rewrites: {
    "schema/README.md": "schema/index.md",
  },

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
          { text: "Schema", link: "/schema/" },
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
                { text: "Worker Service", link: "/en/prd/04-worker-service" },
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
                  text: "GitHub App Integration",
                  link: "/en/adr/09-github-app-integration",
                },
                {
                  text: "TestStatus Data Contract",
                  link: "/en/adr/10-test-status-data-contract",
                },
                {
                  text: "Visibility Access Control",
                  link: "/en/adr/11-community-private-repo-filtering",
                },
                {
                  text: "Worker-Centric Lifecycle",
                  link: "/en/adr/12-worker-centric-analysis-lifecycle",
                },
                {
                  text: "Billing and Quota",
                  link: "/en/adr/13-billing-quota-architecture",
                },
                {
                  text: "AI Spec Generation Pipeline",
                  link: "/en/adr/14-ai-spec-generation-pipeline",
                },
                {
                  text: "Parser Version Tracking",
                  link: "/en/adr/15-parser-version-tracking",
                },
                {
                  text: "Multi-Queue Priority Routing",
                  link: "/en/adr/16-multi-queue-priority-routing",
                },
                {
                  text: "Test File Schema Normalization",
                  link: "/en/adr/17-test-file-schema-normalization",
                },
                {
                  text: "GitHub API Cache Tables",
                  link: "/en/adr/18-github-api-cache-tables",
                },
                {
                  text: "Hierarchical Spec Document Schema",
                  link: "/en/adr/19-hierarchical-spec-document-schema",
                },
                {
                  text: "GitHub App Installation Schema",
                  link: "/en/adr/20-github-app-installation-schema",
                },
                {
                  text: "Quota Reservation",
                  link: "/en/adr/21-quota-reservation",
                },
                {
                  text: "Scheduler Removal",
                  link: "/en/adr/22-scheduler-removal-railway-cron",
                },
                {
                  text: "Railway IaC Deployment",
                  link: "/en/adr/23-railway-iac-deployment",
                },
                {
                  text: "Core",
                  collapsed: true,
                  items: [
                    { text: "Overview", link: "/en/adr/core/" },
                    {
                      text: "Core Library Separation",
                      link: "/en/adr/core/01-core-library-separation",
                    },
                    {
                      text: "Dynamic Test Counting Policy",
                      link: "/en/adr/core/02-dynamic-test-counting-policy",
                    },
                    {
                      text: "Tree-sitter AST Parsing",
                      link: "/en/adr/core/03-tree-sitter-ast-parsing-engine",
                    },
                    {
                      text: "Early-Return Detection",
                      link: "/en/adr/core/04-early-return-framework-detection",
                    },
                    {
                      text: "Parser Pooling Disabled",
                      link: "/en/adr/core/05-parser-pooling-disabled",
                    },
                    {
                      text: "Unified Framework Definition",
                      link: "/en/adr/core/06-unified-framework-definition",
                    },
                    {
                      text: "Source Abstraction Interface",
                      link: "/en/adr/core/07-source-abstraction-interface",
                    },
                    {
                      text: "Shared Parser Modules",
                      link: "/en/adr/core/08-shared-parser-modules",
                    },
                    {
                      text: "Config Scope Resolution",
                      link: "/en/adr/core/09-config-scope-resolution",
                    },
                    {
                      text: "Standard Go Project Layout",
                      link: "/en/adr/core/10-standard-go-project-layout",
                    },
                    {
                      text: "Golden Snapshot Testing",
                      link: "/en/adr/core/11-integration-testing-golden-snapshots",
                    },
                    {
                      text: "Parallel Scanning",
                      link: "/en/adr/core/12-parallel-scanning-worker-pool",
                    },
                    {
                      text: "NaCl SecretBox Encryption",
                      link: "/en/adr/core/13-nacl-secretbox-encryption",
                    },
                    {
                      text: "Indirect Import Unsupported",
                      link: "/en/adr/core/14-indirect-import-unsupported",
                    },
                    {
                      text: "C# Preprocessor Limitation",
                      link: "/en/adr/core/15-csharp-preprocessor-attribute-limitation",
                    },
                    {
                      text: "Domain Hints Extraction",
                      link: "/en/adr/core/16-domain-hints-extraction",
                    },
                    {
                      text: "Swift Testing Support",
                      link: "/en/adr/core/17-swift-testing-framework-support",
                    },
                    {
                      text: "JUnit 4 Separation",
                      link: "/en/adr/core/18-junit4-framework-separation",
                    },
                    {
                      text: "Vitest 4.0+ API Support",
                      link: "/en/adr/core/19-vitest-4-api-support",
                    },
                    {
                      text: "Java 21+ Implicit Class",
                      link: "/en/adr/core/20-java21-implicit-class-detection",
                    },
                    {
                      text: "Rust Macro Test Detection",
                      link: "/en/adr/core/21-rust-macro-test-detection",
                    },
                  ],
                },
                {
                  text: "Worker",
                  collapsed: true,
                  items: [
                    { text: "Overview", link: "/en/adr/worker/" },
                    {
                      text: "Scheduled Re-analysis",
                      link: "/en/adr/worker/01-scheduled-recollection",
                    },
                    {
                      text: "Clean Architecture Layers",
                      link: "/en/adr/worker/02-clean-architecture-layers",
                    },
                    { text: "Graceful Shutdown", link: "/en/adr/worker/03-graceful-shutdown" },
                    {
                      text: "OAuth Token Degradation",
                      link: "/en/adr/worker/04-oauth-token-graceful-degradation",
                    },
                    {
                      text: "Analyzer-Scheduler Separation",
                      link: "/en/adr/worker/05-worker-scheduler-separation",
                    },
                    {
                      text: "Semaphore Clone Concurrency",
                      link: "/en/adr/worker/06-semaphore-clone-concurrency",
                    },
                    { text: "Repository Pattern", link: "/en/adr/worker/07-repository-pattern" },
                    {
                      text: "SpecView Worker Separation",
                      link: "/en/adr/worker/08-specview-worker-separation",
                    },
                    {
                      text: "Phase 1 V1 Quality Enhancement",
                      link: "/en/adr/worker/09-phase1-v1-quality-enhancement",
                    },
                    {
                      text: "Per-User Fairness Middleware",
                      link: "/en/adr/worker/10-per-user-fairness-middleware",
                    },
                    {
                      text: "Data Retention Cleanup",
                      link: "/en/adr/worker/11-data-retention-cleanup",
                    },
                    {
                      text: "Phase 2 Behavior Cache",
                      link: "/en/adr/worker/12-phase2-behavior-cache",
                    },
                  ],
                },
                {
                  text: "Web",
                  collapsed: true,
                  items: [
                    { text: "Overview", link: "/en/adr/web/" },
                    { text: "Go Backend Language", link: "/en/adr/web/01-go-backend-language" },
                    {
                      text: "Next.js + React Selection",
                      link: "/en/adr/web/02-nextjs-react-selection",
                    },
                    { text: "Chi Router Selection", link: "/en/adr/web/03-chi-router-selection" },
                    {
                      text: "TanStack Query Selection",
                      link: "/en/adr/web/04-tanstack-query-selection",
                    },
                    {
                      text: "shadcn/ui + Tailwind",
                      link: "/en/adr/web/05-shadcn-tailwind-selection",
                    },
                    { text: "SQLc Selection", link: "/en/adr/web/06-sqlc-selection" },
                    {
                      text: "Next.js BFF Architecture",
                      link: "/en/adr/web/07-nextjs-bff-architecture",
                    },
                    {
                      text: "Clean Architecture",
                      link: "/en/adr/web/08-clean-architecture-pattern",
                    },
                    { text: "DI Container Pattern", link: "/en/adr/web/09-di-container-pattern" },
                    {
                      text: "StrictServerInterface",
                      link: "/en/adr/web/10-strict-server-interface-contract",
                    },
                    {
                      text: "Feature-Based Modules",
                      link: "/en/adr/web/11-feature-based-module-organization",
                    },
                    {
                      text: "APIHandlers Composition",
                      link: "/en/adr/web/12-apihandlers-composition-pattern",
                    },
                    {
                      text: "Domain Error Handling",
                      link: "/en/adr/web/13-domain-error-handling-pattern",
                    },
                    {
                      text: "slog Structured Logging",
                      link: "/en/adr/web/14-slog-structured-logging",
                    },
                    {
                      text: "React 19 use() Hook",
                      link: "/en/adr/web/15-react-19-use-hook-pattern",
                    },
                    { text: "nuqs URL State", link: "/en/adr/web/16-nuqs-url-state-management" },
                    { text: "next-intl i18n", link: "/en/adr/web/17-next-intl-i18n-strategy" },
                    { text: "next-themes Dark Mode", link: "/en/adr/web/18-next-themes-dark-mode" },
                    {
                      text: "CSS Variable Tokens",
                      link: "/en/adr/web/19-css-variable-design-token-system",
                    },
                    { text: "Skeleton Loading", link: "/en/adr/web/20-skeleton-loading-pattern" },
                    {
                      text: "Anonymous Rate Limiting",
                      link: "/en/adr/web/21-anonymous-rate-limiting",
                    },
                    {
                      text: "React Compiler Adoption",
                      link: "/en/adr/web/22-react-compiler-adoption",
                    },
                    {
                      text: "Window-Level Virtualization",
                      link: "/en/adr/web/23-window-level-virtualization",
                    },
                    {
                      text: "Subscription Period Pro-rata",
                      link: "/en/adr/web/24-subscription-period-prorata",
                    },
                    {
                      text: "OAuth Return URL Handling",
                      link: "/en/adr/web/25-oauth-return-url-handling",
                    },
                    {
                      text: "Zustand Client State",
                      link: "/en/adr/web/26-zustand-client-state",
                    },
                    {
                      text: "Commit-Based Analysis Browsing",
                      link: "/en/adr/web/27-commit-based-analysis-browsing",
                    },
                    {
                      text: "Dynamic Cost Estimation",
                      link: "/en/adr/web/28-dynamic-cost-estimation",
                    },
                  ],
                },
              ],
            },
            {
              text: "Schema",
              collapsed: true,
              items: getSchemaSidebar()[0].items,
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
          { text: "Schema", link: "/schema/" },
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
                { text: "워커 서비스", link: "/ko/prd/04-worker-service" },
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
                { text: "GitHub App 통합", link: "/ko/adr/09-github-app-integration" },
                { text: "TestStatus 데이터 계약", link: "/ko/adr/10-test-status-data-contract" },
                {
                  text: "Visibility 접근 제어",
                  link: "/ko/adr/11-community-private-repo-filtering",
                },
                {
                  text: "Worker 중심 라이프사이클",
                  link: "/ko/adr/12-worker-centric-analysis-lifecycle",
                },
                { text: "빌링 및 쿼터", link: "/ko/adr/13-billing-quota-architecture" },
                { text: "AI 스펙 생성 파이프라인", link: "/ko/adr/14-ai-spec-generation-pipeline" },
                { text: "파서 버전 추적", link: "/ko/adr/15-parser-version-tracking" },
                { text: "멀티큐 우선순위 라우팅", link: "/ko/adr/16-multi-queue-priority-routing" },
                {
                  text: "테스트 파일 스키마 정규화",
                  link: "/ko/adr/17-test-file-schema-normalization",
                },
                { text: "GitHub API 캐시 테이블", link: "/ko/adr/18-github-api-cache-tables" },
                {
                  text: "계층적 스펙 문서 스키마",
                  link: "/ko/adr/19-hierarchical-spec-document-schema",
                },
                {
                  text: "GitHub App 설치 스키마",
                  link: "/ko/adr/20-github-app-installation-schema",
                },
                {
                  text: "할당량 예약",
                  link: "/ko/adr/21-quota-reservation",
                },
                {
                  text: "Scheduler 제거",
                  link: "/ko/adr/22-scheduler-removal-railway-cron",
                },
                {
                  text: "Railway IaC 배포",
                  link: "/ko/adr/23-railway-iac-deployment",
                },
                {
                  text: "Core",
                  collapsed: true,
                  items: [
                    { text: "개요", link: "/ko/adr/core/" },
                    {
                      text: "코어 라이브러리 분리",
                      link: "/ko/adr/core/01-core-library-separation",
                    },
                    {
                      text: "동적 테스트 카운팅 정책",
                      link: "/ko/adr/core/02-dynamic-test-counting-policy",
                    },
                    {
                      text: "Tree-sitter AST 파싱",
                      link: "/ko/adr/core/03-tree-sitter-ast-parsing-engine",
                    },
                    {
                      text: "Early-Return 탐지",
                      link: "/ko/adr/core/04-early-return-framework-detection",
                    },
                    { text: "파서 풀링 비활성화", link: "/ko/adr/core/05-parser-pooling-disabled" },
                    {
                      text: "통합 Framework Definition",
                      link: "/ko/adr/core/06-unified-framework-definition",
                    },
                    {
                      text: "Source 추상화 인터페이스",
                      link: "/ko/adr/core/07-source-abstraction-interface",
                    },
                    { text: "공유 파서 모듈", link: "/ko/adr/core/08-shared-parser-modules" },
                    { text: "Config 스코프 해석", link: "/ko/adr/core/09-config-scope-resolution" },
                    {
                      text: "표준 Go 프로젝트 레이아웃",
                      link: "/ko/adr/core/10-standard-go-project-layout",
                    },
                    {
                      text: "골든 스냅샷 테스트",
                      link: "/ko/adr/core/11-integration-testing-golden-snapshots",
                    },
                    { text: "병렬 스캐닝", link: "/ko/adr/core/12-parallel-scanning-worker-pool" },
                    {
                      text: "NaCl SecretBox 암호화",
                      link: "/ko/adr/core/13-nacl-secretbox-encryption",
                    },
                    {
                      text: "간접 Import 미지원",
                      link: "/ko/adr/core/14-indirect-import-unsupported",
                    },
                    {
                      text: "C# 전처리기 한계",
                      link: "/ko/adr/core/15-csharp-preprocessor-attribute-limitation",
                    },
                    { text: "도메인 힌트 추출", link: "/ko/adr/core/16-domain-hints-extraction" },
                    {
                      text: "Swift Testing 지원",
                      link: "/ko/adr/core/17-swift-testing-framework-support",
                    },
                    { text: "JUnit 4 분리", link: "/ko/adr/core/18-junit4-framework-separation" },
                    { text: "Vitest 4.0+ API 지원", link: "/ko/adr/core/19-vitest-4-api-support" },
                    {
                      text: "Java 21+ 암시적 클래스",
                      link: "/ko/adr/core/20-java21-implicit-class-detection",
                    },
                    {
                      text: "Rust 매크로 테스트 감지",
                      link: "/ko/adr/core/21-rust-macro-test-detection",
                    },
                  ],
                },
                {
                  text: "Worker",
                  collapsed: true,
                  items: [
                    { text: "개요", link: "/ko/adr/worker/" },
                    {
                      text: "스케줄 기반 재분석",
                      link: "/ko/adr/worker/01-scheduled-recollection",
                    },
                    {
                      text: "Clean Architecture 레이어",
                      link: "/ko/adr/worker/02-clean-architecture-layers",
                    },
                    { text: "Graceful Shutdown", link: "/ko/adr/worker/03-graceful-shutdown" },
                    {
                      text: "OAuth 토큰 Degradation",
                      link: "/ko/adr/worker/04-oauth-token-graceful-degradation",
                    },
                    {
                      text: "Analyzer-Scheduler 분리",
                      link: "/ko/adr/worker/05-worker-scheduler-separation",
                    },
                    {
                      text: "Semaphore Clone 동시성",
                      link: "/ko/adr/worker/06-semaphore-clone-concurrency",
                    },
                    { text: "Repository 패턴", link: "/ko/adr/worker/07-repository-pattern" },
                    {
                      text: "SpecView Worker 분리",
                      link: "/ko/adr/worker/08-specview-worker-separation",
                    },
                    {
                      text: "Phase 1 V1 품질 개선",
                      link: "/ko/adr/worker/09-phase1-v1-quality-enhancement",
                    },
                    {
                      text: "사용자별 공정성 미들웨어",
                      link: "/ko/adr/worker/10-per-user-fairness-middleware",
                    },
                    {
                      text: "보존 기간 기반 데이터 정리",
                      link: "/ko/adr/worker/11-data-retention-cleanup",
                    },
                    {
                      text: "Phase 2 행동 캐시",
                      link: "/ko/adr/worker/12-phase2-behavior-cache",
                    },
                  ],
                },
                {
                  text: "Web",
                  collapsed: true,
                  items: [
                    { text: "개요", link: "/ko/adr/web/" },
                    { text: "Go 백엔드 언어", link: "/ko/adr/web/01-go-backend-language" },
                    { text: "Next.js + React 선택", link: "/ko/adr/web/02-nextjs-react-selection" },
                    { text: "Chi 라우터 선택", link: "/ko/adr/web/03-chi-router-selection" },
                    {
                      text: "TanStack Query 선택",
                      link: "/ko/adr/web/04-tanstack-query-selection",
                    },
                    {
                      text: "shadcn/ui + Tailwind",
                      link: "/ko/adr/web/05-shadcn-tailwind-selection",
                    },
                    { text: "SQLc 선택", link: "/ko/adr/web/06-sqlc-selection" },
                    {
                      text: "Next.js BFF 아키텍처",
                      link: "/ko/adr/web/07-nextjs-bff-architecture",
                    },
                    {
                      text: "Clean Architecture",
                      link: "/ko/adr/web/08-clean-architecture-pattern",
                    },
                    { text: "DI Container 패턴", link: "/ko/adr/web/09-di-container-pattern" },
                    {
                      text: "StrictServerInterface",
                      link: "/ko/adr/web/10-strict-server-interface-contract",
                    },
                    {
                      text: "Feature 기반 모듈",
                      link: "/ko/adr/web/11-feature-based-module-organization",
                    },
                    {
                      text: "APIHandlers 합성",
                      link: "/ko/adr/web/12-apihandlers-composition-pattern",
                    },
                    {
                      text: "도메인 에러 처리",
                      link: "/ko/adr/web/13-domain-error-handling-pattern",
                    },
                    { text: "slog 구조화 로깅", link: "/ko/adr/web/14-slog-structured-logging" },
                    {
                      text: "React 19 use() Hook",
                      link: "/ko/adr/web/15-react-19-use-hook-pattern",
                    },
                    { text: "nuqs URL 상태", link: "/ko/adr/web/16-nuqs-url-state-management" },
                    { text: "next-intl i18n", link: "/ko/adr/web/17-next-intl-i18n-strategy" },
                    { text: "next-themes 다크 모드", link: "/ko/adr/web/18-next-themes-dark-mode" },
                    {
                      text: "CSS 변수 토큰",
                      link: "/ko/adr/web/19-css-variable-design-token-system",
                    },
                    { text: "스켈레톤 로딩", link: "/ko/adr/web/20-skeleton-loading-pattern" },
                    { text: "익명 Rate Limiting", link: "/ko/adr/web/21-anonymous-rate-limiting" },
                    { text: "React Compiler 도입", link: "/ko/adr/web/22-react-compiler-adoption" },
                    {
                      text: "윈도우 레벨 가상화",
                      link: "/ko/adr/web/23-window-level-virtualization",
                    },
                    {
                      text: "구독 기간 일할 계산",
                      link: "/ko/adr/web/24-subscription-period-prorata",
                    },
                    {
                      text: "OAuth 리턴 URL 처리",
                      link: "/ko/adr/web/25-oauth-return-url-handling",
                    },
                    {
                      text: "Zustand 클라이언트 상태",
                      link: "/ko/adr/web/26-zustand-client-state",
                    },
                    {
                      text: "커밋 기반 분석 조회",
                      link: "/ko/adr/web/27-commit-based-analysis-browsing",
                    },
                    {
                      text: "동적 비용 예측",
                      link: "/ko/adr/web/28-dynamic-cost-estimation",
                    },
                  ],
                },
              ],
            },
            {
              text: "Schema",
              collapsed: true,
              items: getSchemaSidebar()[0].items,
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

    sidebar: {
      "/schema/": getSchemaSidebar(),
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
