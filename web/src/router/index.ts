import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import Upload from '@/pages/Upload.vue'
import Auth from '@/pages/Auth.vue'
import ResourceDetail from '@/pages/ResourceDetail.vue'
import Dashboard from '@/pages/Dashboard.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/upload', component: Upload },
    { path: '/auth', component: Auth },
    { path: '/resources/:id', component: ResourceDetail, props: true },
    { path: '/dashboard', component: Dashboard },
  ],
})

export default router
