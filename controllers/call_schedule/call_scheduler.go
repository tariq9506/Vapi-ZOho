package callschedule

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"tutree-vapi/constants"
	"tutree-vapi/models"
	"tutree-vapi/service"
)

type LeadInfo struct {
	District string
	ZohoID   string
}

func ScheduleCalls(leads []models.Lead) error {
	const batchSize = 10
	totalLeads := len(leads)

	for i := 0; i < totalLeads; i += batchSize {
		end := i + batchSize
		if end > totalLeads {
			end = totalLeads
		}

		batch := leads[i:end]
		if err := processBatch(batch); err != nil {
			return err
		}
	}

	return nil
}

func processBatch(batch []models.Lead) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(batch))
	defer close(errChan)

	for _, lead := range batch {
		wg.Add(1)
		go func(lead models.Lead) {
			defer wg.Done()

			if err := processLead(lead); err != nil {
				errChan <- err
			}
		}(lead)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func processLead(lead models.Lead) error {
	phone := lead.DialingCode + lead.Phone
	if !strings.HasPrefix(phone, "+") {
		phone = fmt.Sprintf("+%s", phone)
	}

	statusCode, callID, err := service.MakeACallOfZohoColdLeads(phone, lead.District)
	log.Println("MakeACallOfZohoColdLeads LOG", statusCode, callID, err)
	if err != nil {
		log.Printf("Failed to initiate call to %s: %v\n", phone, err)
		return errors.New("failed to initiate call: " + err.Error())
	}

	if statusCode != http.StatusCreated {
		log.Printf("Call initiation failed for %s with status code: %d\n", phone, statusCode)
		return errors.New("call initiation failed with status code: " + http.StatusText(statusCode))
	}

	log.Printf("Successfully initiated the call to %s with ID: %s\n", phone, callID)
	for {
		endedResponse, err := service.GetCallDetails(callID)
		if err != nil {
			log.Printf("Error checking status for call %s: %v", callID, err)
			continue
		}
		status, ok := endedResponse["status"]
		if !ok {
			log.Printf("Invalid status type for call %s", callID)
			return errors.New("invalid status type in response")
		}
		log.Println("------------------------->", endedResponse["offer_response"])
		if endedResponse["offer_response"] == "Interested" {
		}
		if status == "in-progress" {
			log.Println("Call is in progress", status)
		}
		log.Printf("Call %s status: %s", callID, status)

		if status == "ended" {
			var zohoLeadData models.ZohoCRMLeads
			log.Printf("Call %s ended with reason: %v", callID, endedResponse["endedReason"])

			switch endedResponse["endedReason"] {
			case constants.CustomerEndedCall, constants.AssistantEndedCall:
				zohoLeadData.CallEndedReason = endedResponse["offer_response"]
				if zohoLeadData.CallEndedReason == "" {
					zohoLeadData.CallEndedReason = "Not Interested"
				}

			default:
				zohoLeadData.CallEndedReason = endedResponse["endedReason"]
			}

			zohoLeadData.ZohoLeadId = lead.ID

			// if err := zoho.UpdateLeadsOnZohoCRM(zohoLeadData); err != nil {
			// 	log.Printf("Error updating status for call in Zoho CRM lead dashboard %s: %v", callID, err)
			// 	return errors.New("error updating Zoho CRM: " + err.Error())
			// }

			if err := models.UpdateLeadsInfoAfterCall(zohoLeadData); err != nil {
				log.Printf("Error updating status for call in database %s: %v", callID, err)
				return errors.New("error updating database: " + err.Error())
			}
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}
