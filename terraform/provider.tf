terraform {
  required_providers {
    linuxbox = {
      version = "~> 0.1"
      source  = "local/erwinwillems/linuxbox"
    }
  }
}

provider "linuxbox" {

}

