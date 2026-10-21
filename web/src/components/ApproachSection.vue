<script setup lang="ts">
import type { ApproachItem, Section } from '../types'

defineProps<{ section: Section; items: ApproachItem[] }>()
</script>

<template>
  <section id="approach" class="shell section" aria-labelledby="approach-title">
    <header class="section__head">
      <p class="kicker">
        <span class="kicker__rule" aria-hidden="true"></span>
        {{ section.kicker }}
      </p>
      <h2 id="approach-title" class="section__title">{{ section.title }}</h2>
      <p class="section__lead">{{ section.lead }}</p>
    </header>

    <div class="approach">
      <article v-for="(item, i) in items" :key="item.title" class="approach__item">
        <span class="approach__index" aria-hidden="true">{{ String(i + 1).padStart(2, '0') }}</span>
        <div>
          <h3 class="approach__title">{{ item.title }}</h3>
          <p class="approach__body">{{ item.body }}</p>
        </div>
      </article>
    </div>
  </section>
</template>
