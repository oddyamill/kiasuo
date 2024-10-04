import { fetchMock, SELF } from "cloudflare:test"
import { test, expect, beforeAll, afterEach } from "vitest"
import { AUTH_HEADER } from "./config"

let mocked = false

beforeAll(async () => {
	const response = await fetch("https://1.1.1.1/cdn-cgi/trace")

	if ((mocked = !(await response.text()).includes("loc=RU"))) {
		fetchMock.activate()
		fetchMock.disableNetConnect()
	}

	console.log("mocked:", mocked)
})

afterEach(() => mocked && fetchMock.assertNoPendingInterceptors())

test("worker must responds with 401", async () => {
	fetchMock
		.get("https://dnevnik.kiasuo.ru")
		.intercept({ path: "/diary/api/schedule" })
		.reply(401)

	const response = await SELF.fetch("https://example.com/diary/api/schedule")
	expect(response.status).toBe(401)
})

test("worker must response with 404", async () => {
	fetchMock
		.get("https://dnevnik.kiasuo.ru")
		.intercept({ path: "/diary/api/unknown" })
		.reply(404, `
			<body>
					<div class="dialog">
							<h1>Страница не найдена</h1>
							<p>Страница, которую Вы пытаетесь посмотреть не найдена. Возможно, Вы ошиблись при наборе адреса или страница была удалена с сайта.</p>
					</div>
			</body>
		`)

	const response = await SELF.fetch("https://example.com/diary/api/unknown")
	expect(await response.text())
		.include("Страница не найдена")
		.include(
			"Страница, которую Вы пытаетесь посмотреть не найдена. Возможно, Вы ошиблись при наборе адреса или страница была удалена с сайта."
		)
	expect(response.status).toBe(404)
})

test("worker must response with 407", async () => {
	const response = await SELF.fetch("https://example.com/diary/api/recipients", {
		headers: {
			[AUTH_HEADER]: "123456",
		},
	})
	expect(response.status).toBe(407)
})
