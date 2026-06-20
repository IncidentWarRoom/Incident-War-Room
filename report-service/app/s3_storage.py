import os

import boto3
from botocore.client import Config
from dotenv import load_dotenv


load_dotenv()


class S3Storage:
    def __init__(self):
        self.endpoint_url = os.getenv("S3_ENDPOINT_URL")
        self.region = os.getenv("S3_REGION")
        self.bucket_name = os.getenv("S3_BUCKET_NAME")
        self.access_key = os.getenv("S3_ACCESS_KEY")
        self.secret_key = os.getenv("S3_SECRET_KEY")
        self.public_base_url = os.getenv("S3_PUBLIC_BASE_URL")

        self.client = boto3.client(
            "s3",
            endpoint_url=self.endpoint_url,
            region_name=self.region,
            aws_access_key_id=self.access_key,
            aws_secret_access_key=self.secret_key,
            config=Config(signature_version="s3v4"),
        )

    def upload_pdf(self, pdf_bytes: bytes, object_name: str) -> str:
        self.client.put_object(
            Bucket=self.bucket_name,
            Key=object_name,
            Body=pdf_bytes,
            ContentType="application/pdf",
            ACL="public-read",
        )

        return f"{self.public_base_url.rstrip('/')}/{object_name}"