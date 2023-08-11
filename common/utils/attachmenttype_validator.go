package utils

import "pi-inventory/common/consts"

func ValidateAttachmentType(givenType string) (bool, string) {
	switch givenType {
	case consts.ValidAttachmentType1:
		return true, consts.ValidAttachmentType1
	case consts.ValidAttachmentType2:
		return true, consts.ValidAttachmentType2
	case consts.ValidAttachmentType3:
		return true, consts.ValidAttachmentType3
	default:
		return false, ""
	}
}
