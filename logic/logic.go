package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThreatDataService interface {
	CreateThreat(Threat) error
	GetThreat(int) (Threat, error)
}

type ThreatData struct {
	DB *pgxpool.Pool
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
	rows, err := threatData.DB.Query(context.Background(), "SELECT threats.id,filename,sha256,coalesce(NULLIF(array_agg(com),'{NULL}'),'{}') as comments,submitted FROM threats LEFT JOIN comments com ON com.id = ANY(threats.comments) WHERE threats.id = $1 GROUP BY threats.id", threat_id)
	if err != nil {
		return Threat{}, err
	}

	threat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Threat])
	if err != nil {
		return Threat{}, err
	}

	return threat, nil
}

func Add_Comment(threat_id int, comment string) error {
	var id int
	row := DB_INSTANCE.DB.QueryRow(context.Background(), "INSERT INTO comments(text,date) VALUES ($1,$2) RETURNING id", comment, time.Now().Format("2006-01-02"))

	err := row.Scan(&id)
	if err != nil {
		return err
	}

	fmt.Printf("UPDATE threats SET comments = comments || %d WHERE id = %d", id, threat_id)
	_, err = DB_INSTANCE.DB.Exec(context.Background(), "UPDATE threats SET comments = comments || $1::int WHERE id = $2", id, threat_id)
	if err != nil {

		return err
	}

	return nil
}

func Get_Comments(threat_id int) []Comment {
	enriched_comments := make([]Comment, 0)
	selected_threat, err := DB_INSTANCE.GetThreat(threat_id)
	if err != nil {
		return enriched_comments
	}

	for _, comment := range selected_threat.Comments {
		split_string := strings.Split(comment, ",")
		split_string[0] = strings.Replace(split_string[0], "(", "", 1)
		id, _ := strconv.Atoi(split_string[0])
		split_string[1] = strings.Replace(split_string[1], `"`, ``, 2)
		split_string[2] = strings.Replace(split_string[2], ")", "", 1)
		new_comment := Comment{
			ID:   id,
			Text: split_string[1],
			Date: split_string[2],
		}
		enriched_comments = append(enriched_comments, new_comment)
	}
	return enriched_comments
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

func Recent_Threats() []Homepage_Threat {
	rows, err := DB_INSTANCE.DB.Query(context.Background(), "select id,filename from threats ORDER BY submitted DESC LIMIT 5")
	if err != nil {
		log.Printf("%v", err)
		return make([]Homepage_Threat, 0)
	}
	threats, err := pgx.CollectRows(rows, pgx.RowToStructByName[Homepage_Threat])

	if err != nil {
		log.Printf("%v", err)
		return make([]Homepage_Threat, 0)
	}

	return threats
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

func Search(search_term string) ([]Search_Result, int) {
	result := make([]Search_Result, 0)
	id := -1

	if len(search_term) == 64 {
		err := DB_INSTANCE.DB.QueryRow(context.Background(), "SELECT id FROM threats WHERE sha256 = $1", search_term).Scan(&id)
		fmt.Printf("%v", err)
		if err == nil {
			return result, id
		}
	}

	rows, err := DB_INSTANCE.DB.Query(context.Background(), "SELECT id,filename,sha256 FROM threats WHERE LOWER(filename) LIKE '%' || LOWER($1) || '%'", search_term)
	if err != nil {
		fmt.Printf("%v", err)
		return result, id
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[Search_Result])
	if err != nil {
		return result, id
	}

	return results, id
}
