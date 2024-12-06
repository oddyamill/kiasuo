interface YandexError {
	errorMessage: string
	errorType: string
	stackTrace?: YandexErrorStackTrace[]
}

interface YandexErrorStackTrace {
	function: string
	file: string | null
	line: number | null
	column: number | null
}

export function throwYandexError(error: YandexError): never {
	const err = new Error(error.errorMessage)
	err.name = error.errorType

	if (error.stackTrace) {
		err.stack = renderStackTrace(err, error.stackTrace)
	}

	throw err
}

function renderStackTrace(error: Error, trace: YandexErrorStackTrace[]): string {
	return error.toString() + "\n" + trace.map(renderStackFrame).join("\n")
}

function renderStackFrame(frame: YandexErrorStackTrace): string {
	const location = frame.file !== null ? ` (${frame.file}:${frame.line}:${frame.column})` : ""
	return `    at ${frame.function}${location}`
}
