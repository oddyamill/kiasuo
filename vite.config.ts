import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import vueDevTools from "vite-plugin-vue-devtools"
import { cloudflare } from "@cloudflare/vite-plugin"

export default defineConfig({
	plugins: [
		vue(),
		vueDevTools(),
		cloudflare()
	],
	server: {
		host: "0.0.0.0",
		port: 8787,
		allowedHosts: ["brand-new-kiasuo-webapp-indev.oddya.ru"],
	},
	build: {
		sourcemap: true,
	},
})
