import {
	AUTH_HEADER,
	CACHE_HEADER,
	CACHE_TTL,
	CACHED_ROUTES,
	EDGE_HEADER,
	EDGES,
	ORIGIN_DOMAIN,
} from "./config"
import { throwYandexError } from "./errors"

function proxyRequest(url: URL, request: Request, cf: CfProperties, headers?: Headers): Promise<Response> {
	const init: RequestInit = {
		headers: headers ?? request.headers,
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
	if (env.YANDEX !== undefined) {
		return proxyKiasuo(url, request, env, true)
	}

	const edge = (request.headers.get(EDGE_HEADER) || env.EDGE).toLowerCase()
	return proxyRequest(url, request, { resolveOverride: `cloudflare-edge-${edge}.oddya.ru` })
}

async function proxyKiasuo(url: URL, request: Request, env: Env, yandex?: boolean): Promise<Response> {
	const cf: CfProperties = {}

	if (request.headers.has(AUTH_HEADER)) {
		if (request.headers.get(AUTH_HEADER) !== env.AUTH) {
			return new Response(null, { status: 407 })
		}

		if (CACHED_ROUTES.includes(url.pathname) && request.headers.get(CACHE_HEADER) !== "no") {
			cf.cacheEverything = true
			cf.cacheTtlByStatus = { "200-299": CACHE_TTL }
		}
	}

	const headers = new Headers()

	if (request.headers.has("Accept")) {
		headers.set("Accept", request.headers.get("Accept")!)
	}

	if (request.headers.has("Authorization")) {
		headers.set((yandex ? "Kiasuo-" : "") + "Authorization", request.headers.get("Authorization")!)
	}

	if (request.headers.has("Content-Type")) {
		headers.set("Content-Type", request.headers.get("Content-Type")!)
		headers.set("Content-Length", request.headers.get("Content-Length") ?? "0")
	}

	url.hostname = ORIGIN_DOMAIN

	if (yandex) {
		const origin = url.toString()
		url = new URL(env.YANDEX)
		url.searchParams.set("origin", origin)
	}

	const response = await proxyRequest(url, request, cf, headers)

	if (response.status === 502 && response.headers.get("X-Function-Error") === "true") {
		throwYandexError(await response.json())
	}

	return response
}

async function purgeCache(url: URL, request: Request, env: Env): Promise<Response> {
	if (request.method !== "POST") {
		return new Response(null, { status: 405 })
	}

	if (request.headers.get(AUTH_HEADER) !== env.AUTH) {
		return new Response(null, { status: 401 })
	}

	const urls = CACHED_ROUTES.map((route) => `https://${ORIGIN_DOMAIN}${route}${url.search}`)

	if (env.YANDEX !== undefined) {
		urls.push(...urls.map((origin) => `${env.YANDEX}?origin=${encodeURIComponent(origin)}`))
	}

	const response = await fetch(`https://api.cloudflare.com/client/v4/zones/${env.ZONE}/purge_cache`, {
		headers: {
			Authorization: (env.CLOUDFLARE[6] !== " " ? "Bearer " : "") + env.CLOUDFLARE,
			"Content-Type": "application/json",
		},
		method: "POST",
		body: JSON.stringify({ files: urls }),
	})

	return new Response(null, { status: response.status })
}

export default {
	fetch(request, env): Promise<Response> {
		const url = new URL(request.url)

		if (url.pathname === "/internal/purge-cache") {
			return purgeCache(url, request, env)
		}

		if (request.cf !== undefined && !EDGES.includes(request.cf.colo)) {
			return proxyEdge(url, request, env)
		}

		return proxyKiasuo(url, request, env)
	},
} satisfies ExportedHandler<Env>
