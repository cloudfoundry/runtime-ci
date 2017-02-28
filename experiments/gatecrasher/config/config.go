package config

import (
	"os"
	"strconv"
)

type Config struct {
	Target                   string
	Poll_interval_in_seconds int
	Total_number_of_requests int
}

func Load() Config {
	var config Config
	config = Config{
		Target: "example.com",
		Poll_interval_in_seconds: 1,
		Total_number_of_requests: 10,
	}

	if targetString, ok := os.LookupEnv("TARGET"); ok {
		config.Target = targetString
	}

	if value, ok := os.LookupEnv("POLL_INTERVAL_IN_SECONDS"); ok {
		if numValue, err := strconv.Atoi(value); err != nil {
			panic(err)
		} else {
			config.Poll_interval_in_seconds = numValue
		}
	}

	//	if value, ok := os.LookupEnv("REPORT_INTERVAL_IN_SECONDS"); ok {
	//		if numValue, err := strconv.Atoi(value); err != nil {
	//			panic(err)
	//		} else {
	//			config.Report_interval_in_requests = numValue
	//		}
	//	}

	if value, ok := os.LookupEnv("TOTAL_NUMBER_OF_REQUESTS"); ok {
		if numValue, err := strconv.Atoi(value); err != nil {
			panic(err)
		} else {
			config.Total_number_of_requests = numValue
		}
	}

	return config
}
