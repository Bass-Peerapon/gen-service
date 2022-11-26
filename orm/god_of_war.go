package orm

import "github.com/Bass-Peerapon/gen-service/models"

func OrmGodOfWar(ptr *models.GodOfWar, currentRow RowValue) (*models.GodOfWar, error) {
	v, err := fillValue(ptr, currentRow)
	if v != nil {
		return v.(*models.GodOfWar), nil
	}

	return nil, err
}
