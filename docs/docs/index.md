---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Scaffold"
  text: "Project Template Scaffolding Tool"
  tagline: Generate Boilerplate Code with Ease
  actions:
    - theme: brand
      text: What is Scaffold?
      link: /introduction/what-is-scaffold
    - theme: alt
      text: Quick Start
      link: /introduction/quick-start
    - theme: alt
      text: Creating Templates
      link: /configuration/scaffold-file
  image: /imgs/scaffold-gopher.webp

features:
  - title: Feature Flag Support
    link: /configuration/scaffold-file#features
    details: Guard large sections of your templates with feature flags to allow for easy feature toggling.

  - title: Powerful Interactive Forms
    link: /configuration/scaffold-file#prompts
    details: Grouped Inputs, Conditional Fields and more all powered by Charm.sh.

  - title: In Project Scaffolding
    link: /configuration/scaffold-file#template-scaffolds
    details: Scaffolds can run within your project creating files across the project, and even injecting code into existing files.

  - title: Supports Testing Templates
    link: /advanced/testing-scaffolds
    details: Scaffold provides several utilities and features to help you test your scaffolds and make sure
---
