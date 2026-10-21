<script setup lang="ts">
const model = defineModel<string[]>({ required: true })
const { placeholder = '', addLabel = 'Добавить' } = defineProps<{
  label: string
  placeholder?: string
  addLabel?: string
}>()

function add() {
  model.value = [...model.value, '']
}

function remove(index: number) {
  model.value = model.value.filter((_, i) => i !== index)
}

function move(index: number, delta: number) {
  const next = [...model.value]
  const target = index + delta
  if (target < 0 || target >= next.length) return
  ;[next[index], next[target]] = [next[target], next[index]]
  model.value = next
}

function update(index: number, value: string) {
  const next = [...model.value]
  next[index] = value
  model.value = next
}
</script>

<template>
  <div class="field field--list">
    <span class="field__label">{{ label }}</span>

    <p v-if="!model.length" class="field__empty">Пусто — блок на сайте не отрисуется.</p>

    <div v-for="(item, index) in model" :key="index" class="row">
      <input
        class="input"
        type="text"
        :value="item"
        :placeholder="placeholder"
        @input="update(index, ($event.target as HTMLInputElement).value)"
      />
      <div class="row__tools">
        <button class="icon-btn" type="button" title="Выше" :disabled="index === 0" @click="move(index, -1)">
          ↑
        </button>
        <button
          class="icon-btn"
          type="button"
          title="Ниже"
          :disabled="index === model.length - 1"
          @click="move(index, 1)"
        >
          ↓
        </button>
        <button class="icon-btn icon-btn--danger" type="button" title="Удалить" @click="remove(index)">
          ✕
        </button>
      </div>
    </div>

    <button class="btn btn--ghost btn--sm" type="button" @click="add">+ {{ addLabel }}</button>
  </div>
</template>
