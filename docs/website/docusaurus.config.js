// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("./src/css/prism-theme");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

const companyWebsiteUrl =
  process.env.COMPANY_WEBSITE_URL ?? "https://permit.io";

const permitUrl = process.env.PERMIT_DOCS_URL ?? "https://docs.permit.io/";

const algoliaAppId = process.env.APPLICATION_ID || "";
const algoliaApiKey = process.env.API_KEY || "";

let searchConfig = {};
if (algoliaAppId && algoliaApiKey) {
  searchConfig = {
    algolia: {
      appId: process.env.APPLICATION_ID,
      apiKey: process.env.API_KEY,
      indexName: "docs",
    },
  };
}

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Permit.io",
  tagline: "Fullstack Authorization",
  url: "https://docs.permit.io",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "permitio", // Usually your GitHub org/user name.
  projectName: "docs", // Usually your repo name.

  presets: [
    [
      "@docusaurus/preset-classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          routeBasePath: "/",
        },
        blog: false, // disables docusaurus blog
        theme: {
          customCss: require.resolve("./src/css/custom.scss"),
        },
      }),
    ],
  ],

  plugins: ["docusaurus-plugin-sass"],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        logo: {
          alt: "OPToggles",
          src: "img/optoggles-logo.svg",
          srcDark: "img/optoggles-logo-dark.svg",
          width: 200,
        },
        items: [
          {
            type: "doc",
            docId: "home",
            position: "left",
            label: "Docs",
          },
          {
            to: `${permitUrl}`,
            label: "Permit docs",
            position: "left",
            target: "_self",
          },
          {
            type: "doc",
            docId: "changelog",
            position: "left",
            label: "Changelog",
          },
          {
            href: "https://github.com/permitio/OPToggles",
            position: "right",
            className: "header-github-link",
            "aria-label": "Github",
          },
          {
            href: "https://bit.ly/opal-slack",
            position: "right",
            className: "header-slack-link",
            "aria-label": "Slack community",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [],
        copyright: `Copyright © ${new Date().getFullYear()} Permit Inc.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      //search configuration
      ...searchConfig,
      announcementBar: {
        id: "support_us",
        content:
          '⭐️ If you like OPToggles, give it a star on <a href="https://www.github.com/permitio/optoggles" target="_blank" rel="noopener noreferrer">GitHub</a> and follow us on <a href="https://www.twitter.com/permit_io" target="_blank" rel="noopener noreferrer">Twitter</a>',
        backgroundColor: "#6851ff",
        textColor: "#FFFFFF",
        isCloseable: true,
      },
    }),
};

module.exports = config;
