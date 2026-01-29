terraform {
  required_version = "~>1.3"
  required_providers {
    csd = {
      source  = "idealo/csd"
      version = "~>2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~>6.28"
    }
  }
}

provider "aws" {
  region              = "eu-central-1"
  allowed_account_ids = ["433744410943"]
}

provider "csd" {}
