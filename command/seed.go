package command

import (
	"github.com/spf13/cobra"
	// attachmentSeeder "pi-inventory/modules/attachment/seed"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run Seed",
	Run: func(cmd *cobra.Command, args []string) {
		//gormDB := dic.Container.Get(commonConst.DbService).(*gorm.DB)

		//stockSeeder.CategorySeed(gormDB)
		//warehouseSeeder.WarehouseSeed(gormDB)
		//stockSeeder.UnitSeed(gormDB)
		//groupItemSeeder.VariantSeed(gormDB)
		//stockSeeder.StockSeed(gormDB)
		//stockSeeder.StockActivitySeed(gormDB)
		//stockSeeder.PurposeSeed(gormDB)
		// attachmentSeeder.AttachmentSeed(gormDB)
	},
}
