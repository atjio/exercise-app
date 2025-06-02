package global

import "os"

var (
	PORT string = ""
	// LOAD_BALANCER_URL is the address of the load balancer.
	// It defaults to "http://127.0.0.1:8080" if the LOADBALANCER_URL environment variable is not set.
	LOAD_BALANCER_URL string
	REGISTER_PATH   string = "/register"
	HEALTHCHECK_PATH string = "/healthcheck"
	REGISTER_URL    string
)

func init() {
	LOAD_BALANCER_URL = os.Getenv("LOADBALANCER_URL")
	if LOAD_BALANCER_URL == "" {
		LOAD_BALANCER_URL = "http://127.0.0.1:8080"
	}
	REGISTER_URL = LOAD_BALANCER_URL + REGISTER_PATH
}

func GetHealthcheckUrl() string {
	// This healthcheck URL is for the service itself, not the load balancer
	return "http://localhost" + PORT + HEALTHCHECK_PATH
}