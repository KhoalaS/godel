export interface List<T extends object> {
  rows: Row<T>[]
  columns: (Column | TypedColumn<T>)[]
}

export interface Row<T> {
  data: T
  id: string | number
}

export type TypedColumn<T extends object> = Column & {
  key: keyof T
}

export interface Column {
  key: string
  width: number
  headerName: string
}
