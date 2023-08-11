package seed

// import (
// 	"fmt"
// 	"gorm.io/gorm"
// 	"pi-inventory/modules/Composite/schema"
// 	"time"

// 	CompositeConst "pi-inventory/modules/Composite/consts"

// 	fakeData "github.com/brianvoe/gofakeit/v6"
// )

// func CompositeSeed(db *gorm.DB) {

// 	composites := make([]*schema.Composite, 0)

// 	for index := 0; index < 10; index++ {
// 		Composite := schema.Composite{
// 			Title:     fakeData.Company(),
// 			Address: fakeData.Address().City,
// 			OwnerID: 1,
// 			CreatedAt: time.Now(),
// 			CreatedBy: 1,
// 		}
// 		composites = append(composites, &Composite)
// 	}
// 	result := db.Table(CompositeConst.CompositeTable).Create(&composites)
// 	if result.Error != nil {
// 		panic(fmt.Errorf("failed generating statement: %w", result.Error))
// 	}
// }
