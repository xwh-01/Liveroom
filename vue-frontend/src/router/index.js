import { createRouter, createWebHistory } from 'vue-router'
import LiveRoom from '../views/LiveRoom.vue'

const routes = [
  { path: '/', name: 'LiveRoom', component: LiveRoom }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
