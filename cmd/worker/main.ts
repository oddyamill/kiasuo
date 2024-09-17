function proxy(url: URL, request: Request, cf?: CfProperties) {
	const init: RequestInit = {
		headers: request.headers,
		method: request.method,
		redirect: "manual",
		cf,
	}

	if (request.method === "POST") {
		init.body = request.body
	}

	return fetch(url, init)
}

function proxyEdge(url: URL, request: Request, env: Env) {
	return proxy(url, request, { resolveOverride: `cloudflare-edge-${env.EDGE}.oddya.ru` })
}

async function proxyKiasuo(url: URL, request: Request) {
	url.hostname = "dnevnik.kiasuo.ru"
	return proxy(url, request)
}

export default {
	async fetch(request, env): Promise<Response> {
		const url = new URL(request.url)

		if (request.cf !== undefined && !["KJA", "KLD", "LED"].includes(request.cf.colo)) {
			return proxyEdge(url, request, env)
		}

		return proxyKiasuo(url, request)
	},
} satisfies ExportedHandler<Env>
