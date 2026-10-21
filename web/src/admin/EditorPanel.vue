<script setup lang="ts">
import { computed } from 'vue'
import type { Content, Project } from '../types'
import FieldText from './FieldText.vue'
import ProjectEditor from './ProjectEditor.vue'

const props = defineProps<{ content: Content }>()
const emit = defineEmits<{ (e: 'change'): void }>()

const c = computed(() => props.content)

function blankProject(): Project {
  const base = 'project'
  const taken = new Set(c.value.projects.map((p) => p.id))
  let id = base
  let n = 2
  while (taken.has(id)) {
    id = `${base}-${n++}`
  }
  return {
    id,
    name: 'Новый проект',
    kind: '',
    tagline: '',
    summary: '',
    bullets: [],
    tags: [],
    demo: '',
    demoText: '',
    year: String(new Date().getFullYear()),
    status: 'черновик',
    featured: false,
  }
}

function addProject() {
  c.value.projects.push(blankProject())
  emit('change')
}

function duplicateProject(index: number) {
  const source = c.value.projects[index]
  const copy: Project = JSON.parse(JSON.stringify(source))
  const taken = new Set(c.value.projects.map((p) => p.id))
  let id = `${source.id}-copy`
  let n = 2
  while (taken.has(id)) {
    id = `${source.id}-copy-${n++}`
  }
  copy.id = id
  copy.featured = false
  c.value.projects.splice(index + 1, 0, copy)
  emit('change')
}

function removeProject(index: number) {
  const project = c.value.projects[index]
  if (!window.confirm(`Удалить «${project.name}»? Это применится после сохранения.`)) return
  c.value.projects.splice(index, 1)
  emit('change')
}

function moveProject(index: number, delta: number) {
  const target = index + delta
  const list = c.value.projects
  if (target < 0 || target >= list.length) return
  ;[list[index], list[target]] = [list[target], list[index]]
  emit('change')
}

function addFact() {
  c.value.hero.facts.push({ key: '', value: '' })
}

function removeFact(index: number) {
  c.value.hero.facts.splice(index, 1)
}

function addStat() {
  c.value.stats.push({ value: '', label: '' })
}

function removeStat(index: number) {
  c.value.stats.splice(index, 1)
}

function addNav() {
  c.value.nav.push({ label: '', href: '#' })
}

function removeNav(index: number) {
  c.value.nav.splice(index, 1)
}

function addApproach() {
  c.value.approach.items.push({ title: '', body: '' })
}

function removeApproach(index: number) {
  c.value.approach.items.splice(index, 1)
}

function addContact() {
  c.value.contact.items.push({ label: '', value: '', href: '' })
}

function removeContact(index: number) {
  c.value.contact.items.splice(index, 1)
}
</script>

