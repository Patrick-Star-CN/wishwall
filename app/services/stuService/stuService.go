package stuService

import (
	"wishwall/app/models"
	"wishwall/config/database"
)

func GetStu(num string) (*models.Stu, error) {
	var stu *models.Stu

	result := database.DB.Where(
		&models.Stu{
			Num: num,
		}).Find(&stu)
	if result.Error != nil {
		return nil, result.Error
	}
	return stu, nil
}
