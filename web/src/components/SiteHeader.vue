<script setup lang="ts">
import type { Link } from '../types'

defineProps<{
  nav: Link[]
  lang: string
  langs: readonly string[]
}>()

const emit = defineEmits<{ (e: 'update:lang', value: string): void }>()

const LABELS: Record<string, string> = { ru: 'RU', en: 'EN' }
</script>

<template>
  <header class="header">
    <div class="shell header__inner">
      <a class="brand" href="#main">
        <span class="brand__mark" aria-hidden="true">M</span>
        <span class="brand__text">
          <span class="brand__name">MaxDiWay</span>
          <span class="brand__role">Максим · разработчик</span>
        </span>
      </a>

      <nav class="nav" aria-label="Основная навигация">
        <a v-for="item in nav" :key="item.href" class="nav__link" :href="item.href">
          {{ item.label }}
        </a>
      </nav>

      <div class="lang" role="group" aria-label="Язык интерфейса">
        <button
          v-for="code in langs"
          :key="code"
          class="lang__btn"
          :class="{ 'lang__btn--on': code === lang }"
          type="button"
          :aria-pressed="code === lang"
          :lang="code"
          @click="emit('update:lang', code)"
        >
          {{ LABELS[code] ?? code.toUpperCase() }}
        </button>
      </div>
    </div>
  </header>
</template>
