import type { Cookies } from "@sveltejs/kit";
import { COOKIE_MAX_AGE } from "./constants";
import { dev } from "$app/environment";

export async function request<Response = void>(url: string, method: string, body?: unknown) {
    const req: RequestInit = {
        method: method,
    };
    if (body !== undefined) {
        req.body = JSON.stringify(body);
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
