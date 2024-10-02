package mysqlaccesscontrol

import (
	"gameApp/entity"
	"gameApp/repository/mysql"
	"time"
)

func scanPermission(scanner mysql.Scanner) (entity.Permission, error) {
	var createdAt time.Time
	per := entity.Permission{}
	err := scanner.Scan(&per.ID, &per.Title, &createdAt)

	return per, err

}
