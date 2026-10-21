import { createApp } from 'vue'
import App from './App.vue'
import './styles.css'

// /admin — отдельный бандл: посетителям сайта код редактора не грузится.
const path = window.location.pathname.replace(/\/+$/, '')

if (path.endsWith('/admin')) {
  void import('./admin/AdminApp.vue').then(({ default: AdminApp }) => {
    createApp(AdminApp).mount('#app')
  })
} else {
  createApp(App).mount('#app')
}
