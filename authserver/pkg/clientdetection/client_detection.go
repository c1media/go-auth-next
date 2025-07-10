package clientdetection

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ClientType represents the type of client making the request
type ClientType string

const (
	ClientTypeNextJS      ClientType = "nextjs"
	ClientTypeReactNative ClientType = "react-native"
	ClientTypeAPI         ClientType = "api"
	ClientTypeWeb         ClientType = "web"
	ClientTypeUnknown     ClientType = "unknown"
)

// ClientInfo contains information about the client
type ClientInfo struct {
	Type   ClientType
	Origin string
}

// DetectClient determines the client type from the request
func DetectClient(c *gin.Context) ClientInfo {
	// 1. Check explicit client type header (most reliable)
	if clientType := c.GetHeader("X-Client-Type"); clientType != "" {
		return ClientInfo{
			Type:   ClientType(strings.ToLower(clientType)),
			Origin: c.GetHeader("Origin"),
		}
	}

	// 2. Check User-Agent for known patterns
	userAgent := strings.ToLower(c.GetHeader("User-Agent"))

	// Next.js patterns
	if strings.Contains(userAgent, "next.js") || strings.Contains(userAgent, "nextjs") {
		return ClientInfo{
			Type:   ClientTypeNextJS,
			Origin: c.GetHeader("Origin"),
		}
	}

	// React Native patterns
	if strings.Contains(userAgent, "react-native") || strings.Contains(userAgent, "reactnative") {
		return ClientInfo{
			Type:   ClientTypeReactNative,
			Origin: c.GetHeader("Origin"),
		}
	}

	// API client patterns (common HTTP libraries)
	if strings.Contains(userAgent, "curl") ||
		strings.Contains(userAgent, "postman") ||
		strings.Contains(userAgent, "insomnia") ||
		strings.Contains(userAgent, "axios") ||
		strings.Contains(userAgent, "fetch") {
		return ClientInfo{
			Type:   ClientTypeAPI,
			Origin: c.GetHeader("Origin"),
		}
	}

	// 3. Check Referer/Origin for development patterns
	referer := c.GetHeader("Referer")
	origin := c.GetHeader("Origin")

	// Next.js development server patterns
	if strings.Contains(referer, "localhost:3000") ||
		strings.Contains(origin, "localhost:3000") ||
		strings.Contains(referer, "vercel.app") ||
		strings.Contains(origin, "vercel.app") {
		return ClientInfo{
			Type:   ClientTypeNextJS,
			Origin: origin,
		}
	}

	// 4. Check for browser-specific headers
	if c.GetHeader("Sec-Fetch-Mode") != "" || c.GetHeader("Sec-Fetch-Site") != "" {
		return ClientInfo{
			Type:   ClientTypeWeb,
			Origin: origin,
		}
	}

	// 5. Default to unknown
	return ClientInfo{
		Type:   ClientTypeUnknown,
		Origin: origin,
	}
}

// RequiresCSRF determines if the client type requires CSRF protection
func (ci ClientInfo) RequiresCSRF() bool {
	switch ci.Type {
	case ClientTypeWeb:
		return true
	case ClientTypeNextJS, ClientTypeReactNative, ClientTypeAPI:
		return false
	case ClientTypeUnknown:
		// Default to requiring CSRF for unknown clients (safer)
		return true
	default:
		return true
	}
}

// IsMobileClient checks if the client is a mobile app
func (ci ClientInfo) IsMobileClient() bool {
	return ci.Type == ClientTypeReactNative
}

// IsWebClient checks if the client is a web browser
func (ci ClientInfo) IsWebClient() bool {
	return ci.Type == ClientTypeNextJS || ci.Type == ClientTypeWeb
}

// IsAPIClient checks if the client is an API consumer
func (ci ClientInfo) IsAPIClient() bool {
	return ci.Type == ClientTypeAPI
}
