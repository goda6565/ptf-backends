resource "aws_ssm_parameter" "ssm_parameter" {
  for_each = toset(var.secrets) # 重複を防ぐためsetに変換するのも良いです
  name     = "/${var.app_name}/${each.value}"
  type     = "secure_string"
  value    = "uninitialized"
  lifecycle {
    ignore_changes = [value]
  }
}
