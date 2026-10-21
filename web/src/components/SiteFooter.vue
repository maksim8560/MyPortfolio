<script setup lang="ts">
import type { Health } from '../types'

defineProps<{
  footer: { note: string; built: string }
  health: Health | null
}>()
</script>

<template>
  <footer class="footer">
    <div class="shell footer__inner">
      <p class="footer__note">{{ footer.note }}</p>

      <p class="footer__status">
        <span class="dot" :class="health ? 'dot--ok' : 'dot--idle'" aria-hidden="true"></span>
        <span v-if="health">
          API {{ health.status }} · {{ health.runtime }} · {{ footer.built }}
        </span>
        <span v-else>API недоступен</span>
      </p>
    </div>
  </footer>
</template>
