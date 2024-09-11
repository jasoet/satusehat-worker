package entity

type MappingStatus string

const (
	Ready      MappingStatus = "READY"      // ready to push to SatuSehat
	Incomplete MappingStatus = "INCOMPLETE" // need to re-fetch the detail, if data exists but incomplete
	Invalid    MappingStatus = "INVALID"    // mark the data invalid, no need to continue the process
)

type PublishStatus string

const (
	Success        PublishStatus = "SUCCESS"         // Successfully Publish
	PayloadInvalid PublishStatus = "PAYLOAD_INVALID" // Successfully Publish but PayloadError
	RequestError   PublishStatus = "ERROR"           // Unsuccessfully Publish can be retried
	Preparing      PublishStatus = "PREPARING"
)
