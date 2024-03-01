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
