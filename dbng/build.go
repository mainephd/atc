package dbng

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/event"
)

type BuildStatus string

const (
	BuildStatusPending   BuildStatus = "pending"
	BuildStatusStarted   BuildStatus = "started"
	BuildStatusAborted   BuildStatus = "aborted"
	BuildStatusSucceeded BuildStatus = "succeeded"
	BuildStatusFailed    BuildStatus = "failed"
	BuildStatusErrored   BuildStatus = "errored"
)

type Build interface {
	ID() int
	Interceptible() (bool, error)

	SaveStatus(s BuildStatus) error
	SaveImageResourceVersion(planID atc.PlanID, resourceVersion atc.Version, resourceHash string) error
	SetInterceptible(bool) error

	Finish(s BuildStatus) error
	Delete() (bool, error)
}

type build struct {
	id int

	pipelineID int
	teamID     int
	conn       Conn
}

var ErrBuildDisappeared = errors.New("build-disappeared-from-db")

func (b *build) ID() int { return b.id }

func (b *build) Interceptible() (bool, error) {
	var interceptible bool

	err := psql.Select("interceptible").
		From("builds").
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		QueryRow().Scan(&interceptible)

	if err != nil {
		return true, err
	}

	return interceptible, nil
}

func (b *build) SetInterceptible(i bool) error {
	rows, err := psql.Update("builds").
		Set("interceptible", i).
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrBuildDisappeared
	}

	return nil

}

func (b *build) SaveStatus(s BuildStatus) error {
	rows, err := psql.Update("builds").
		Set("status", string(s)).
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrBuildDisappeared
	}

	return nil
}

func (b *build) Finish(s BuildStatus) error {
	tx, err := b.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var endTime time.Time

	err = tx.QueryRow(`
		UPDATE builds
		SET status = $2, end_time = now(), completed = true
		WHERE id = $1
		RETURNING end_time
	`, b.id, string(s)).Scan(&endTime)
	if err != nil {
		return err
	}

	err = b.saveEvent(tx, event.Status{
		Status: atc.BuildStatus(s),
		Time:   endTime.Unix(),
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`
		DROP SEQUENCE %s
	`, buildEventSeq(b.id)))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (b *build) Delete() (bool, error) {
	rows, err := psql.Delete("builds").
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return false, err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected == 0 {
		return false, ErrBuildDisappeared
	}

	return true, nil
}

func (b *build) SaveImageResourceVersion(planID atc.PlanID, resourceVersion atc.Version, resourceHash string) error {
	version, err := json.Marshal(resourceVersion)
	if err != nil {
		return err
	}

	return safeCreateOrUpdate(
		b.conn,
		func(tx Tx) (sql.Result, error) {
			return psql.Insert("image_resource_versions").
				Columns("version", "build_id", "plan_id", "resource_hash").
				Values(version, b.id, string(planID), resourceHash).
				RunWith(tx).
				Exec()
		},
		func(tx Tx) (sql.Result, error) {
			return psql.Update("image_resource_versions").
				Set("version", version).
				Set("resource_hash", resourceHash).
				Where(sq.Eq{
					"build_id": b.id,
					"plan_id":  string(planID),
				}).
				RunWith(tx).
				Exec()
		},
	)
}

func (b *build) saveEvent(tx Tx, event atc.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	table := fmt.Sprintf("team_build_events_%d", b.teamID)
	if b.pipelineID != 0 {
		table = fmt.Sprintf("pipeline_build_events_%d", b.pipelineID)
	}

	_, err = tx.Exec(fmt.Sprintf(`
		INSERT INTO %s (event_id, build_id, type, version, payload)
		VALUES (nextval('%s'), $1, $2, $3, $4)
	`, table, buildEventSeq(b.id)), b.id, string(event.EventType()), string(event.Version()), payload)
	if err != nil {
		return err
	}

	return nil
}

func createBuildEventSeq(tx Tx, buildid int) error {
	_, err := tx.Exec(fmt.Sprintf(`
		CREATE SEQUENCE %s MINVALUE 0
	`, buildEventSeq(buildid)))
	return err
}

func buildEventSeq(buildid int) string {
	return fmt.Sprintf("build_event_id_seq_%d", buildid)
}
