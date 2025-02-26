.PHONY: run-minio stop-minio restart-minio

# Start MinIO container
run-minio:
	docker run -d --name minio \
	  -p 9000:9000 -p 9001:9001 \
	  -e "MINIO_ROOT_USER=admin" \
	  -e "MINIO_ROOT_PASSWORD=secretpass" \
	  quay.io/minio/minio server /data --console-address ":9001"

# Stop and remove MinIO container
stop-minio:
	docker stop minio || true
	docker rm minio || true

# Restart MinIO container
restart-minio: stop-minio run-minio
