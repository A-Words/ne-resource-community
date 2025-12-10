import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import Upload from '@/pages/Upload.vue'
import Auth from '@/pages/Auth.vue'
import ResourceDetail from '@/pages/ResourceDetail.vue'
import Dashboard from '@/pages/Dashboard.vue'
import AdminAudit from '@/pages/AdminAudit.vue'
import AdminReports from '@/pages/AdminReports.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/upload', component: Upload },
    { path: '/auth', component: Auth },
    { path: '/resources/:id', component: ResourceDetail, props: true },
    { path: '/dashboard', component: Dashboard },
    { path: '/admin/audit', component: AdminAudit },
    { path: '/admin/reports', component: AdminReports },
  ],
})

export default router
