import { createRouter, createWebHistory } from 'vue-router'
import About from '../src/views/About.vue'
import Widget from '../src/views/Widget.vue';
import MoscowTalks from '../src/views/MoscowTalks.vue';

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: About
        },
        {
            path: '/widget',
            component: Widget
        },
        {
            path: '/moscowTalks',
            component: MoscowTalks
        },
    ],
})

export default router