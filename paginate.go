package paginator

import "gorm.io/gorm"

// Paginate applies pagination to any GORM query with metadata.
func Paginate(model interface{}, query *gorm.DB, params Params, preload ...string) (interface{}, Metadata, error) {
	var total int64
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Model(model).Count(&total).Error; err != nil {
		return nil, Metadata{}, err
	}

	for _, pre := range preload {
		query = query.Preload(pre)
	}

	if params.OrderBy != "" {
		query = query.Order(params.OrderBy)
	}

	offset := (params.Page - 1) * params.Limit
	query = query.Limit(params.Limit).Offset(offset)

	resultSlice := makeSlice(model)
	if err := query.Find(resultSlice).Error; err != nil {
		return nil, Metadata{}, err
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return dereferenceSlice(resultSlice), Metadata{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: total,
		TotalPages:   totalPages,
		HasNext:      params.Page < totalPages,
		HasPrev:      params.Page > 1,
	}, nil
}
