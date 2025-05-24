package platform

import (
	managementv1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
)

func GetUserName(self *managementv1.Self) string {
	if self.Status.User != nil {
		return self.Status.User.Name
	}

	if self.Status.Team != nil {
		return self.Status.Team.Name
	}

	return self.Status.Subject
}
