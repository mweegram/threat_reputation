package logic

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mweegram/threat_reputation/database"
)

type ThreatDataService interface {
	CreateThreat(database.Threat) error
	GetThreat(int) (database.Threat, error)
}

type ThreatData struct {
	db *pgx.Conn
}

func (threatData *ThreatData) CreateThreat(new_threat database.Threat) error {
	_, err := threatData.db.Exec(context.Background(), "INSERT INTO threats(filename,sha256,submitted) VALUES ($1,$2,$3)", new_threat.Filename, new_threat.Sha256, new_threat.Submitted)
	if err != nil {
		return err
	}
	return nil
}

func (threatData *ThreatData) GetThreat(threat_id int) (database.Threat, error) {
	rows, err := threatData.db.Query(context.Background(), "SELECT threats.id,filename,sha256,array_agg(com),submitted FROM threats INNER JOIN comments com ON com.id = ANY(threats.comments) WHERE threats.id = $1", threat_id)
	if err != nil {
		return database.Threat{}, err
	}

	threat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[database.Threat])
	if err != nil {
		return database.Threat{}, err
	}

	return threat, nil
}
