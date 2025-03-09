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
            { text: "What is Scaffold?", link: "/introduction/what-is-scaffold", },
            { text: "Quick Start", link: "/introduction/quick-start" },
            { text: "Terminology", link: "/introduction/terminology" },
          ],
        },
        {
          text: "Project Scaffolds",
          items:[
            { text: "Usage", link:"/projects/using-projects"},
            { text: "Creating", link:"/projects/creating-projects"},
            { text: "Available Templates", link:"/projects/featured-scaffolds"}
          ]
        },
        {
          text: "Template Scaffolds",
          items:[
            { text: "Usage", link:"/templates/usage"},
          ]
        },
        {
          text: "Configuration",
          items:[
            { text: "Scaffold Config", link: "/configuration/scaffold-file"},
            { text: "User Config", link: "/configuration/scaffold-rc"}
          ]
        },
        {
          text: "Template System",
          items:[
            { text: "Template Engine", link:"/template-system/template-engine"},
            { text: "Partials", link:"/template-system/partials"},
          ]
        },
        {
          text: "Advanced",
          items:[
            { text: "Editor Support", link: "/advanced/editor-support"},
            { text: "Hooks", link: "/advanced/hooks"},
            { text: "Scaffold Resolution", link: "/advanced/scaffold-resolution"},
            { text: "Testing Scaffolds", link: "/advanced/testing-scaffolds"},
          ]
        },
      ],

      socialLinks: [
        { icon: "github", link: "https://github.com/hay-kot/scaffold" },
      ],
    },
  }),
);
