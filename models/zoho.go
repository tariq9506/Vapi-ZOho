package models

import (
	"fmt"
	"log"
	"strings"
	"time"
	"tutree-vapi/config"
)

type ZohoLeads struct {
	District string
	Phone    string
}
type ZohoCRMLeads struct {
	ZohoLeadId      string
	CallEndedReason string
}

type Lead struct {
	ID          string     `json:"id"`                     // Unique identifier for the lead
	FullName    string     `json:"full_name"`              // Full name of the lead
	LeadStatus  string     `json:"lead_status,omitempty"`  // Status of the lead
	Phone       string     `json:"phone,omitempty"`        // Phone number
	CreatedAt   *time.Time `json:"created_at,omitempty"`   // Creation date and time
	District    string     `json:"district,omitempty"`     // District name
	Gender      string     `json:"gender,omitempty"`       // Gender (male, female, other)
	QueryParam  string     `json:"query_param,omitempty"`  // Query parameter
	SessionType string     `json:"session_type,omitempty"` // Session type
	DialingCode string     `json:"dialing_code,omitempty"`
}

func GetDestrictByPhone(phone string) (string, string, error) {
	district := ""
	zohoLeadID := ""
	db, err := config.GetDB2()
	if err != nil {
		log.Println("[ERROR] GetDestrictByPhone: Failed while connecting with the database :", err)
		return district, zohoLeadID, err
	}
	defer db.Close()
	query := `SELECT district, id
			FROM leads
			WHERE phone = $1`
	err = db.QueryRow(query, phone).Scan(&district, &zohoLeadID)
	if err != nil {
		log.Println("[ERROR] GetDestrictByPhone: Failed to get district from the database with error :", err)
		return district, zohoLeadID, err
	}
	return district, zohoLeadID, nil
}
func GetZoheColdLeads() ([]Lead, error) {
	leads := []Lead{}
	db, err := config.GetDB2()
	if err != nil {
		log.Println("[ERROR] GetZoheColdLeads: Failed while connecting with the database :", err)
		return leads, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, phone, district, dialing_code FROM leads WHERE call_ended_reasone IS NULL LIMIT 1000")
	if err != nil {
		log.Println("[ERROR] GetZoheColdLeads: Failed to get cold leads from the database with error :", err)
		return leads, err
	}
	for rows.Next() {
		lead := Lead{}
		err = rows.Scan(&lead.ID, &lead.Phone, &lead.District, &lead.DialingCode)
		if err != nil {
			log.Println("[ERROR] GetZoheColdLeads: Failed to scan cold leads from the database with error :", err)
			return leads, err
		}
		leads = append(leads, lead)
	}
	return leads, nil
}
func InsertAllLeads(leads []Lead) error {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("[ERROR] InsertAllLeads: Failed while connecting with the database :", err)
		return err
	}
	defer db.Close()
	var value []interface{}
	var placeHolder []string
	// Dynamically construct placeholders and append values
	for i, lead := range leads {

		// Create placeholders for each set of values
		placeHolder = append(placeHolder, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d,'%s')",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5, "+1"))

		// Append the actual values
		value = append(value, lead.ID, lead.FullName, lead.Phone, lead.District, lead.Gender)
	}

	// Construct the final query by joining the placeholders
	query := fmt.Sprintf(`INSERT INTO leads (id, full_name, phone, district, gender,dialing_code) VALUES  %s
						ON CONFLICT(id)
						DO UPDATE SET
						full_name = EXCLUDED.full_name,
						phone = EXCLUDED.phone,
						district = EXCLUDED.district,
						gender = EXCLUDED.gender,
						dialing_code = EXCLUDED.dialing_code
								`, strings.Join(placeHolder, ", "))

	// Execute the query
	_, err = db.Query(query, value...)
	if err != nil {
		log.Println("[ERROR] InsertAllLeads : Failed to execute query with error :", err)
		return err
	}
	return nil
}
func UpdateLeadsInfoAfterCall(zohoLeadData ZohoCRMLeads) error {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("[ERROR] UpdateLeadsInfoAfterCall: Failed while connecting with the database :", err)
		return err
	}
	defer db.Close()
	query := `UPDATE leads SET call_ended_reasone = $1 WHERE id = $2`
	// Execute the query
	_, err = db.Query(query, zohoLeadData.CallEndedReason, zohoLeadData.ZohoLeadId)
	if err != nil {
		log.Println("[ERROR] UpdateLeadsInfoAfterCall : Failed to execute query with error :", err)
		return err
	}
	return nil
}
