locals {
  tags = merge(
    var.tags,
    {
      "provisioner" = "terraform"
  })
}
