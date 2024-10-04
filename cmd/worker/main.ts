function proxy(url: URL, request: Request, cf?: CfProperties): Promise<Response> {
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

function proxyEdge(url: URL, request: Request, env: Env): Promise<Response>{
	return proxy(url, request, { resolveOverride: `cloudflare-edge-${env.EDGE}.oddya.ru` })
}

const CACHE_HEADER = "Worker-Cache", AUTH_HEADER = "Worker-Authorization"

async function proxyKiasuo(url: URL, request: Request, env: Env): Promise<Response> {
	const cf: CfProperties = {}

	if (request.headers.get(CACHE_HEADER) === "true") {
		if (request.headers.get(AUTH_HEADER) !== env.AUTH) {
			return new Response(null, { status: 401 })
		}

		request.headers.delete(CACHE_HEADER)
		request.headers.delete(AUTH_HEADER)
		cf.cacheEverything = true
		cf.cacheTtlByStatus = { "200-299": 86400 }
	}

	url.hostname = "dnevnik.kiasuo.ru"
	return proxy(url, request, cf)
}

export default {
	async fetch(request, env): Promise<Response> {
		const url = new URL(request.url)

		if (request.cf !== undefined && !["KJA", "KLD", "LED"].includes(request.cf.colo)) {
			return proxyEdge(url, request, env)
		}

		return proxyKiasuo(url, request, env)
	},
} satisfies ExportedHandler<Env>
