---
scaffolds:
  - name: hay-kot/scaffold-go-cli
    href: https://github.com/hay-kot/scaffold-go-cli
    description: |
      Scaffold for stubbing out an opinionated CLI project using <a href="https://github.com/urfave/cli">urfave/cli</a>
    features:
      - CI/CD w/ Github Actions or Drone.io
      - GoReleaser for releases
      - GolangCI-Lint for linting
      - Build/Version/Commit injection on build
---

# Featured Scaffolds

<script setup>
import Featured from '../components/Featured.vue'
import { useData } from 'vitepress'
const { page } = useData()
</script>

<Featured :features="page.frontmatter.scaffolds"/>
