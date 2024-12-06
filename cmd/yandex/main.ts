import type { Handler } from "@yandex-cloud/function-types"

const ORIGIN_DOMAIN = "dnevnik.kiasuo.ru"

const handler: Handler.Http = async (event) => {
	let origin: URL

	try {
		origin = new URL(event.queryStringParameters.origin || event.headers["Origin"])
	} catch {
		return {
			statusCode: 400,
			body: "Invalid origin",
		}
	}

	if (origin.hostname !== ORIGIN_DOMAIN) {
		return {
			statusCode: 403,
			body: "Host not allowed",
		}
	}

	const headers = new Headers()

	if (event.headers["Kiasuo-Authorization"]) {
		headers.set("Authorization", event.headers["Kiasuo-Authorization"])
	}

	if (event.httpMethod === "POST") {
		headers.set("Content-Type", event.headers["Content-Type"] || "application/json")
		headers.set("Content-Length", event.headers["Content-Length"] || "0")
	}

	const response = await fetch(origin, {
		headers,
		method: event.httpMethod,
		// ookay it works
		body: event.body ? Buffer.from(event.body, event.isBase64Encoded ? "base64" : "utf8").toString() : undefined,
	})

	return {
		statusCode: response.status,
		body: response.ok ? await response.text() : undefined,
		headers: {
			"Content-Encoding": response.headers.get("Content-Encoding") ?? "identity",
			"Content-Type": response.headers.get("Content-Type") ?? "text/plain",
			"Content-Length": response.headers.get("Content-Length") ?? "0",
		},
	}
}

module.exports = { handler }
