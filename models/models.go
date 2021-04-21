package modles

const (
	TIME_FORMAT_1 string = "2006-01-02T15:04:00+07:00"
)

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Data struct {
	ReceivedBy string    `json:"receivedBy"`
	Histories  []History `json:"histories"`
}

type History struct {
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	Formatted   Format `json:"formatted"`
}

type Format struct {
	CreatedAt string `json:"receivedBy"`
}
