export interface ServerResponse<T> {
  error: boolean;
  msg: string;
  data: T;
}
