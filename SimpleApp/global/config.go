package global

var (
	PORT string = ""
	LOAD_BALANCER_URL string = "http://127.0.0.1:8080"
	REGISTER_PATH string = "/register"
	HEALTHCHECK_PATH string = "/healthcheck"
	REGISTER_URL string = LOAD_BALANCER_URL + REGISTER_PATH
)

func GetHealthcheckUrl() string {
	return "http://localhost" + PORT + HEALTHCHECK_PATH
}