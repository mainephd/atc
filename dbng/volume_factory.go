package dbng

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/nu7hatch/gouuid"
)

//go:generate counterfeiter . VolumeFactory

type VolumeFactory interface {
	GetTeamVolumes(teamID int) ([]CreatedVolume, error)

	CreateContainerVolume(int, Worker, CreatingContainer, string) (CreatingVolume, error)
	FindContainerVolume(int, Worker, CreatingContainer, string) (CreatingVolume, CreatedVolume, error)

	FindBaseResourceTypeVolume(int, *UsedWorkerBaseResourceType) (CreatingVolume, CreatedVolume, error)
	CreateBaseResourceTypeVolume(int, *UsedWorkerBaseResourceType) (CreatingVolume, error)

	FindResourceCacheVolume(Worker, *UsedResourceCache) (CreatingVolume, CreatedVolume, error)
	FindResourceCacheInitializedVolume(Worker, *UsedResourceCache) (CreatedVolume, bool, error)
	CreateResourceCacheVolume(Worker, *UsedResourceCache) (CreatingVolume, error)

	FindVolumesForContainer(CreatedContainer) ([]CreatedVolume, error)
	GetOrphanedVolumes() ([]CreatedVolume, []DestroyingVolume, error)
	GetDuplicateResourceCacheVolumes() ([]CreatingVolume, []CreatedVolume, []DestroyingVolume, error)

	FindCreatedVolume(handle string) (CreatedVolume, bool, error)
}

type volumeFactory struct {
	conn Conn
}

func NewVolumeFactory(conn Conn) VolumeFactory {
	return &volumeFactory{
		conn: conn,
	}
}

