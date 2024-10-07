import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/home', name: 'Home', component: () => import('../views/Home.vue') },
    { path: '/edit', name: 'Edit', component: () => import('../views/Edit.vue') },
    { path: '/', name: 'Login', component: () => import('../views/Login.vue') },
  ]
})

export default router
