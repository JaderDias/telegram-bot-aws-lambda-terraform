resource "aws_s3_bucket" "bucket" {
  bucket = "${terraform.workspace}-telegram-bot-${random_pet.this.id}"
  tags = {
    environment = terraform.workspace
  }
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.bucket.id
  acl    = "private"
}
