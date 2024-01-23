package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HydrologyStatsRepository struct {
	dbPool *pgxpool.Pool
}

func NewHydrologyStatsRepository(pool *pgxpool.Pool) *HydrologyStatsRepository {
	return &HydrologyStatsRepository{dbPool: pool}
}

func (r *HydrologyStatsRepository) AddControlValue(ctx context.Context, value model.ControlValue) error {

	sql := `INSERT INTO control_values (id, post_code, type, date_start, value)
            VALUES ($1, $2, $3, $4, $5)
            ON CONFLICT (post_code, type, date_start) DO NOTHING;`

	commandTag, err := r.dbPool.Exec(ctx, sql, value.ID, value.PostCode, value.Type, value.DateStart, value.Value)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("control value already exists")
	}

	return nil
}

func (r *HydrologyStatsRepository) RemoveControlValue(ctx context.Context, id string) error {

	sql := `DELETE FROM control_values WHERE id = $1;`

	commandTag, err := r.dbPool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no control value found with id: %s", id)
	}

	return nil
}

func (r *HydrologyStatsRepository) UpdateControlValues(ctx context.Context, values []model.ControlValue) error {

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, value := range values {
		sql := `UPDATE control_values SET post_code = $1, type = $2, date_start = $3, value = $4 WHERE id = $5 AND NOT EXISTS (
                    SELECT 1 FROM control_values WHERE post_code = $1 AND type = $2 AND date_start = $3 AND id != $5
                );`

		commandTag, err := tx.Exec(ctx, sql, value.PostCode, value.Type, value.DateStart, value.Value, value.ID)
		if err != nil {
			return err
		}

		if commandTag.RowsAffected() == 0 {
			return fmt.Errorf("update failed for control value with id: %s", value.ID)
		}
	}

	return tx.Commit(ctx)
}

func (r *HydrologyStatsRepository) GetControlValues(ctx context.Context, postCode string, controlType model.ControlValueType, page, pageSize int) ([]model.ControlValue, int, error) {

	var controlValues []model.ControlValue
	var totalCount int

	query := `SELECT id, post_code, type, date_start, value FROM control_values WHERE post_code = $1 AND type = $2 ORDER BY date_start DESC LIMIT $3 OFFSET $4`
	rows, err := r.dbPool.Query(ctx, query, postCode, controlType, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cv model.ControlValue
		if err := rows.Scan(&cv.ID, &cv.PostCode, &cv.Type, &cv.DateStart, &cv.Value); err != nil {
			return nil, 0, err
		}
		controlValues = append(controlValues, cv)
	}

	countQuery := `SELECT COUNT(*) FROM control_values WHERE post_code = $1 AND type = $2`
	err = r.dbPool.QueryRow(ctx, countQuery, postCode, controlType).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return controlValues, totalCount, nil
}

func (r *HydrologyStatsRepository) GetControlValuesByDay(ctx context.Context, postCode string, date time.Time) ([]model.ControlValue, error) {

	var controlValues []model.ControlValue

	query := `
	SELECT id, post_code, type, date_start, value
	FROM (
    SELECT DISTINCT ON (type) id, post_code, type, date_start, value
    FROM control_values
    WHERE post_code = $1 AND date_start <= $2
    ORDER BY type, date_start DESC, value DESC
	) AS subquery
	ORDER BY value DESC;
    `

	rows, err := r.dbPool.Query(ctx, query, postCode, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cv model.ControlValue
		if err := rows.Scan(&cv.ID, &cv.PostCode, &cv.Type, &cv.DateStart, &cv.Value); err != nil {
			return nil, err
		}
		controlValues = append(controlValues, cv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return controlValues, nil
}

func (r *HydrologyStatsRepository) GetControlValuesByDateInterval(ctx context.Context, postCode string, dateStart, dateEnd time.Time) ([]model.ControlValue, error) {
	var controlValues []model.ControlValue

	query := `
    WITH latest_before_start AS (
        SELECT DISTINCT ON (type) id, post_code, type, date_start, value
        FROM control_values
        WHERE post_code = $1 AND date_start < $2
        ORDER BY type, date_start DESC
    )
    SELECT id, post_code, type, date_start, value
    FROM control_values
    WHERE post_code = $1 AND date_start BETWEEN $2 AND $3
    UNION ALL
    SELECT * FROM latest_before_start
    ORDER BY date_start
    `

	rows, err := r.dbPool.Query(ctx, query, postCode, dateStart, dateEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cv model.ControlValue
		if err := rows.Scan(&cv.ID, &cv.PostCode, &cv.Type, &cv.DateStart, &cv.Value); err != nil {
			return nil, err
		}
		controlValues = append(controlValues, cv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return controlValues, nil
}

func (r *HydrologyStatsRepository) AddWaterlevel(ctx context.Context, value model.Waterlevel) error {
	sqlUpdate := `
        UPDATE waterlevels SET waterlevel = $1 WHERE date = $2 AND post_code = $3;
    `
	result, err := r.dbPool.Exec(ctx, sqlUpdate, value.Waterlevel, value.Date, value.PostCode)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		sqlInsert := `
            INSERT INTO waterlevels (id, post_code, date, waterlevel) VALUES ($1, $2, $3, $4);
        `
		_, err := r.dbPool.Exec(ctx, sqlInsert, value.Id, value.PostCode, value.Date, value.Waterlevel)
		return err
	}

	return nil
}

func (r *HydrologyStatsRepository) GetStartInterval(ctx context.Context, postCode string) (time.Time, error) {
	sql := `
        SELECT MIN(date) FROM waterlevels WHERE post_code = $1;
    `

	var startDate time.Time
	err := r.dbPool.QueryRow(ctx, sql, postCode).Scan(&startDate)
	if err != nil {
		return time.Time{}, err
	}

	return startDate, nil
}

func (r *HydrologyStatsRepository) GetWaterlevelsByDateInterval(ctx context.Context, postCode string, dateStart, dateEnd time.Time) ([]model.Waterlevel, error) {
	var waterlevels []model.Waterlevel

	query := `
        WITH latest_before_start AS (
            SELECT DISTINCT ON (post_code) id, post_code, date, waterlevel
            FROM waterlevels
            WHERE post_code = $1 AND date < $2
            ORDER BY post_code, date DESC
        )
        SELECT id, post_code, date, waterlevel
        FROM waterlevels
        WHERE post_code = $1 AND date BETWEEN $2 AND $3
        UNION ALL
        SELECT * FROM latest_before_start
        ORDER BY date
    `

	rows, err := r.dbPool.Query(ctx, query, postCode, dateStart, dateEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wl model.Waterlevel
		if err := rows.Scan(&wl.Id, &wl.PostCode, &wl.Date, &wl.Waterlevel); err != nil {
			return nil, err
		}
		waterlevels = append(waterlevels, wl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return waterlevels, nil
}
