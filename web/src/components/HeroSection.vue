<script setup lang="ts">
import type { Content } from '../types'

defineProps<{ hero: Content['hero'] }>()
</script>

<template>
  <section class="hero shell" aria-labelledby="hero-title">
    <div class="hero__main">
      <p class="kicker">
        <span class="kicker__rule" aria-hidden="true"></span>
        {{ hero.kicker }}
      </p>

      <h1 id="hero-title" class="hero__title">
        <span class="hero__title-line">{{ hero.title1 }}</span>
        <span class="hero__title-line hero__title-line--accent">{{ hero.title2 }}</span>
      </h1>

      <p class="hero__lead">{{ hero.lead }}</p>

      <div class="hero__actions">
        <a class="btn btn--primary" :href="hero.primary.href">{{ hero.primary.label }}</a>
        <a
          class="btn btn--ghost"
          :href="hero.secondary.href"
          :target="hero.secondary.external ? '_blank' : undefined"
          :rel="hero.secondary.external ? 'noopener noreferrer' : undefined"
        >
          {{ hero.secondary.label }}
          <span class="btn__arrow" aria-hidden="true">↗</span>
        </a>
      </div>
    </div>

    <aside class="hero__aside" aria-label="Технический стек">
      <div class="fact-card">
        <div class="fact-card__head">
          <span class="dot" aria-hidden="true"></span>
          <span>сейчас в работе</span>
        </div>
        <dl class="facts">
          <div v-for="fact in hero.facts" :key="fact.key" class="facts__row">
            <dt class="facts__key">{{ fact.key }}</dt>
            <dd class="facts__value">{{ fact.value }}</dd>
          </div>
        </dl>
        <p class="fact-card__foot">{{ hero.note }}</p>
      </div>
    </aside>
  </section>
</template>
