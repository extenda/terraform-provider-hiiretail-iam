package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// groupNameValidator validates an IAM group name according to HiiRetail's rules
type groupNameValidator struct {
}

// Description returns a plain text description of the validator's behavior
func (v groupNameValidator) Description(_ context.Context) string {
	return "group name must be between 3 and 128 characters, start with a letter, and contain only letters, numbers, hyphens, and underscores"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior
func (v groupNameValidator) MarkdownDescription(_ context.Context) string {
	return "Group name must be between 3 and 128 characters, start with a letter, and contain only letters, numbers, hyphens (`-`), and underscores (`_`)."
}

// ValidateString performs the validation
func (v groupNameValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	pattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{2,127}$`)

	if !pattern.MatchString(value) {
		response.Diagnostics.AddError(
			"Invalid Group Name",
			fmt.Sprintf(
				"Group name %q is invalid. Group names must be between 3 and 128 characters, start with a letter, and contain only letters, numbers, hyphens (-), and underscores (_).",
				value,
			),
		)
	}
}

// groupDescriptionValidator validates an IAM group description
type groupDescriptionValidator struct {
}

// Description returns a plain text description of the validator's behavior
func (v groupDescriptionValidator) Description(_ context.Context) string {
	return "group description must be between 1 and 256 characters"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior
func (v groupDescriptionValidator) MarkdownDescription(_ context.Context) string {
	return "Group description must be between 1 and 256 characters."
}

// ValidateString performs the validation
func (v groupDescriptionValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	if len(value) < 1 || len(value) > 256 {
		response.Diagnostics.AddError(
			"Invalid Group Description",
			fmt.Sprintf(
				"Group description length must be between 1 and 256 characters, got %d characters.",
				len(value),
			),
		)
	}
}