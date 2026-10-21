<script setup lang="ts">
import type { Section, Project } from '../types'
import ProjectCard from './ProjectCard.vue'

defineProps<{ section: Section; projects: Project[] }>()
</script>

<template>
  <section id="work" class="shell section" aria-labelledby="work-title">
    <header class="section__head">
      <p class="kicker">
        <span class="kicker__rule" aria-hidden="true"></span>
        {{ section.kicker }}
      </p>
      <h2 id="work-title" class="section__title">{{ section.title }}</h2>
      <p class="section__lead">{{ section.lead }}</p>
    </header>

    <div class="projects">
      <ProjectCard
        v-for="project in projects"
        :key="project.id"
        :project="project"
        :class="{ 'project--wide': project.featured }"
      />
    </div>
  </section>
</template>
