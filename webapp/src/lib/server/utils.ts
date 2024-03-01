export async function request<T = void>(url: string, method: string, body?: T) {
    const req: RequestInit = {
        method: method,
    };
    if (body !== undefined) {
        req.body = JSON.stringify(body);
    }
    const res = await fetch(url, req);
    return await res.json();
}
