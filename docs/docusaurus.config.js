// @ts-check
// Note: type annotations allow type checking and IDEs autosuggestion

const lightCodeTheme = require("prism-react-renderer/themes/nightOwlLight");
const darkCodeTheme = require("prism-react-renderer/themes/nightOwl");

// #79BBA5
// #47AFA2
// #55AFC3
// #AABEBF
// #6D3573

async function createConfig() {
  const gfm = (await import("remark-gfm")).default;

  /** @type {import('@docusaurus/types').Config} */
  const config = {
    title: "Bubbleprompt",
    tagline: "Prompts for your terminal",
    url: "https://aschey.tech",
    baseUrl: "/bubbleprompt",
    onBrokenLinks: "throw",
    onBrokenMarkdownLinks: "warn",
    favicon: "img/favicon.ico",

    // Even if you don't use internalization, you can use this field to set useful
    // metadata like html lang. For example, if your site is Chinese, you may want
    // to replace "en" with "zh-Hans".
    i18n: {
      defaultLocale: "en",
      locales: ["en"],
    },

    presets: [
      [
        "classic",
        /** @type {import('@docusaurus/preset-classic').Options} */
        ({
          docs: {
            routeBasePath: "/",
            // Remove this to remove the "edit this page" links.
            editUrl: "https://github.com/aschey/bubbleprompt/tree/main/docs/",
            remarkPlugins: [gfm],
          },
          blog: false,
          theme: {
            customCss: require.resolve("./src/css/custom.css"),
          },
        }),
      ],
    ],

    themeConfig:
      /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
      ({
        navbar: {
          title: "Bubbleprompt",
          logo: {
            alt: "Bubbleprompt Logo",
            src: "img/logo.svg",
          },
          items: [
            {
              type: "doc",
              docId: "Intro",
              position: "left",
              label: "Documentation",
            },
            {
              href: "https://github.com/aschey/bubbleprompt",
              label: "GitHub",
              position: "right",
            },
          ],
        },
        footer: {},
        prism: {
          theme: lightCodeTheme,
          darkTheme: darkCodeTheme,
        },
      }),
  };

  return config;
}

module.exports = createConfig;
