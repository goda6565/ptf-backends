terraform {
  backend "s3" {
    bucket         = "ptf-terraform-tfstate-backends"
    key            = "terraform.tfstate"
    region         = "ap-northeast-1"
    dynamodb_table = "ptf-terraform-tfstate-locking"
    encrypt        = true // ステートファイルを暗号化する
  }
}