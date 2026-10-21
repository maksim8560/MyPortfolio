<script setup lang="ts">
import { ref } from 'vue'
import type { ContactItem, Section } from '../types'

defineProps<{ section: Section; items: ContactItem[]; note: string }>()

const copiedLabel = ref('')
let resetTimer: number | undefined

async function copy(value: string) {
  const ok = await writeClipboard(value)
  copiedLabel.value = ok ? value : ''
  window.clearTimeout(resetTimer)
  resetTimer = window.setTimeout(() => {
    copiedLabel.value = ''
  }, 2000)
}

// navigator.clipboard есть только в защищённом контексте (https или localhost),
// поэтому нужен запасной путь — иначе кнопка молча не работает.
async function writeClipboard(value: string): Promise<boolean> {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(value)
      return true
    }
  } catch {
    // падаем в фолбэк ниже
  }
  try {
    const area = document.createElement('textarea')
    area.value = value
    area.setAttribute('readonly', '')
    area.style.position = 'fixed'
    area.style.top = '-1000px'
    document.body.appendChild(area)
    area.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(area)
    return ok
  } catch {
    return false
  }
}
</script>

<template>
  <section id="contact" class="shell section" aria-labelledby="contact-title">
    <div class="contact">
      <header class="contact__head">
        <p class="kicker">
          <span class="kicker__rule" aria-hidden="true"></span>
          {{ section.kicker }}
        </p>
        <h2 id="contact-title" class="section__title">{{ section.title }}</h2>
        <p class="section__lead">{{ section.lead }}</p>
      </header>

      <ul class="contact__list">
        <li v-for="item in items" :key="item.label" class="contact__row">
          <span class="contact__label">{{ item.label }}</span>

          <a
            v-if="item.href"
            class="contact__value contact__value--link"
            :href="item.href"
            target="_blank"
            rel="noopener noreferrer"
          >
            {{ item.value }}
            <span aria-hidden="true">↗</span>
          </a>

          <button
            v-else-if="item.copy"
            class="contact__value contact__value--button"
            type="button"
            @click="copy(item.value)"
          >
            {{ item.value }}
            <span class="contact__copy">
              {{ copiedLabel === item.value ? (item.copied || 'OK') : 'копировать' }}
            </span>
          </button>

          <span v-else class="contact__value">{{ item.value }}</span>
        </li>
      </ul>

      <p class="contact__note">{{ note }}</p>
    </div>
  </section>
</template>
