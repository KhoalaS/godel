package types

type DownloadJobTransformer func(job DownloadJob) (DownloadJob, error)
