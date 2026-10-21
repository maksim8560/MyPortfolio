<script setup lang="ts">
import type { Stat } from '../types'

defineProps<{ stats: Stat[] }>()
</script>

<template>
  <section v-if="stats.length" class="shell stats" aria-label="Коротко о работе">
    <div v-for="stat in stats" :key="stat.label" class="stat">
      <span class="stat__value">{{ stat.value }}</span>
      <span class="stat__label">{{ stat.label }}</span>
    </div>
  </section>
</template>
