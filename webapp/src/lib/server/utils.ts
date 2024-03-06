import type { Cookies } from "@sveltejs/kit";
import { COOKIE_MAX_AGE, SESSION_ID_HEADER } from "./constants";
import { dev } from "$app/environment";

type RequestOptions = {
    sessionId?: string;
    body?: unknown;
};

export async function request<Response = void>(url: string, method: string, opts?: RequestOptions) {
    const req: RequestInit = {
        method: method,
    };
    if (opts !== undefined && opts.body !== undefined) {
        req.body = JSON.stringify(opts.body);
    }
    if (opts !== undefined && opts.sessionId !== undefined) {
        req.headers = {
            [SESSION_ID_HEADER]: opts.sessionId,
        };
    }
    try {
        const res = await fetch(url, req);
        return (await res.json()) as Response;
    } catch (error) {
        console.error("request error", error);
        return undefined;
    }
}

export function setCookie(cookies: Cookies, key: string, value: string) {
    cookies.set(key, value, {
        path: "/",
        httpOnly: true,
        sameSite: "strict",
        secure: !dev,
        maxAge: COOKIE_MAX_AGE,
    });
}
