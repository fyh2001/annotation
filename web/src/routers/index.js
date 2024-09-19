import { createRouter, createWebHashHistory, createWebHistory } from "vue-router";

const routes = [
    {
        path: "",
        name: "Workbanch",
        component: () => import("../views/workbench/index.vue")
    }
]

const router = createRouter({
    // history: createWebHashHistory("/annotation"),
    history: createWebHashHistory(),
    routes,
});

export default router;
