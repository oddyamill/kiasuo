import { AUTH_HEADER, CACHE_HEADER, CACHE_ROUTES, CACHE_TTL, ORIGIN_DOMAIN } from "./config"

function proxyRequest(url: URL, request: Request, cf?: CfProperties): Promise<Response> {
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

function proxyEdge(url: URL, request: Request, env: Env): Promise<Response> {
	return proxyRequest(url, request, { resolveOverride: `cloudflare-edge-${env.EDGE}.oddya.ru` })
}

async function proxyKiasuo(url: URL, request: Request, env: Env): Promise<Response> {
	const cf: CfProperties = {}

	if (request.headers.has(AUTH_HEADER)) {
		if (request.headers.get(AUTH_HEADER) !== env.AUTH) {
			return new Response(null, { status: 407 })
		}

		if (CACHE_ROUTES.includes(url.pathname) && request.headers.get(CACHE_HEADER) !== "no") {
			cf.cacheEverything = true
			cf.cacheTtlByStatus = { "200-299": CACHE_TTL }
		}
	}

	const headers = new Headers()

	if (request.headers.has("Accept-Encoding")) {
		headers.set("Accept-Encoding", request.headers.get("Accept-Encoding")!)
	}

	if (request.headers.has("Authorization")) {
		headers.set("Authorization", request.headers.get("Authorization")!)
	}

	if (request.headers.has("Content-Type")) {
		headers.set("Content-Type", request.headers.get("Content-Type")!)
		headers.set("Content-Length", request.headers.get("Content-Length") ?? "0")
	}

	// @ts-expect-error
	request.headers = headers
	url.hostname = ORIGIN_DOMAIN
	return proxyRequest(url, request, cf)
}

async function purgeCache(url: URL, request: Request, env: Env): Promise<Response> {
	if (request.method !== "POST") {
		return new Response(null, { status: 405 })
	}

	if (request.headers.get(AUTH_HEADER) !== env.AUTH) {
		return new Response(null, { status: 401 })
	}

	const response = await fetch(`https://api.cloudflare.com/client/v4/zones/${env.ZONE}/purge_cache`, {
		headers: {
			Authorization: "Bearer " + env.CLOUDFLARE,
			"Content-Type": "application/json",
		},
		method: "POST",
		body: JSON.stringify({
			files: CACHE_ROUTES.map((route) => `https://${ORIGIN_DOMAIN}${route}${url.search}`),
		}),
	})

	return new Response(null, { status: response.status })
}

export default {
	async fetch(request, env): Promise<Response> {
		const url = new URL(request.url)

		if (url.pathname === "/purge-cache") {
			return purgeCache(url, request, env)
		}

		if (request.cf !== undefined && !["KJA", "KLD", "LED"].includes(request.cf.colo)) {
			return proxyEdge(url, request, env)
		}

		return proxyKiasuo(url, request, env)
	},
} satisfies ExportedHandler<Env>
