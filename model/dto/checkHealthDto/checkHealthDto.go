package checkHealthDto

type (
	VersionRequest struct {
		Version string `json:"version" binding:"required"`
	}
)
