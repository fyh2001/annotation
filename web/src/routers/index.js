import { createRouter, createWebHashHistory, } from "vue-router";

const routes = [
    {
        path: "",
        name: "Workbanch",
        component: () => import("../views/workbench/index.vue")
    }
]

const router = createRouter({
    history: createWebHashHistory("/annotation"),
    routes,
});

export default router;
