package structs

// Don't worry about structs right now
// It will be used to convert the JSON responses into native Go structs

type Initial struct {
	Balance   int `json:"balance"`
	CheckWait int `json:"check_wait"`
	Cost      int `json:"cost"`
	ID        int `json:"id"`
}

type CheckStatus struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type Transcript struct {
	Job      `json:"job"`
	Speakers []Speaker `json:"speakers"`
	Words    []Word    `json:"words"`
}

type Job struct {
	Lang      string `json:"lang"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`
	ID        int    `json:"id"`
}

type Speaker struct {
	Duration   string `json:"duration"`
	Confidence string `json:"confidence"`
	Name       string `json:"name"`
	Time       string `json:"time"`
}

type Word struct {
	Duration   string `json:"duration"`
	Confidence string `json:"confidence"`
	Name       string `json:"name"`
	Time       string `json:"time"`
}
