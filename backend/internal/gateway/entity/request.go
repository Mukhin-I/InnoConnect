package entity

type CreateRequestRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	RequesterAddress string `json:"requester_address"`
	Type             string `json:"type"`
	Deadline         string `json:"deadline"`
}

type CreateRequestResponse struct {
	ID               int64 `json:"id"`
	Creator          User  `json:"creator"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	RequesterAddress string `json:"requester_address"`
	Type             string `json:"type"`
	Deadline         string `json:"deadline"`
}

type RequestShort struct {
	RequestID int64  `json:"request_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Deadline  string `json:"deadline"`
}

type RequestFull struct {
	RequestID        int64 `json:"request_id"`
	Creator          User  `json:"creator"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	RequesterAddress string `json:"requester_address"`
	Type             string `json:"type"`
	Deadline         string `json:"deadline"`
}

type GetRequestsResponse struct {
	Requests []RequestShort `json:"requests"`
}
