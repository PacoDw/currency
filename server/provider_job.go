package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PacoDw/currency/repository"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// RunProviderJob is a job that will run during all the life cicle of the server making request
// to the Currency Provider, the configuration is determinated by the server config.
// Note: every N time the job is going to make a request to the Currency Provider but this request
// is limited by a timeout, if the timeout is reached the request will be cancelled.
func (s *Server) RunProviderJob(ctx context.Context) {
	s.logger.Info(fmt.Sprintf("Running Currency Provider each %s", s.currencyRequestInterval))

	ticker := time.NewTicker(s.currencyRequestInterval)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// starts the job by runing infinite loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("context is finishing the Currency Provider")

				return
			case <-ticker.C:
				// getting the latest data from the Currency Provider
				meta, blob, errReq := s.currencyClient.GetLatestExchangeRates(ctx)

				// logging the stats
				s.logger.Info("Request",
					zap.String("url", cast.ToString(meta.URL)),
					zap.String("time_elapsed", cast.ToString(meta.Elapsed)),
					zap.String("status", cast.ToString(meta.Status)),
					zap.String("requested_at", cast.ToString(meta.RequestedAt.Format(time.RFC3339))),
					zap.String("details", cast.ToString(errReq)),
				)

				// creating the stats to be saved in the database
				reqStats := repository.RequestStatus{
					URL:         meta.URL,
					TimeElapsed: meta.Elapsed,
					Status:      meta.Status,
					RequestedAt: meta.RequestedAt,
				}

				// checking the database connection and saving the stats into the database
				requestID, err := s.repo.CheckConn().RequestStatus.Insert(reqStats)
				if err != nil {
					s.logger.Warn(fmt.Sprintf("error trying to insert request status: %s", err))
				}

				// if there is an error or the request to the Currency Provider was failed then
				// cancel the current request and wait for the next one according to the Request
				// Interval
				if reqStats.Status == "failure" || errReq != nil {
					continue
				}

				// so far the previous requst was successed then we need to map the data
				cvals := make([]repository.CurrencyValue, 0)

				rootM := map[string]interface{}{}
				if err := json.Unmarshal(blob, &rootM); err != nil {
					s.logger.Warn(fmt.Sprintf("error trying to unmarshall: %s", err))

					return
				}

				lastUpdated, err := time.Parse("2006-01-02T15:04:05Z", rootM["meta"].(map[string]interface{})["last_updated_at"].(string))
				if err != nil {
					log.Printf("err time.Parse: %#+v\n", err)

					return
				}

				data := rootM["data"].(map[string]interface{})

				for _, val := range data {
					cm := val.(map[string]interface{})

					cvals = append(cvals, repository.CurrencyValue{
						Name:         cast.ToString(cm["code"]),
						RequestID:    requestID,
						Value:        cast.ToFloat64(cm["value"]),
						LastUdatedAt: lastUpdated,
					})
				}

				// Insert the data in baches into the database
				if err := s.repo.CheckConn().CurrencyValue.BulkInsert(cvals); err != nil {
					s.logger.Warn(fmt.Sprintf("error make a bulkinsert: %s", err))
				}
			}
		}
	}()

	<-s.quitCurrency
	log.Println("Currency Provider Job is closed successfully")
}
