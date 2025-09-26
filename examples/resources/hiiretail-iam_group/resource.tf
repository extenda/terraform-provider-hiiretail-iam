# Basic group creation
resource "hiiretail-iam_group" "basic" {
  name        = "developers"
  description = "Development team group"
}

# Comprehensive group example
resource "hiiretail-iam_group" "complete" {
  name        = "platform-team"
  description = "Platform engineering team with access to infrastructure services"
}