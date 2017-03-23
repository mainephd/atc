package db

const containerColumns = "c.worker_name, resource_id, check_type, check_source, build_id, plan_id, stage, handle, b.name as build_name, r.name as resource_name, p.id as pipeline_id, p.name as pipeline_name, j.name as job_name, step_name, type, working_directory, env_variables, attempts, process_user, c.id, resource_type_version, c.team_id"

const containerJoins = `
		LEFT JOIN pipelines p
		  ON p.id = c.pipeline_id
		LEFT JOIN resources r
			ON r.id = c.resource_id
		LEFT JOIN builds b
		  ON b.id = c.build_id
		LEFT JOIN jobs j
		  ON j.id = b.job_id`

// this saves off metadata and other things not yet expressed by dbng (best_if_used_by)
// func (db *SQLDB) PutTheRestOfThisCrapInTheDatabaseButPleaseRemoveMeLater(handle string, metadata ContainerMetadata, maxLifetime time.Duration) error {
// 	tx, err := db.conn.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback()

// 	maxLifetimeValue := "NULL"
// 	if maxLifetime > 0 {
// 		maxLifetimeValue = fmt.Sprintf(`NOW() + '%d second'::INTERVAL`, int(maxLifetime.Seconds()))
// 	}

// 	var attempts sql.NullString
// 	if len(metadata.Attempts) > 0 {
// 		attemptsBlob, err := json.Marshal(metadata.Attempts)
// 		if err != nil {
// 			return err
// 		}
// 		attempts.Valid = true
// 		attempts.String = string(attemptsBlob)
// 	}

// 	var id int
// 	err = tx.QueryRow(`
// 		UPDATE containers SET (
// 			best_if_used_by,
// 			process_user,
// 			attempts,
// 			pipeline_id
// 		) = (
// 			`+maxLifetimeValue+`,
// 			$2,
// 			$3,
// 			$4
// 		)
// 		WHERE handle = $1
// 		RETURNING id`,
// 		handle,
// 		metadata.User,
// 		attempts,
// 		metadata.PipelineID,
// 	).Scan(&id)
// 	if err != nil {
// 		return err
// 	}

// 	return tx.Commit()
// }
