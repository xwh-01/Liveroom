import { createRouter, createWebHistory } from 'vue-router'
import RoomLobby from '../views/RoomLobby.vue'
import LiveRoom from '../views/LiveRoom.vue'

const routes = [
  { path: '/', redirect: '/rooms' },
  { path: '/rooms', name: 'RoomLobby', component: RoomLobby },
  { path: '/room/:roomId', name: 'LiveRoom', component: LiveRoom, props: true }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