<template>
  <div class="admin-editor">
    <details class="panel" open>
      <summary class="panel__head">
        <h2 class="panel__title">Страница и меню</h2>
        <span class="panel__hint">заголовок вкладки, описание, пункты меню</span>
      </summary>
      <div class="panel__body">
        <FieldText v-model="c.meta.title" label="Заголовок страницы" :maxlength="200" />
        <FieldText v-model="c.meta.description" label="Описание для поиска" :rows="2" :maxlength="400" />

        <div class="field">
          <span class="field__label">Меню</span>
          <div v-for="(item, i) in c.nav" :key="i" class="row">
            <input v-model="item.label" class="input" type="text" placeholder="Проекты" :maxlength="60" />
            <input v-model="item.href" class="input input--mono" type="text" placeholder="#work" :maxlength="400" />
            <div class="row__tools">
              <button class="icon-btn icon-btn--danger" type="button" title="Удалить" @click="removeNav(i)">
                ✕
              </button>
            </div>
          </div>
          <button class="btn btn--ghost btn--sm" type="button" @click="addNav">+ пункт меню</button>
        </div>
      </div>
    </details>

    <details class="panel" open>
      <summary class="panel__head">
        <h2 class="panel__title">Первый экран</h2>
        <span class="panel__hint">заголовок, описание, кнопки, карточка стека</span>
      </summary>
      <div class="panel__body">
        <div class="admin-grid">
          <FieldText v-model="c.hero.kicker" label="Надзаголовок" :maxlength="80" />
          <FieldText v-model="c.hero.title1" label="Заголовок, строка 1" :maxlength="120" />
          <FieldText v-model="c.hero.title2" label="Заголовок, строка 2 (акцент)" :maxlength="120" />
        </div>
        <FieldText v-model="c.hero.lead" label="Лид" :rows="3" :maxlength="800" />

        <div class="admin-grid">
          <FieldText v-model="c.hero.primary.label" label="Кнопка 1 — текст" :maxlength="80" />
          <FieldText v-model="c.hero.primary.href" label="Кнопка 1 — ссылка" :maxlength="400" />
          <FieldText v-model="c.hero.secondary.label" label="Кнопка 2 — текст" :maxlength="80" />
          <FieldText v-model="c.hero.secondary.href" label="Кнопка 2 — ссылка" :maxlength="400" />
        </div>
        <label class="check">
          <input v-model="c.hero.secondary.external" type="checkbox" />
          <span>Кнопка 2 открывается в новой вкладке</span>
        </label>

        <FieldText v-model="c.hero.note" label="Примечание в карточке" :rows="2" :maxlength="300" />

        <div class="field">
          <span class="field__label">Стек на первом экране</span>
          <div v-for="(fact, i) in c.hero.facts" :key="i" class="row row--pair">
            <input v-model="fact.key" class="input" type="text" placeholder="фронтенд" :maxlength="60" />
            <input v-model="fact.value" class="input" type="text" placeholder="Vue 3" :maxlength="200" />
            <div class="row__tools">
              <button class="icon-btn icon-btn--danger" type="button" title="Удалить" @click="removeFact(i)">
                ✕
              </button>
            </div>
          </div>
          <button class="btn btn--ghost btn--sm" type="button" @click="addFact">+ строка</button>
        </div>
      </div>
    </details>

    <details class="panel">
      <summary class="panel__head">
        <h2 class="panel__title">Цифры</h2>
        <span class="panel__hint">полоса под первым экраном</span>
      </summary>
      <div class="panel__body">
        <div v-for="(stat, i) in c.stats" :key="i" class="row row--pair">
          <input v-model="stat.value" class="input" type="text" placeholder="3" :maxlength="40" />
          <input v-model="stat.label" class="input" type="text" placeholder="проекта в проде" :maxlength="120" />
          <div class="row__tools">
            <button class="icon-btn icon-btn--danger" type="button" title="Удалить" @click="removeStat(i)">
              ✕
            </button>
          </div>
        </div>
        <button class="btn btn--ghost btn--sm" type="button" @click="addStat">+ цифра</button>
      </div>
    </details>

    <details class="panel" open>
      <summary class="panel__head">
        <h2 class="panel__title">Проекты</h2>
        <span class="panel__hint">{{ c.projects.length }} шт. · порядок как на сайте</span>
      </summary>
      <div class="panel__body">
        <div class="admin-grid">
          <FieldText v-model="c.work.kicker" label="Надзаголовок секции" :maxlength="80" />
          <FieldText v-model="c.work.title" label="Заголовок секции" :maxlength="160" />
        </div>
        <FieldText v-model="c.work.lead" label="Лид секции" :rows="2" :maxlength="800" />

        <ProjectEditor
          v-for="(project, i) in c.projects"
          :key="project.id + '-' + i"
          :project="project"
          :index="i"
          :total="c.projects.length"
          @move="moveProject(i, $event)"
          @duplicate="duplicateProject(i)"
          @remove="removeProject(i)"
        />

        <button class="btn btn--primary" type="button" @click="addProject">+ Добавить проект</button>
      </div>
    </details>

    <details class="panel">
      <summary class="panel__head">
        <h2 class="panel__title">Подход</h2>
        <span class="panel__hint">{{ c.approach.items.length }} принципов</span>
      </summary>
      <div class="panel__body">
        <div class="admin-grid">
          <FieldText v-model="c.approach.kicker" label="Надзаголовок" :maxlength="80" />
          <FieldText v-model="c.approach.title" label="Заголовок" :maxlength="160" />
        </div>
        <FieldText v-model="c.approach.lead" label="Лид" :rows="2" :maxlength="800" />

        <div v-for="(item, i) in c.approach.items" :key="i" class="subcard">
          <div class="subcard__head">
            <span class="subcard__index">{{ String(i + 1).padStart(2, '0') }}</span>
            <button class="icon-btn icon-btn--danger" type="button" title="Удалить" @click="removeApproach(i)">
              ✕
            </button>
          </div>
          <FieldText v-model="item.title" label="Заголовок" :maxlength="200" />
          <FieldText v-model="item.body" label="Текст" :rows="3" :maxlength="1200" />
        </div>

        <button class="btn btn--ghost btn--sm" type="button" @click="addApproach">+ принцип</button>
      </div>
    </details>

    <details class="panel">
      <summary class="panel__head">
        <h2 class="panel__title">Контакты</h2>
        <span class="panel__hint">{{ c.contact.items.length }} пункта</span>
      </summary>
      <div class="panel__body">
        <div class="admin-grid">
          <FieldText v-model="c.contact.kicker" label="Надзаголовок" :maxlength="80" />
          <FieldText v-model="c.contact.title" label="Заголовок" :maxlength="160" />
        </div>
        <FieldText v-model="c.contact.lead" label="Лид" :rows="2" :maxlength="800" />
        <FieldText v-model="c.contact.note" label="Примечание" :rows="2" :maxlength="300" />

        <div v-for="(item, i) in c.contact.items" :key="i" class="subcard">
          <div class="subcard__head">
            <span class="subcard__index">{{ item.label || 'контакт' }}</span>
            <button class="icon-btn icon-btn--danger" type="button" title="Удалить" @click="removeContact(i)">
              ✕
            </button>
          </div>
          <div class="admin-grid">
            <FieldText v-model="item.label" label="Подпись" :maxlength="60" />
            <FieldText v-model="item.value" label="Значение" :maxlength="200" />
            <FieldText v-model="item.href" label="Ссылка" hint="пусто — не кликабельно" :maxlength="400" />
            <FieldText v-model="item.copied" label="Текст после копирования" :maxlength="60" placeholder="Скопировано" />
          </div>
          <label class="check">
            <input v-model="item.copy" type="checkbox" />
            <span>Кнопка «копировать» вместо ссылки</span>
          </label>
        </div>

        <button class="btn btn--ghost btn--sm" type="button" @click="addContact">+ контакт</button>
      </div>
    </details>

    <details class="panel">
      <summary class="panel__head">
        <h2 class="panel__title">Подвал</h2>
      </summary>
      <div class="panel__body">
        <FieldText v-model="c.footer.note" label="Примечание" :rows="2" :maxlength="300" />
        <FieldText v-model="c.footer.built" label="Пометка о стеке" :maxlength="80" />
      </div>
    </details>
  </div>
</template>
