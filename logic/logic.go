package logic

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

type ThreatDataService interface {
	CreateThreat(Threat) error
	GetThreat(int) (Threat, error)
}

type ThreatData struct {
	DB *pgx.Conn
}

var DB_INSTANCE ThreatData

func (threatData *ThreatData) CreateThreat(new_threat Threat) error {
	_, err := threatData.DB.Exec(context.Background(), "INSERT INTO threats(filename,sha256,submitted) VALUES ($1,$2,$3)", new_threat.Filename, new_threat.Sha256, new_threat.Submitted)
	if err != nil {
		return err
	}
	return nil
}

func (threatData *ThreatData) GetThreat(threat_id int) (Threat, error) {
	rows, err := threatData.DB.Query(context.Background(), "SELECT threats.id,filename,sha256,array_agg(com),submitted FROM threats INNER JOIN comments com ON com.id = ANY(threats.comments) WHERE threats.id = $1", threat_id)
	if err != nil {
		return Threat{}, err
	}

	threat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Threat])
	if err != nil {
		return Threat{}, err
	}

	return threat, nil
}

func (threat *Threat) Validate() error {
	if threat.Filename == "" || threat.Sha256 == "" {
		return errors.New("cannot submit blank filename or filehash")
	}

	if len(threat.Sha256) != 64 {
		return errors.New("invalid hash, SHA256 hashes are 64 hex chars long")
	}

	return nil
}

func HomepageStats() []Stats {
	results := make([]Stats, 0)
	total_threats := Stats{
		Title: "Total Number of Threats in Database",
		Value: 0,
	}
	recent_threats := Stats{
		Title: "New Threats Last 30 Days",
		Value: 0,
	}
	comments := Stats{
		Title: "Total Comments Left On Threats",
		Value: 0,
	}

	err := DB_INSTANCE.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM threats").Scan(&total_threats.Value)
	if err != nil {
		log.Printf("%v", err)
		total_threats.Value = 0
	}

	err = DB_INSTANCE.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM threats WHERE current_date-submitted::date <= 30").Scan(&recent_threats.Value)
	if err != nil {
		log.Printf("%v", err)
		recent_threats.Value = 0
	}

	err = DB_INSTANCE.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM comments").Scan(&comments.Value)
	if err != nil {
		log.Printf("%v", err)
		comments.Value = 0
	}
	results = append(results, total_threats, recent_threats, comments)

	return results
}
