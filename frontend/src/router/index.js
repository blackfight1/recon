import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Targets from '../views/Targets.vue'
import Subdomains from '../views/Subdomains.vue'
import Changes from '../views/Changes.vue'
import Tasks from '../views/Tasks.vue'
import QuickScan from '../views/QuickScan.vue'

const routes = [
    {
        path: '/',
        name: 'Dashboard',
        component: Dashboard
    },
    {
        path: '/quick-scan',
        name: 'QuickScan',
        component: QuickScan
    },
    {
        path: '/targets',
        name: 'Targets',
        component: Targets
    },
    {
        path: '/subdomains',
        name: 'Subdomains',
        component: Subdomains
    },
    {
        path: '/changes',
        name: 'Changes',
        component: Changes
    },
    {
        path: '/tasks',
        name: 'Tasks',
        component: Tasks
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
