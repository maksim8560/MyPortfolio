<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{ (e: 'submit', password: string): void }>()
const props = defineProps<{ error: string; busy: boolean }>()

const password = ref('')
const input = ref<HTMLInputElement | null>(null)

function submit() {
  if (!password.value || props.busy) return
  emit('submit', password.value)
}
</script>

<template>
  <div class="admin-login">
    <form class="admin-login__card" @submit.prevent="submit">
      <p class="kicker">
        <span class="kicker__rule" aria-hidden="true"></span>
        админка
      </p>
      <h1 class="admin-login__title">Вход</h1>
      <p class="admin-login__lead">
        Пароль задаётся переменной <code>ADMIN_PASSWORD</code> или берётся из
        <code>data/admin-password.txt</code>.
      </p>

      <label class="field">
        <span class="field__label">Пароль</span>
        <input
          ref="input"
          v-model="password"
          class="input"
          type="password"
          autocomplete="current-password"
          autofocus
          :disabled="busy"
          required
        />
      </label>

      <p v-if="error" class="alert" role="alert">{{ error }}</p>

      <button class="btn btn--primary admin-login__submit" type="submit" :disabled="busy">
        {{ busy ? 'Проверяем…' : 'Войти' }}
      </button>

      <a class="admin-login__back" href="/">← вернуться на сайт</a>
    </form>
  </div>
</template>
