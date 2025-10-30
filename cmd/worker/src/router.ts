import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
		{
			path: "/webapp/marks",
			component: () => import("./views/MarksView.vue"),
		},
  ],
})

export default router
