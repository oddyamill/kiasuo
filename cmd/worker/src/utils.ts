export function humanizeLesson(lesson: string): string {
	switch (lesson.toLowerCase()) {
		case "основы безопасности жизнедеятельности":
		case "основы безопасности и защиты родины":
			return "ОБЖ"
		case "физическая культура":
			return "Физкультура"
		case "алгебра и начала математического анализа":
			return "Алгебра"
		case "искусственный интеллект":
			return "Информатика (ИИ)"
		case "россия - мои горизонты":
			return "Профориентация"
		case "учимся писать сочинение":
			return "Русский язык (сочинение)"
		case "трудные вопросы орфографии и пунктуации":
		case "секреты русского языка":
			return "Русский язык (консультация)"
		case "родной язык (русский язык)":
			return "Родной язык"
		case "труд (технология)":
			return "Труд"
		case "решение задач повышенной сложности":
			return "Математика (консультация)"
		case "основы духовно-нравственной культуры народов россии":
			return "ОДНКНР"
		case "изобразительное искусство":
			return "ИЗО"
		case "мировая художественная культура":
			return "МХК"
		case "основы религиозных культур и светской этики":
			return "ОРКСЭ"
		case "иностранный язык (английский язык)":
			return "Английский язык"
		default:
			return lesson
	}
}

export function error(message: string): Promise<any> {
	window.Telegram.WebApp.showAlert(message, () => window.Telegram.WebApp.close())
	// :cry:
	return new Promise(() => {})
}
