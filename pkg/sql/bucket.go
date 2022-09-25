package kvsql

type SqlBucket struct {
	validTablename string
	keyMin         string
	keyMax         string
}

type SqlBucketBuilder func(bucket string) (SqlBucket, error)
