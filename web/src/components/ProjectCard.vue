<script setup lang="ts">
import { computed } from 'vue'
import type { Project } from '../types'

const props = defineProps<{ project: Project }>()

const externalLink = computed(() => props.project.demo ?? '')
</script>

<template>
  <article class="project" :class="{ 'project--featured': project.featured }">
    <div class="project__top">
      <p class="project__kind">
        <span v-if="project.featured" class="badge">{{ project.status }}</span>
        <span>{{ project.kind }}</span>
        <span class="project__year">{{ project.year }}</span>
      </p>
      <h3 class="project__name">{{ project.name }}</h3>
      <p class="project__tagline">{{ project.tagline }}</p>
    </div>

    <div class="project__body">
      <p class="project__summary">{{ project.summary }}</p>

      <ul v-if="project.bullets.length" class="project__list">
        <li v-for="point in project.bullets" :key="point">{{ point }}</li>
      </ul>

      <ul class="tags" aria-label="Технологии">
        <li v-for="tag in project.tags" :key="tag" class="tag">{{ tag }}</li>
      </ul>
    </div>

    <footer v-if="project.demo" class="project__foot">
      <a
        class="project__link"
        :href="externalLink"
        target="_blank"
        rel="noopener noreferrer"
      >
        {{ project.demoText || project.name }}
        <span aria-hidden="true">↗</span>
        <span class="sr-only">(откроется в новой вкладке)</span>
      </a>
      <span class="project__url">{{ project.demo.replace(/^https?:\/\//, '') }}</span>
    </footer>
  </article>
</template>
