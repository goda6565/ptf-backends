resource "aws_ssm_parameter" "ssm_parameter" {
  for_each = var.secrets
  name  = "/${var.app_name}/${each.key}"
  type  = "secure_string"
  value = "uninitialized"
  lifecycle {
    ignore_changes = [value]
  }
}