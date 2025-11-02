package activity

import (
	"context"
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letterBytes[b[i]%byte(len(letterBytes))]
	}
	return string(b)
}

func CopyTableActivity(ctx context.Context, table string) error {
	// 1. Postgres
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("postgres connect: %w", err)
	}
	defer db.Close()

	// 2. query
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	// 3. local CSV
	tmpFile := fmt.Sprintf("/tmp/%s_%d_%s.csv", table, time.Now().Unix(), randString(6))
	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("create tmp csv: %w", err)
	}
	defer os.Remove(tmpFile)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	cols, _ := rows.Columns()
	if err := w.Write(cols); err != nil {
		return fmt.Errorf("write headers: %w", err)
	}

	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		strRow := make([]string, len(row))
		for i, v := range row {
			strRow[i] = fmt.Sprint(v)
		}
		if err := w.Write(strRow); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	// 4. flush & rewind file before upload
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("csv flush: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("rewind file: %w", err)
	}

	// 5. MinIO upload
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: os.Getenv("S3_ENDPOINT")}, nil
				})),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("S3_KEY"),
				os.Getenv("S3_SECRET"),
				"")),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	bucket := os.Getenv("S3_BUCKET")
	key := fmt.Sprintf("sync-loop/%s_%d.csv", table, time.Now().Unix())
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("s3 upload: %w", err)
	}
	return nil
}