<script setup lang="ts">
import type { Project } from '../types'
import FieldText from './FieldText.vue'
import StringList from './StringList.vue'

const props = defineProps<{ project: Project; index: number; total: number }>()

const emit = defineEmits<{
  (e: 'move', delta: number): void
  (e: 'duplicate'): void
  (e: 'remove'): void
}>()

const p = props.project
</script>

<template>
  <details class="admin-card" :open="index === 0">
    <summary class="admin-card__summary">
      <span class="admin-card__title">
        <span class="admin-card__index">{{ index + 1 }}</span>
        {{ p.name || 'Без названия' }}
      </span>
      <span class="admin-card__meta">
        <code>{{ p.id || '—' }}</code>
        <span v-if="p.featured" class="badge">широкая</span>
        <span v-if="!p.demo" class="admin-card__warn">без ссылки</span>
      </span>
    </summary>

    <div class="admin-card__body">
      <div class="admin-grid">
        <FieldText v-model="p.id" label="ID" hint="латиница, цифры, дефис" :maxlength="41" />
        <FieldText v-model="p.name" label="Название" :maxlength="120" />
        <FieldText v-model="p.kind" label="Тип" hint="Мессенджер, Игра…" :maxlength="80" />
        <FieldText v-model="p.year" label="Годы" hint="2025—2026" :maxlength="40" />
        <FieldText v-model="p.status" label="Статус" hint="в проде, черновик…" :maxlength="40" />
        <label class="check">
          <input v-model="p.featured" type="checkbox" />
          <span>Широкая карточка (главный проект)</span>
        </label>
      </div>

      <FieldText v-model="p.tagline" label="Короткое описание" :rows="2" :maxlength="400" />
      <FieldText v-model="p.summary" label="Описание" :rows="4" :maxlength="4000" />

      <div class="admin-grid">
        <FieldText
          v-model="p.demo"
          label="Ссылка на проект"
          hint="https://…"
          :maxlength="400"
          placeholder="https://"
        />
        <FieldText v-model="p.demoText" label="Текст кнопки" :maxlength="80" placeholder="Открыть проект" />
      </div>

      <StringList v-model="p.tags" label="Теги" placeholder="Vue" add-label="тег" />
      <StringList
        v-model="p.bullets"
        label="Ключевые пункты"
        placeholder="Что было сложным и как решено"
        add-label="пункт"
      />

      <div class="admin-card__actions">
        <button class="btn btn--ghost btn--sm" type="button" :disabled="index === 0" @click="emit('move', -1)">
          ↑ Выше
        </button>
        <button
          class="btn btn--ghost btn--sm"
          type="button"
          :disabled="index === total - 1"
          @click="emit('move', 1)"
        >
          ↓ Ниже
        </button>
        <button class="btn btn--ghost btn--sm" type="button" @click="emit('duplicate')">⧉ Дублировать</button>
        <button class="btn btn--danger btn--sm" type="button" @click="emit('remove')">Удалить проект</button>
      </div>
    </div>
  </details>
</template>
