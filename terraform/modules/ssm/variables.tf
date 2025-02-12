variable "app_name" {
  type        = string
  description = "The name of the application"
}

variable "secrets" {
  type        = map(string)
  description = "The names of the secret"
}