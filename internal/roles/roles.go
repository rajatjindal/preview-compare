package roles

import (
	"os"
)

type Role string

const (
	RolePreviewMain      Role = "preview-main"
	RolePreviewFirst     Role = "preview-1"
	RolePreviewSecond    Role = "preview-2"
	RolePreviewUndefined Role = "undefined"
)

func GetRole() Role {
	role := os.Getenv("app_role")
	if role == "" {
		return RolePreviewUndefined
	}

	return Role(role)
}
