package settings

import (
	"context"
	"database/sql"
	"errors"

	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/models"
)

type SettingRepository struct {
}

func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

func (r *SettingRepository) GetSettingByType(ctx context.Context, model string, modelId *int64, settingType enums.SettingType) (map[string]any, error) {
	q := xqb.Model[models.Setting]().
		WithContext(ctx).
		Where("type", "=", settingType)

	if modelId != nil && model != "" {
		q.Where("model", "=", model).
			Where("model_id", "=", *modelId)
	}

	settings, err := q.First()
	if err != nil {
		if errors.Is(err, xqb.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return settings.Settings, nil
}

func (r *SettingRepository) GetAllSettingsByType(ctx context.Context, settingType enums.SettingType) ([]models.Setting, error) {
	return xqb.Model[models.Setting]().
		WithContext(ctx).
		Where("type", "=", settingType).
		Get()
}

func (r *SettingRepository) SaveSetting(ctx context.Context, model string, modelId *int64, settingType enums.SettingType, settings map[string]any, tx *sql.Tx) error {
	values := map[string]any{
		"type":       settingType,
		"settings":   settings,
		"updated_at": time.Now(),
	}

	if modelId != nil && model != "" {
		values["model"] = model
		values["model_id"] = *modelId
	} else {
		values["model"] = "global"
		values["model_id"] = 0
	}

	_, err := xqb.Table("settings").
		WithContext(ctx).
		WithTx(tx).
		Upsert(
			[]map[string]any{values},
			[]string{"type", "model", "model_id"},
			[]string{"settings", "updated_at"},
		)
	if err != nil {
		return err
	}

	return nil
}
