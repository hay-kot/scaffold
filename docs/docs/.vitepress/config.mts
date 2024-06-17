import { defineConfig } from "vitepress";
import { withMermaid } from "vitepress-plugin-mermaid";
// https://vitepress.dev/reference/site-config
export default withMermaid(
  defineConfig({
    base: "/scaffold/",
    title: "Scaffold",
    description: "A Project and Template Scaffolding Tool",
    head: [["link", { rel: "icon", href: "/scaffold/favicon.webp" }]],
    themeConfig: {
      search: {
        provider: "local",
        options: {
          detailedView: true,
        },
      },
      // https://vitepress.dev/reference/default-theme-config
      nav: [
        { text: "Home", link: "/" },
        { text: "Docs", link: "/introduction/what-is-scaffold" },
      ],
      outline: "deep",

      sidebar: [
        {
          text: "Introduction",
          items: [
            {
              text: "What is Scaffold?",
              link: "/introduction/what-is-scaffold",
            },
            { text: "Quick Start", link: "/introduction/quick-start" },
            { text: "Terminology", link: "/introduction/terminology" },
          ],
        },
        {
          text: "User Guide",
          items: [
            { text: "User Configuration", link: "/user-guide/scaffold-rc" },
            {
              text: "Scaffold Resolution",
              link: "/user-guide/scaffold-resolution",
            },
            {
              text: "Featured Scaffolds",
              link: "/user-guide/featured-scaffolds",
            },
          ],
        },
        {
          text: "Creating Scaffolds",
          items: [
            { text: "Scaffold File", link: "/templates/scaffold-file" },
            { text: "File Reference", link: "/templates/config-reference" },
            { text: "Template Engine", link: "/templates/template-engine" },
            { text: "Hooks", link: "/templates/hooks" },
            { text: "Testing Scaffolds", link: "/templates/testing-scaffolds" },
          ],
        },
        {
          text: "Advanced",
          items: [{ text: "Editor Support", link: "/advanced/editor-support" }],
        },
      ],

      socialLinks: [
        { icon: "github", link: "https://github.com/hay-kot/scaffold" },
      ],
    },
  }),
);
