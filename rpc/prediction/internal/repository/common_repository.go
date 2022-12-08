package repository

import (
	"github.com/go-playground/validator/v10"
	"github.com/maniizu3110/attendance/rpc/prediction/external/mycache"
	"github.com/maniizu3110/attendance/rpc/prediction/internal/repository/query"
	"gorm.io/gorm"
)

const (
	NOT_FOUND_ERROR string = "record not found"
)

type CommonRepository[T any] interface {
	GetList(config query.FetchConfig) ([]*T, int64, error)
	GetByID(id uint64, preload ...string) (*T, error)
	Create(data *T) (*T, error)
	Update(id uint64, data *T) (*T, error)
	SoftDelete(id uint64, data *T) error
	HardDelete(id uint64, data *T) error
	Restore(id uint64, data *T) (*T, error)
	BatchCreate(data []*T) ([]*T, error)
}

const (
	DefaultLimit int = 10
	MaxLimit     int = 50
)

// repositoryで扱う基本的な構造体
type commonRepository[T any] struct {
	db        *gorm.DB
	Validator *validator.Validate
	cache     mycache.Cache
}

func NewCommonRepository[T any](db *gorm.DB, validator *validator.Validate, cache mycache.Cache) CommonRepository[T] {
	res := &commonRepository[T]{
		db:        db.Session(&gorm.Session{}), // reusable
		Validator: validator,
		cache:     cache,
	}
	return res
}

// Fetch: FetchConfigの内容に従った情報を全て取得
func (cr *commonRepository[T]) GetList(fetchConfig query.FetchConfig) ([]*T, int64, error) {
	var (
		dataNum int64          // 取得データ総数(offset/limitの設定がない場合に取得される件数.paginationに使用)
		data    []*T  = []*T{} // 返却データ
	)
	model := new(T)

	db, err := applyFetchConfig(*model, cr.db, fetchConfig, cr.cache) // FetchConfigをgormのDBに適用
	if err != nil {
		return nil, 0, err
	}

	err = implementGetList(db, &data, &dataNum, fetchConfig)
	if err != nil {
		return nil, 0, err
	}

	return data, dataNum, nil
}

// GetByID: 指定IDを持つデータとそれに紐づく指定preloadのデータを取得（論理削除対象も含む）
func (cr *commonRepository[T]) GetByID(id uint64, preload ...string) (*T, error) {
	data := new(T)
	db := cr.db.Unscoped()
	// DBにPreload適用
	db, err := query.ApplyPreloads(data, preload, db, cr.cache)
	if err != nil {
		return nil, err
	}
	// データの取得
	if err := db.Unscoped().Where("id = ?", id).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Create: データの新規作成
func (cr *commonRepository[T]) Create(data *T) (*T, error) {
	err := cr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}
		err := cr.Validator.Struct(data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Update: 既存データの内容の更新
func (cr *commonRepository[T]) Update(id uint64, data *T) (*T, error) {
	err := cr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(data).Where("id = ?", id).Updates(data).Error; err != nil {
			return err
		}
		if err := tx.Model(data).Where("id = ?", id).First(data).Error; err != nil {
			return err
		}
		err := cr.Validator.Struct(data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SoftDelete:指定データを論理削除（復元可）する．
func (cr *commonRepository[T]) SoftDelete(id uint64, data *T) error {
	if err := cr.db.Delete(data).Where("id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// HardDelete:指定データを物理削除（復元不可）する．
func (cr *commonRepository[T]) HardDelete(id uint64, data *T) error {
	if err := cr.db.Unscoped().Delete(data).Where("id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// BatchCreate: データの複数作成
func (r *commonRepository[T]) BatchCreate(data []*T) ([]*T, error) {
	var res []*T
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range data {
			if err := tx.Create(&v).Error; err != nil {
				return nil
			}
			res = append(res, v)
		}
		return nil
	})
	return res, err
}

// Restore:論理削除されたデータを復元する．
func (cr *commonRepository[T]) Restore(id uint64, data *T) (*T, error) {
	err := cr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Model(data).Update("deleted_at", nil).Error; err != nil {
			return err
		}
		err := cr.Validator.Struct(data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// FetchConfigの内容をgormのdbに適用
// FetchConfigを適用する順番を変えると挙動が変わることに注意
func applyFetchConfig[T any](model T, db *gorm.DB, fetchConfig query.FetchConfig, cache mycache.Cache) (*gorm.DB, error) {
	db = db.Model(&model)

	if fetchConfig.WithTrashed {
		db = db.Unscoped()
	}

	db, err := query.ApplyOrders(model, fetchConfig.Order, db) // Order適用
	if err != nil {
		return nil, err
	}

	db, err = query.ApplyPreloads(&model, fetchConfig.Preload, db, cache) // Preload適用
	if err != nil {
		return nil, err
	}

	limit := DefaultLimit
	if (fetchConfig.Limit != 0) && (fetchConfig.Limit < MaxLimit) {
		limit = fetchConfig.Limit
	}
	db = db.Limit(limit).Offset(fetchConfig.Offset)                  // Limit, Offset適用
	db, err = query.ApplyJoins(&model, fetchConfig.Joins, db, cache) // Join適用
	if err != nil {
		return nil, err
	}
	db, err = query.ApplyQueries(model, fetchConfig, db) // Query適用
	if err != nil {
		return nil, err
	}

	return db, nil
}

func implementGetList[T any](db *gorm.DB, data *[]*T, dataNum *int64, fetchConfig query.FetchConfig) error {
	if len(fetchConfig.Joins) > 0 {
		if err := db.Find(data).Error; err != nil {
			return err
		}
		if err := db.Count(dataNum).Error; err != nil {
			return err
		}
	} else {
		if err := db.Find(data).Count(dataNum).Error; err != nil {
			return err
		}
	}
	return nil
}
