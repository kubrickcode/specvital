import { defineConfig } from "vitepress";

export default defineConfig({
  title: "Specvital",
  description: "Open-source test coverage insights for code review",

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
                { text: "Competitive Differentiation", link: "/en/adr/02-competitive-differentiation" },
                { text: "API Worker Separation", link: "/en/adr/04-api-worker-service-separation" },
                { text: "Queue-Based Processing", link: "/en/adr/05-queue-based-async-processing" },
                { text: "Repository Strategy", link: "/en/adr/06-repository-strategy" },
                { text: "PaaS-First Infrastructure", link: "/en/adr/07-paas-first-infrastructure" },
                { text: "Shared Infrastructure", link: "/en/adr/08-shared-infrastructure" },
                {
                  text: "Core",
                  collapsed: false,
                  items: [
                    { text: "Overview", link: "/en/adr/core/" },
                    { text: "Core Library Separation", link: "/en/adr/core/01-core-library-separation" },
                  ],
                },
                {
                  text: "Collector",
                  collapsed: false,
                  items: [
                    { text: "Overview", link: "/en/adr/collector/" },
                    { text: "Scheduled Re-collection", link: "/en/adr/collector/01-scheduled-recollection" },
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
                { text: "API/Worker 분리", link: "/ko/adr/04-api-worker-service-separation" },
                { text: "큐 기반 비동기 처리", link: "/ko/adr/05-queue-based-async-processing" },
                { text: "리포지토리 전략", link: "/ko/adr/06-repository-strategy" },
                { text: "PaaS 우선 인프라", link: "/ko/adr/07-paas-first-infrastructure" },
                { text: "공유 인프라", link: "/ko/adr/08-shared-infrastructure" },
                {
                  text: "Core",
                  collapsed: false,
                  items: [
                    { text: "개요", link: "/ko/adr/core/" },
                    { text: "코어 라이브러리 분리", link: "/ko/adr/core/01-core-library-separation" },
                  ],
                },
                {
                  text: "Collector",
                  collapsed: false,
                  items: [
                    { text: "개요", link: "/ko/adr/collector/" },
                    { text: "스케줄 기반 재수집", link: "/ko/adr/collector/01-scheduled-recollection" },
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
    logo: "/logo.svg",

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
