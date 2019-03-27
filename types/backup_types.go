package types

type BackupConfig struct {
	// Backup interval in hours
	IntervalHours int `yaml:"interval_hours" json:"intervalHours,omitempty" norman:"default=12"`
	// Number of backups to keep
	Retention int `yaml:"retention" json:"retention,omitempty" norman:"default=6"`
	// s3 target
	S3BackupConfig *S3BackupConfig `yaml:",omitempty" json:"s3BackupConfig"`
}

type S3BackupConfig struct {
	// Access key ID
	AccessKey string `yaml:"access_key" json:"accessKey,omitempty"`
	// Secret access key
	SecretKey string `yaml:"secret_key" json:"secretKey,omitempty" norman:"required,type=password" `
	// name of the bucket to use for backup
	BucketName string `yaml:"bucket_name" json:"bucketName,omitempty"`
	// AWS Region, AWS spcific
	Region string `yaml:"region" json:"region,omitempty"`
	// Endpoint is used if this is not an AWS API
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}
