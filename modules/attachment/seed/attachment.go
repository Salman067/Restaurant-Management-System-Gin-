package seed

import (
	"fmt"
	"pi-inventory/modules/attachment/schema"
	"time"

	"gorm.io/gorm"

	attachmentConst "pi-inventory/modules/attachment/consts"

	fakeData "github.com/brianvoe/gofakeit/v6"
)

func AttachmentSeed(db *gorm.DB) {

	attachments := make([]*schema.Attachment, 0)

	for index := 0; index < 10; index++ {
		attachment := schema.Attachment{
			Name:          fakeData.Name(),
			Path:          fakeData.URL(),
			AttachmentKey: fakeData.UUID(),
			OwnerID:       1,
			CreatedAt:     time.Now(),
			CreatedBy:     1,
		}
		attachments = append(attachments, &attachment)
	}
	result := db.Table(attachmentConst.AttachmentTable).Create(&attachments)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
