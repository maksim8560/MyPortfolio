<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import SiteHeader from './components/SiteHeader.vue'
import HeroSection from './components/HeroSection.vue'
import StatStrip from './components/StatStrip.vue'
import WorkSection from './components/WorkSection.vue'
import ApproachSection from './components/ApproachSection.vue'
import ContactSection from './components/ContactSection.vue'
import SiteFooter from './components/SiteFooter.vue'
import { fetchContent, fetchHealth } from './api'
import type { Content, Health } from './types'

const LANGS = ['ru', 'en'] as const
const LANG_KEY = 'portfolio.lang'

const lang = ref<string>(readLang())
const content = ref<Content | null>(null)
const health = ref<Health | null>(null)
const error = ref<string>('')
const loading = ref(true)

let controller: AbortController | null = null

function readLang(): string {
  try {
    const saved = localStorage.getItem(LANG_KEY)
    if (saved && (LANGS as readonly string[]).includes(saved)) return saved
  } catch {
    // приватный режим — просто идём на язык браузера
  }
  const nav = navigator.language.slice(0, 2).toLowerCase()
  return (LANGS as readonly string[]).includes(nav) ? nav : 'ru'
}

function setLang(next: string) {
  if (next === lang.value) return
  lang.value = next
  try {
    localStorage.setItem(LANG_KEY, next)
  } catch {
    // не критично: язык сбросится при следующем визите
  }
}

async function load(target: string) {
  controller?.abort()
  const current = new AbortController()
  controller = current
  loading.value = true
  error.value = ''

  try {
    content.value = await fetchContent(target, current.signal)
    applyMeta(content.value)
  } catch (err) {
    if (current.signal.aborted) return
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить контент'
  } finally {
    if (!current.signal.aborted) loading.value = false
  }
}

function applyMeta(data: Content) {
  document.title = data.meta.title
  const desc = document.querySelector('meta[name="description"]')
  if (desc) desc.setAttribute('content', data.meta.description)
  document.documentElement.lang = data.lang
}

watch(lang, (next) => void load(next))

onMounted(async () => {
  await load(lang.value)
  try {
    health.value = await fetchHealth()
  } catch {
    // Бейдж бэкенда просто не покажется — страница от этого не ломается.
  }
})

onBeforeUnmount(() => controller?.abort())

const showSkeleton = computed(() => loading.value && !content.value)
</script>

<template>
  <a class="skip-link" href="#work">Перейти к проектам</a>

  <div class="page">
    <div class="page__glow" aria-hidden="true"></div>

    <SiteHeader
      v-if="content"
      :nav="content.nav"
      :lang="content.lang"
      :langs="LANGS"
      @update:lang="setLang"
    />

    <main id="main">
      <template v-if="showSkeleton">
        <div class="shell">
          <div class="skeleton" aria-hidden="true">
            <span class="skeleton__line skeleton__line--sm"></span>
            <span class="skeleton__line skeleton__line--xl"></span>
            <span class="skeleton__line skeleton__line--lg"></span>
            <span class="skeleton__line"></span>
          </div>
          <p class="sr-only" role="status">Загружаем содержимое</p>
        </div>
      </template>

      <div v-else-if="error" class="shell">
        <div class="failure" role="alert">
          <h1 class="failure__title">Контент не приехал</h1>
          <p class="failure__text">
            {{ error }}. Это ошибка API, а не «пусто» — страницу рисовать нечем.
          </p>
          <p class="failure__hint">
            <code>go run ./server</code> — если бэкенд не запущен.
          </p>
          <button class="btn btn--primary" type="button" @click="load(lang)">Повторить</button>
        </div>
      </div>

      <template v-else-if="content">
        <HeroSection :hero="content.hero" />
        <StatStrip :stats="content.stats" />
        <WorkSection :section="content.work" :projects="content.projects" />
        <ApproachSection
          :section="content.approach"
          :items="content.approach.items"
        />
        <ContactSection
          :section="content.contact"
          :items="content.contact.items"
          :note="content.contact.note"
        />
      </template>
    </main>

    <SiteFooter v-if="content" :footer="content.footer" :health="health" />
  </div>
</template>
