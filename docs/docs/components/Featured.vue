<template>
  <div class="feature-grid">
    <a
      :href="feature.href"
      target="_blank"
      v-for="feature in features"
      class="vp-raw"
    >
      <div class="feature-card">
        <p class="title">{{ feature.name }}</p>
        <p class="description" v-html="feature.description"></p>
        <ul class="list">
          <li v-for="item in feature.features">
            <span>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 16 16"
              >
                <path
                  fill-rule="evenodd"
                  d="M12.416 3.376a.75.75 0 0 1 .208 1.04l-5 7.5a.75.75 0 0 1-1.154.114l-3-3a.75.75 0 0 1 1.06-1.06l2.353 2.353l4.493-6.74a.75.75 0 0 1 1.04-.207"
                  clip-rule="evenodd"
                />
              </svg>
            </span>
            {{ item }}
          </li>
        </ul>
      </div>
    </a>
  </div>
</template>

<script setup lang="ts">
type Props = {
  features: {
    href: string;
    name: string;
    description: string;
    deatures: string[];
  }[];
};

defineProps<Props>();
</script>

<style lang="css" scoped>
.vp-raw {
  all: unset !important;
}

a {
  cursor: pointer;
}

.feature-grid {
  width: 100%;
  display: grid;
  margin: 2rem 0;
  gap: 1rem;
}

.vp-doc p {
  margin: 0;  
}

.feature-card {
  background-color: #202126;
  padding: 1rem;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.3s;
}

.feature-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.feature-card .title {
  font-size: 1.5rem;
  font-weight: bold;
}

.feature-card .description {
  font-size: 14px;
  margin-bottom: 0.75rem;
  margin-top: 0.25rem; 
  line-height: 1.5;
  color: var(--vp-c-text-2);
}

.feature-card .list {
  margin: 0;
  padding: 0;
  list-style: none;
}

.list li {
  display: flex;
  align-items: start;
  font-size: 14px;
  padding-left: 0.75rem;
  margin: 0.5rem 0;
}

li svg {
  fill: var(--vp-c-brand-1);
  margin-right: 0.25rem;
  margin-bottom: auto;
}
</style>
