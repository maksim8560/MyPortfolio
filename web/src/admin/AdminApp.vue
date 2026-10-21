<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import LoginScreen from './LoginScreen.vue'
import EditorPanel from './EditorPanel.vue'
import { adminApi, AdminError } from './adminApi'
import type { AdminDocument } from './adminApi'
import './admin.css'

const LANGS = [
  { code: 'ru', label: 'Русский' },
  { code: 'en', label: 'English' },
] as const

const checking = ref(true)
const authed = ref(false)
const doc = ref<AdminDocument | null>(null)
const snapshot = ref('')
const rev = ref(0)
const savedAt = ref('')
const lang = ref<'ru' | 'en'>('ru')

const loginError = ref('')
const loginBusy = ref(false)
const busy = ref(false)
const saveError = ref('')
const toast = ref('')
let toastTimer = 0

const dirty = computed(() => doc.value !== null && JSON.stringify(doc.value) !== snapshot.value)

const savedLabel = computed(() => {
  if (!savedAt.value) return ''
  const d = new Date(savedAt.value)
  if (Number.isNaN(d.getTime()) || d.getFullYear() < 2000) return ''
  return d.toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' })
})

async function checkSession() {
  try {
    const { authenticated } = await adminApi.session()
    authed.value = authenticated
    if (authenticated) await load()
  } catch {
    authed.value = false
  } finally {
    checking.value = false
  }
}

async function load() {
  const data = await adminApi.load()
  doc.value = data.content
  snapshot.value = JSON.stringify(data.content)
  rev.value = data.rev
  savedAt.value = data.savedAt
}

async function login(password: string) {
  loginBusy.value = true
  loginError.value = ''
  try {
    await adminApi.login(password)
    authed.value = true
    await load()
  } catch (err) {
    loginError.value = err instanceof Error ? err.message : 'Не удалось войти'
  } finally {
    loginBusy.value = false
  }
}

async function logout() {
  if (dirty.value && !window.confirm('Есть несохранённые правки. Выйти без сохранения?')) return
  try {
    await adminApi.logout()
  } finally {
    authed.value = false
    doc.value = null
  }
}

async function save() {
  if (!doc.value || busy.value) return
  busy.value = true
  saveError.value = ''
  try {
    const result = await adminApi.save(doc.value)
    snapshot.value = JSON.stringify(doc.value)
    rev.value = result.rev
    savedAt.value = result.savedAt
    showToast('Сохранено — сайт уже отдаёт новый контент')
  } catch (err) {
    if (err instanceof AdminError && err.status === 401) {
      authed.value = false
      loginError.value = 'Сессия истекла, войди заново'
      return
    }
    saveError.value = err instanceof Error ? err.message : 'Не удалось сохранить'
  } finally {
    busy.value = false
  }
}

function showToast(text: string) {
  toast.value = text
  window.clearTimeout(toastTimer)
  toastTimer = window.setTimeout(() => {
    toast.value = ''
  }, 3200)
}

function onKeydown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 's') {
    event.preventDefault()
    void save()
  }
}

function onBeforeUnload(event: BeforeUnloadEvent) {
  if (!dirty.value) return
  event.preventDefault()
  event.returnValue = ''
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
  window.addEventListener('beforeunload', onBeforeUnload)
  void checkSession()
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  window.removeEventListener('beforeunload', onBeforeUnload)
  window.clearTimeout(toastTimer)
})
</script>

<template>
  <div class="admin">
    <p v-if="checking" class="admin-checking">Проверяем сессию…</p>

    <LoginScreen
      v-else-if="!authed"
      :error="loginError"
      :busy="loginBusy"
      @submit="login"
    />

    <template v-else-if="doc">
      <header class="admin-bar">
        <div class="admin-bar__inner">
          <a class="admin-bar__brand" href="/">← на сайт</a>
          <div class="admin-bar__titles">
            <span class="admin-bar__title">Админка портфолио</span>
            <span class="admin-bar__meta">
              <span :class="dirty ? 'pill pill--warn' : 'pill pill--ok'">
                {{ dirty ? 'не сохранено' : 'сохранено' }}
              </span>
              <span v-if="savedLabel">изменение {{ rev }} · {{ savedLabel }}</span>
            </span>
          </div>
          <div class="admin-bar__actions">
            <button class="btn btn--ghost btn--sm" type="button" @click="logout">Выйти</button>
            <button class="btn btn--primary btn--sm" type="button" :disabled="!dirty || busy" @click="save">
              {{ busy ? 'Сохраняем…' : 'Сохранить' }}
            </button>
          </div>
        </div>
      </header>

      <div class="admin-subbar">
        <div class="admin-subbar__inner">
          <div class="tabs" role="tablist" aria-label="Язык контента">
            <button
              v-for="item in LANGS"
              :key="item.code"
              class="tabs__btn"
              :class="{ 'tabs__btn--on': lang === item.code }"
              type="button"
              role="tab"
              :aria-selected="lang === item.code"
              @click="lang = item.code"
            >
              {{ item.label }}
              <span class="tabs__count">{{ doc[item.code].projects.length }}</span>
            </button>
          </div>
          <p class="admin-subbar__hint">
            Языки редактируются отдельно · ⌘/Ctrl + S
          </p>
          <a class="admin-subbar__link" href="/" target="_blank" rel="noopener noreferrer">
            открыть сайт ↗
          </a>
        </div>
      </div>

      <main class="admin-main">
        <p v-if="saveError" class="alert" role="alert">{{ saveError }}</p>
        <EditorPanel :content="doc[lang]" />
      </main>

      <transition name="toast">
        <p v-if="toast" class="toast" role="status">{{ toast }}</p>
      </transition>
    </template>
  </div>
</template>
