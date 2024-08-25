import { SELF } from "cloudflare:test"
import { test, expect } from "vitest"

test("worker must responds with 401", async () => {
	const response = await SELF.fetch("https://example.com/diary/api/schedule")
	expect(response.status).toBe(401)
})

test("worker must response with 404", async () => {
	const response = await SELF.fetch("https://example.com/diary/api/unknown")
	expect(await response.text())
		.include("Страница не найдена")
		.include(
			"Страница, которую Вы пытаетесь посмотреть не найдена. Возможно, Вы ошиблись при наборе адреса или страница была удалена с сайта."
		)
	expect(response.status).toBe(404)
})
