package storage

func ListPopularTags() ([]string, error) {
	var res []string
	err := gormDB.Table("popular_tags").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