func (factory *volumeFactory) GetTeamVolumes(teamID int) ([]CreatedVolume, error) {
	query, args, err := psql.Select(volumeColumns...).
		From("volumes v").
		LeftJoin("workers w ON v.worker_name = w.name").
		LeftJoin("containers c ON v.container_id = c.id").
		LeftJoin("volumes pv ON v.parent_id = pv.id").
		LeftJoin("worker_resource_caches wrc ON wrc.id = v.worker_resource_cache_id").
		Where(sq.Or{
			sq.Eq{
				"v.team_id": teamID,
			},
			sq.Eq{
				"v.team_id": nil,
			},
		}).
		Where(sq.Eq{
			"v.state": "created",
		}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := factory.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdVolumes := []CreatedVolume{}

	for rows.Next() {
		_, createdVolume, _, err := scanVolume(rows, factory.conn)
		if err != nil {
			return nil, err
		}

		createdVolumes = append(createdVolumes, createdVolume)
	}

	return createdVolumes, nil
}

func (factory *volumeFactory) CreateResourceCacheVolume(worker Worker, resourceCache *UsedResourceCache) (CreatingVolume, error) {
	var workerResourcCache *UsedWorkerResourceCache
	err := safeFindOrCreate(factory.conn, func(tx Tx) error {
		var err error
		workerResourcCache, err = WorkerResourceCache{
			WorkerName:    worker.Name(),
			ResourceCache: resourceCache,
		}.FindOrCreate(tx)
		return err
	})
	if err != nil {
		return nil, err
	}

	volume, err := factory.createVolume(
		0,
		worker,
		map[string]interface{}{"worker_resource_cache_id": workerResourcCache.ID},
		VolumeTypeResource,
	)

	volume.resourceCacheID = resourceCache.ID
	return volume, nil
}

func (factory *volumeFactory) CreateBaseResourceTypeVolume(teamID int, uwbrt *UsedWorkerBaseResourceType) (CreatingVolume, error) {
	volume, err := factory.createVolume(
		teamID,
		uwbrt.Worker,
		map[string]interface{}{
			"worker_base_resource_type_id": uwbrt.ID,
			"initialized":                  true,
		},
		VolumeTypeResourceType,
	)
	if err != nil {
		return nil, err
	}

	volume.workerBaseResourceTypeID = uwbrt.ID
	return volume, nil
}

func (factory *volumeFactory) CreateContainerVolume(teamID int, worker Worker, container CreatingContainer, mountPath string) (CreatingVolume, error) {
	volume, err := factory.createVolume(
		teamID,
		worker,
		map[string]interface{}{
			"container_id": container.ID(),
			"path":         mountPath,
			"initialized":  true,
		},
		VolumeTypeContainer,
	)
	if err != nil {
		return nil, err
	}

	volume.path = mountPath
	volume.containerHandle = container.Handle()
	volume.teamID = teamID
	return volume, nil
}

func (factory *volumeFactory) FindVolumesForContainer(container CreatedContainer) ([]CreatedVolume, error) {
	query, args, err := psql.Select(volumeColumns...).
		From("volumes v").
		LeftJoin("workers w ON v.worker_name = w.name").
		LeftJoin("containers c ON v.container_id = c.id").
		LeftJoin("volumes pv ON v.parent_id = pv.id").
		LeftJoin("worker_resource_caches wrc ON wrc.id = v.worker_resource_cache_id").
		Where(sq.Eq{
			"v.state":        VolumeStateCreated,
			"v.container_id": container.ID(),
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := factory.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdVolumes := []CreatedVolume{}

	for rows.Next() {
		_, createdVolume, _, err := scanVolume(rows, factory.conn)
		if err != nil {
			return nil, err
		}

		createdVolumes = append(createdVolumes, createdVolume)
	}

	return createdVolumes, nil
}

func (factory *volumeFactory) FindContainerVolume(teamID int, worker Worker, container CreatingContainer, mountPath string) (CreatingVolume, CreatedVolume, error) {
	return factory.findVolume(teamID, worker, map[string]interface{}{
		"v.container_id": container.ID(),
		"v.path":         mountPath,
	})
}

func (factory *volumeFactory) FindBaseResourceTypeVolume(teamID int, uwbrt *UsedWorkerBaseResourceType) (CreatingVolume, CreatedVolume, error) {
	return factory.findVolume(teamID, uwbrt.Worker, map[string]interface{}{
		"v.worker_base_resource_type_id": uwbrt.ID,
	})
}

func (factory *volumeFactory) FindResourceCacheVolume(worker Worker, resourceCache *UsedResourceCache) (CreatingVolume, CreatedVolume, error) {
	workerResourceCache, found, err := WorkerResourceCache{
		WorkerName:    worker.Name(),
		ResourceCache: resourceCache,
	}.Find(factory.conn)
	if err != nil {
		return nil, nil, err
	}

	if !found {
		return nil, nil, nil
	}

	return factory.findVolume(0, worker, map[string]interface{}{
		"v.worker_resource_cache_id": workerResourceCache.ID,
	})
}

func (factory *volumeFactory) FindResourceCacheInitializedVolume(worker Worker, resourceCache *UsedResourceCache) (CreatedVolume, bool, error) {
	workerResourceCache, found, err := WorkerResourceCache{
		WorkerName:    worker.Name(),
		ResourceCache: resourceCache,
	}.Find(factory.conn)
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	_, createdVolume, err := factory.findVolume(0, worker, map[string]interface{}{
		"v.worker_resource_cache_id": workerResourceCache.ID,
		"v.initialized":              true,
	})
	if err != nil {
		return nil, false, err
	}

	if createdVolume == nil {
		return nil, false, nil
	}

	return createdVolume, true, nil
}

func (factory *volumeFactory) FindCreatedVolume(handle string) (CreatedVolume, bool, error) {
	_, createdVolume, err := factory.findVolume(0, nil, map[string]interface{}{
		"v.handle": handle,
	})
	if err != nil {
		return nil, false, err
	}

	if createdVolume == nil {
		return nil, false, nil
	}

	return createdVolume, true, nil
}

func (factory *volumeFactory) GetOrphanedVolumes() ([]CreatedVolume, []DestroyingVolume, error) {
	query, args, err := psql.Select(volumeColumns...).
		From("volumes v").
		LeftJoin("workers w ON v.worker_name = w.name").
		LeftJoin("containers c ON v.container_id = c.id").
		LeftJoin("volumes pv ON v.parent_id = pv.id").
		LeftJoin("worker_resource_caches wrc ON wrc.id = v.worker_resource_cache_id").
		Where(sq.Eq{
			"v.initialized":                  true,
			"v.worker_resource_cache_id":     nil,
			"v.worker_base_resource_type_id": nil,
			"v.container_id":                 nil,
		}).
		Where(sq.Or{
			sq.Eq{"w.state": string(WorkerStateRunning)},
			sq.Eq{"w.state": string(WorkerStateLanding)},
			sq.Eq{"w.state": string(WorkerStateRetiring)},
		}).
		ToSql()
	if err != nil {
		return nil, nil, err
	}

	rows, err := factory.conn.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	createdVolumes := []CreatedVolume{}
	destroyingVolumes := []DestroyingVolume{}

	for rows.Next() {
		_, createdVolume, destroyingVolume, err := scanVolume(rows, factory.conn)
		if err != nil {
			return nil, nil, err
		}

		if createdVolume != nil {
			createdVolumes = append(createdVolumes, createdVolume)
		}

		if destroyingVolume != nil {
			destroyingVolumes = append(destroyingVolumes, destroyingVolume)
		}
	}

	return createdVolumes, destroyingVolumes, nil
}

func (factory *volumeFactory) GetDuplicateResourceCacheVolumes() ([]CreatingVolume, []CreatedVolume, []DestroyingVolume, error) {
	query, args, err := psql.Select(volumeColumns...).
		From("volumes v").
		LeftJoin("workers w ON v.worker_name = w.name").
		LeftJoin("containers c ON v.container_id = c.id").
		LeftJoin("volumes pv ON v.parent_id = pv.id").
		LeftJoin("volumes dv ON v.worker_resource_cache_id = dv.worker_resource_cache_id").
		LeftJoin("worker_resource_caches wrc ON wrc.id = v.worker_resource_cache_id").
		Where(sq.Eq{
			"v.initialized":  false,
			"dv.initialized": true,
		}).
		Where(sq.Or{
			sq.Eq{"w.state": string(WorkerStateRunning)},
			sq.Eq{"w.state": string(WorkerStateLanding)},
			sq.Eq{"w.state": string(WorkerStateRetiring)},
		}).
		ToSql()
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := factory.conn.Query(query, args...)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	creatingVolumes := []CreatingVolume{}
	createdVolumes := []CreatedVolume{}
	destroyingVolumes := []DestroyingVolume{}

	for rows.Next() {
		creatingVolume, createdVolume, destroyingVolume, err := scanVolume(rows, factory.conn)
		if err != nil {
			return nil, nil, nil, err
		}

		if creatingVolume != nil {
			creatingVolumes = append(creatingVolumes, creatingVolume)
		}

		if createdVolume != nil {
			createdVolumes = append(createdVolumes, createdVolume)
		}

		if destroyingVolume != nil {
			destroyingVolumes = append(destroyingVolumes, destroyingVolume)
		}
	}

	return creatingVolumes, createdVolumes, destroyingVolumes, nil
}

// 1. open tx
// 2. lookup cache id
//   * if not found, create.
//     * if fails (unique violation; concurrent create), goto 1.
// 3. insert into volumes in 'initializing' state
//   * if fails (fkey violation; preexisting cache id was removed), goto 1.
// 4. commit tx

var ErrWorkerResourceTypeNotFound = errors.New("worker resource type no longer exists (stale?)")

// 1. open tx
// 2. lookup worker resource type id
//   * if not found, fail; worker must have new version or no longer supports type
// 3. insert into volumes in 'initializing' state
//   * if fails (fkey violation; worker type gone), fail for same reason as 2.
// 4. commit tx
func (factory *volumeFactory) createVolume(
	teamID int,
	worker Worker,
	columns map[string]interface{},
	volumeType VolumeType,
) (*creatingVolume, error) {
	var volumeID int
	handle, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	columnNames := []string{"worker_name", "handle"}
	columnValues := []interface{}{worker.Name(), handle.String()}
	for name, value := range columns {
		columnNames = append(columnNames, name)
		columnValues = append(columnValues, value)
	}

	if teamID != 0 {
		columnNames = append(columnNames, "team_id")
		columnValues = append(columnValues, teamID)
	}

	err = psql.Insert("volumes").
		Columns(columnNames...). // hey, replace this with SetMap plz
		Values(columnValues...).
		Suffix("RETURNING id").
		RunWith(factory.conn).
		QueryRow().
		Scan(&volumeID)
	if err != nil {
		return nil, err
	}

	return &creatingVolume{
		worker: worker,

		id:     volumeID,
		handle: handle.String(),
		typ:    volumeType,
		teamID: teamID,

		conn: factory.conn,
	}, nil
}

func (factory *volumeFactory) findVolume(teamID int, worker Worker, columns map[string]interface{}) (CreatingVolume, CreatedVolume, error) {
	whereClause := sq.Eq{}
	if teamID != 0 {
		whereClause["v.team_id"] = teamID
	}
	if worker != nil {
		whereClause["v.worker_name"] = worker.Name()
	}

	for name, value := range columns {
		whereClause[name] = value
	}

	row := psql.Select(volumeColumns...).
		From("volumes v").
		LeftJoin("workers w ON v.worker_name = w.name").
		LeftJoin("containers c ON v.container_id = c.id").
		LeftJoin("volumes pv ON v.parent_id = pv.id").
		LeftJoin("worker_resource_caches wrc ON wrc.id = v.worker_resource_cache_id").
		Where(whereClause).
		RunWith(factory.conn).
		QueryRow()
	creatingVolume, createdVolume, _, err := scanVolume(row, factory.conn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return creatingVolume, createdVolume, nil
}

var volumeColumns = []string{
	"v.id",
	"v.handle",
	"v.state",
	"w.name",
	"w.addr",
	"w.baggageclaim_url",
	"v.path",
	"c.handle",
	"pv.handle",
	"v.team_id",
	"wrc.resource_cache_id",
	"v.worker_base_resource_type_id",
	`case when v.container_id is not NULL then 'container'
	  when v.worker_resource_cache_id is not NULL then 'resource'
		when v.worker_base_resource_type_id is not NULL then 'resource-type'
		else 'unknown'
	end`,
}

func scanVolume(row sq.RowScanner, conn Conn) (CreatingVolume, CreatedVolume, DestroyingVolume, error) {
	var id int
	var handle string
	var state string
	var workerName string
	var sqWorkerAddress sql.NullString
	var sqWorkerBaggageclaimURL sql.NullString
	var sqPath sql.NullString
	var sqContainerHandle sql.NullString
	var sqParentHandle sql.NullString
	var sqTeamID sql.NullInt64
	var sqResourceCacheID sql.NullInt64
	var sqWorkerBaseResourceTypeID sql.NullInt64

	var volumeType VolumeType

	err := row.Scan(
		&id,
		&handle,
		&state,
		&workerName,
		&sqWorkerAddress,
		&sqWorkerBaggageclaimURL,
		&sqPath,
		&sqContainerHandle,
		&sqParentHandle,
		&sqTeamID,
		&sqResourceCacheID,
		&sqWorkerBaseResourceTypeID,
		&volumeType,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	var path string
	if sqPath.Valid {
		path = sqPath.String
	}

	var containerHandle string
	if sqContainerHandle.Valid {
		containerHandle = sqContainerHandle.String
	}

	var parentHandle string
	if sqParentHandle.Valid {
		parentHandle = sqParentHandle.String
	}

	var workerBaggageclaimURL string
	if sqWorkerBaggageclaimURL.Valid {
		workerBaggageclaimURL = sqWorkerBaggageclaimURL.String
	}

	var workerAddress string
	if sqWorkerAddress.Valid {
		workerAddress = sqWorkerAddress.String
	}

	var teamID int
	if sqTeamID.Valid {
		teamID = int(sqTeamID.Int64)
	}

	var resourceCacheID int
	if sqResourceCacheID.Valid {
		resourceCacheID = int(sqResourceCacheID.Int64)
	}

	var workerBaseResourceTypeID int
	if sqWorkerBaseResourceTypeID.Valid {
		workerBaseResourceTypeID = int(sqWorkerBaseResourceTypeID.Int64)
	}

	switch state {
	case VolumeStateCreated:
		return nil, &createdVolume{
			id:     id,
			handle: handle,
			typ:    volumeType,
			path:   path,
			teamID: teamID,
			worker: &worker{
				name:            workerName,
				gardenAddr:      &workerAddress,
				baggageclaimURL: &workerBaggageclaimURL,
			},
			containerHandle:          containerHandle,
			parentHandle:             parentHandle,
			resourceCacheID:          resourceCacheID,
			workerBaseResourceTypeID: workerBaseResourceTypeID,
			conn: conn,
		}, nil, nil
	case VolumeStateCreating:
		return &creatingVolume{
			id:     id,
			handle: handle,
			typ:    volumeType,
			path:   path,
			teamID: teamID,
			worker: &worker{
				name:            workerName,
				gardenAddr:      &workerAddress,
				baggageclaimURL: &workerBaggageclaimURL,
			},
			containerHandle:          containerHandle,
			parentHandle:             parentHandle,
			resourceCacheID:          resourceCacheID,
			workerBaseResourceTypeID: workerBaseResourceTypeID,
			conn: conn,
		}, nil, nil, nil
	case VolumeStateDestroying:
		return nil, nil, &destroyingVolume{
			id:     id,
			handle: handle,
			worker: &worker{
				name:            workerName,
				gardenAddr:      &workerAddress,
				baggageclaimURL: &workerBaggageclaimURL,
			},
			conn: conn,
		}, nil
	}

	return nil, nil, nil, nil
}
