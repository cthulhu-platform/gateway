package pkg

import (
	"github.com/cthulhu-platform/common/pkg/env"
)

const (
	BODY_LIMIT_MB = 500
)

var (
	FILE_FOLDER = env.GetEnv("FILE_FOLDER", "./app/fileDump")
	PORT        = env.GetEnv("PORT", "4000")
	CORS_ORIGIN = env.GetEnv("CORS_ORIGIN", "http://localhost:3000")

	// AMQP config
	AMPQ_USER  = env.GetEnv("AMQP_USER", "guest")
	AMPQ_PASS  = env.GetEnv("AMQP_PASS", "guest")
	AMPQ_HOST  = env.GetEnv("AMQP_HOST", "localhost")
	AMPQ_PORT  = env.GetEnv("AMQP_PORT", "5672")
	AMPQ_VHOST = env.GetEnv("AMQP_VHOST", "/")

	// Local config
	LOCAL_AUTH_REPO    = env.GetEnv("LOCAL_REPO_LOC", "./db/auth.db")
	LOCAL_FILE_REPO    = env.GetEnv("LOCAL_FILE_REPO", "./db/file.db")
	LOCAL_GATEWAY_REPO = env.GetEnv("LOCAL_GATEWAY_REPO", "./db/gateway.db")

	// SECRETS
	GITHUB_CLIENT_ID     = env.GetEnv("GITHUB_CLIENT_ID", "")
	GITHUB_CLIENT_SECRET = env.GetEnv("GITHUB_CLIENT_SECRET", "")
	GITHUB_REDIRECT_URI  = env.GetEnv("GITHUB_REDIRECT_URI", "http://localhost:7777/auth/oauth/github/callback")
	
	// JWT Config
	JWT_SECRET       = env.GetEnv("JWT_SECRET", "")
	JWT_ACCESS_EXPIRY  = env.GetEnv("JWT_ACCESS_EXPIRY", "15m")
	JWT_REFRESH_EXPIRY = env.GetEnv("JWT_REFRESH_EXPIRY", "168h") // 7 days
)
