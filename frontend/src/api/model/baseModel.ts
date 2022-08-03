export interface BasicPageParams {
  page: number;
  pageSize: number;
}

export interface BasicFetchResult<T> {
  items: T[];
  total: number;
}

export interface ResponseData<T> {
  code: number;
  msg: string;
  data: T;
}
