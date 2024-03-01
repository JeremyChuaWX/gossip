export type Response<T> =
    | ({ error: false; message: string } & T)
    | { error: true; message: string };

export type User = {
    userId: string;
    username: string;
};
