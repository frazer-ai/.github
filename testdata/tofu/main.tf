# A fixture for the shared OpenTofu workflow: no providers, no credentials.

variable "name" {
  description = "Name to greet"
  type        = string
  default     = "Frazer"
}

resource "terraform_data" "greeting" {
  input = "hello ${var.name}"
}

output "greeting" {
  description = "The greeting"
  value       = terraform_data.greeting.output
}
