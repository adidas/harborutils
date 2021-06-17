package server

type APIError struct {
	Msg string `json:"msg" example:"error msg`
}

type Token struct {
	Token string `json:"token" example:"123.456.789"`
}

type ArtifactSha struct {
	Image   string `json:"image" example:"ngnix:latest"`
	Project string `json:"project" example:"pea-cicd"`
	Sha     string `json:"sha" example:"sha256:a1c2d5c775a3b7ebc7af29c77241819a86cd1222b1931d0712afdcd69c7dcbd5"`
}
type ArtifactCheckSha struct {
	Image        string `json:"image" example:"ngnix:latest"`
	Project      string `json:"project" example:"pea-cicd"`
	Sha          string `json:"sha" example:"sha256:a1c2d5c775a3b7ebc7af29c77241819a86cd1222b1931d0712afdcd69c7dcbd5"`
	TargetDigset string `json:"targetDigset" example:"sha256:a1c2d5c775a3b7ebc7af29c77241819a86cd1222b1931d0712afdcd69c7dcbd5"`
	Equals       bool   `json:"equals" example:"true"`
}

type HealthStatus struct {
	Status string `json:"status" example:"healthy"`
}
