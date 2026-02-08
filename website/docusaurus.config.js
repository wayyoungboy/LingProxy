// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'LingProxy',
  tagline: 'High-performance AI API Gateway',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://lingproxy.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',  // Update to '/lingproxy/' if deploying to GitHub Pages subdirectory

  // GitHub pages deployment config.
  organizationName: 'lingproxy',
  projectName: 'lingproxy',

  onBrokenLinks: 'warn',  // 改为警告，允许构建通过

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'zh'],
    localeConfigs: {
      en: {
        label: 'English',
        direction: 'ltr',
        htmlLang: 'en-US',
      },
      zh: {
        label: '中文',
        direction: 'ltr',
        htmlLang: 'zh-CN',
      },
    },
  },

  // 确保首页路由正确
  plugins: [],

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          path: '../docs',  // 从项目根目录的 docs/ 读取文档，Docusaurus 会自动根据语言读取 en/ 或 zh/ 文件夹
          routeBasePath: 'docs',  // URL 路径保持 /docs
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/lingproxy/lingproxy/tree/main/',
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
          sidebarCollapsible: true,
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themes: ['@docusaurus/theme-mermaid'],

  markdown: {
    mermaid: true,
    hooks: {
      onBrokenMarkdownLinks: 'ignore',  // 忽略外部链接警告（如 docker/README.md）
    },
  },

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/lingproxy-social-card.jpg',
      navbar: {
        title: 'LingProxy',
        logo: {
          alt: 'LingProxy Logo',
          src: 'img/logo.svg',
          srcDark: 'img/logo.svg',
        },
        style: 'primary',
        hideOnScroll: true,
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'docs',
            position: 'left',
            label: 'Docs',
          },
          {
            type: 'localeDropdown',
            position: 'right',
          },
          {
            href: 'https://github.com/lingproxy/lingproxy',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'light',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Introduction',
                to: '/docs/introduction',
              },
              {
                label: 'Quick Start',
                to: '/docs/quick-start',
              },
              {
                label: 'API Reference',
                to: '/docs/api-reference',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/lingproxy/lingproxy',
              },
              {
                label: 'Issues',
                href: 'https://github.com/lingproxy/lingproxy/issues',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/lingproxy/lingproxy',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} LingProxy. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ['bash', 'go', 'javascript', 'json', 'yaml'],
      },
      colorMode: {
        defaultMode: 'light',
        disableSwitch: false,
        respectPrefersColorScheme: false,
      },
    }),
};

module.exports = config;
