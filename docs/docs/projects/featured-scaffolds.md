---
scaffolds:
  - name: hay-kot/scaffold-go-cli
    href: https://github.com/hay-kot/scaffold-go-cli
    description: |
      template for stubbing out an opinionated CLI project using <a href="https://github.com/urfave/cli">urfave/cli</a>
    features:
      - CI/CD w/ Github Actions or Drone.io
      - GoReleaser for releases
      - GolangCI-Lint for linting
      - Build/Version/Commit injection on build
  - name: hay-kot/scaffold-go-pkg
    href: https://github.com/hay-kot/scaffold-go-pkg
    description: |
     template for creating a new go package
    features:
      - Taskfile for common tasks
      - Renovate for dependency updates
      - Github Actions for CI/CD (Testing and Linting)
      - Release Drafter for automated release notes
---

# Featured Scaffolds

<script setup>
import Featured from '../components/Featured.vue'
import { useData } from 'vitepress'
const { page } = useData()
</script>

<Featured :features="page.frontmatter.scaffolds"/>
