<script setup lang="ts">
import { computed } from 'vue'

const model = defineModel<string | undefined>({ required: true })

const props = withDefaults(
  defineProps<{
    label: string
    hint?: string
    rows?: number
    placeholder?: string
    maxlength?: number
  }>(),
  { rows: 0, hint: '', placeholder: '', maxlength: 1000 },
)

const multiline = computed(() => props.rows > 0)
const id = `f-${Math.random().toString(36).slice(2, 9)}`
</script>

<template>
  <div class="field">
    <label class="field__label" :for="id">
      {{ label }}
      <span v-if="hint" class="field__hint">{{ hint }}</span>
    </label>

    <textarea
      v-if="multiline"
      :id="id"
      v-model="model"
      class="input input--area"
      :rows="rows"
      :placeholder="placeholder"
      :maxlength="maxlength"
    ></textarea>

    <input
      v-else
      :id="id"
      v-model="model"
      class="input"
      type="text"
      :placeholder="placeholder"
      :maxlength="maxlength"
    />
  </div>
</template>
