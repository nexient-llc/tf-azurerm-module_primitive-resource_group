variable "resource_group" {
  description = "resource group"
  type = object({
    location = string
  })
}

variable "resource_group_name" {
  description = "name of the resource"
  type        = string
  nullable    = false
}

variable "tags" {
  description = "tags to be assigned to child resource"
  type        = map(string)
  default     = {}
}