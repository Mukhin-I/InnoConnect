package entity

// Request payload for creating a help request
type CreateRequestRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	RequesterAddress string `json:"requester_address"`
	Type             string `json:"type"`
	Deadline         string `json:"deadline"`
}

// Full request details response after creation
type CreateRequestResponse struct {
	ID               int64  `json:"id"`
	Creator          User   `json:"creator"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	RequesterAddress string `json:"requester_address"`
	Type             string `json:"type"`
	Deadline         string `json:"deadline"`
}

// Request preview for lists
type RequestShort struct {
	RequestID int64  `json:"request_id"`
	Title     string `json:"title"`
	CreatorID int64  `json:"creator_id"`
	Type      string `json:"type"`
	Deadline  string `json:"deadline"`
}

// Complete request with creator details
type RequestFull struct {
	RequestID        int64  `json:"request_id"`
	Creator          User   `json:"creator"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	RequesterAddress string `json:"requester_address"`
	Type             string `json:"type"`
	Deadline         string `json:"deadline"`
}

// Response wrapper for requests list
type GetRequestsResponse struct {
	Requests []RequestShort `json:"requests"`
}
