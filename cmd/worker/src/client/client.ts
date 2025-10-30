import type { Marks, Student } from "./types"
import { error } from "../utils"

let baseURL: string = ""

if (import.meta.env.DEV) {
	baseURL = "https://brand-new-kiasuo-app-indev.oddya.ru"
}

async function getWithAuth<T>(pathname: string): Promise<T> {
	try {
		const response = await fetch(baseURL + pathname, {
			headers: {
				"Telegram-Init": window.Telegram.WebApp.initData,
			},
		})

		switch (response.status) {
			case 200:
				return response.json()
			case 403:
				return error("Требуется обновление токена. Подробнее: /start")
			case 401:
				return error("Взлом?")
			default:
				return error("Неизвестная ошибка")
		}
	} catch {
		return error("Ошибка подключения")
	}
}

export function getStudent(): Promise<Student> {
	return getWithAuth("/internal/webapp/student")
}

export async function getMarks(): Promise<Marks> {
	return getWithAuth("/internal/webapp/marks")
}
