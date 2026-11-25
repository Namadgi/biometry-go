# biometry-go

[![Go Reference](https://pkg.go.dev/badge/github.com/Namadgi/biometry-go.svg)](https://pkg.go.dev/github.com/Namadgi/biometry-go)
[![CI](https://github.com/Namadgi/biometry-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Namadgi/biometry-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Namadgi/biometry-go)](go.mod)

Official Go SDK for the [Biometry](https://biometrysolutions.com) identity verification platform.

## Installation

```bash
go get github.com/Namadgi/biometry-go
```

Requires **Go 1.21** or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    biometry "github.com/Namadgi/biometry-go"
)

func main() {
    client := biometry.NewClient("your-api-key",
        biometry.WithUserFullname("John Doe"),
    )

    sessionID, err := client.Sessions.Start(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Session started:", sessionID)
}
```

## Full Verification Flow

A typical identity verification combines session management, consent, document
authentication, face enrollment, video processing, and face matching:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    biometry "github.com/Namadgi/biometry-go"
)

func main() {
    ctx := context.Background()

    client := biometry.NewClient("your-api-key",
        biometry.WithUserFullname("Jane Smith"),
    )

    // 1. Start session
    sessionID, err := client.Sessions.Start(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Session:", sessionID)

    // 2. Record consent
    err = client.Consents.Approve(ctx, "consent-id", "user-id",
        biometry.WithSessionID(sessionID),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 3. Authenticate document
    docFile, _ := os.Open("passport.jpg")
    defer docFile.Close()

    docResp, err := client.DocAuth.Check(ctx, docFile, "passport.jpg",
        biometry.WithSessionID(sessionID),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("DocAuth score: %.2f\n", docResp.Score)

    // 4. Enroll face from document
    docFace, _ := os.Open("passport.jpg")
    defer docFace.Close()

    _, err = client.Enrollment.Face(ctx, biometry.EnrollFaceInput{
        Face:       docFace,
        FileName:   "passport.jpg",
        IsDocument: true,
    }, biometry.WithSessionID(sessionID))
    if err != nil {
        log.Fatal(err)
    }

    // 5. Process video (voice + liveness)
    videoFile, _ := os.Open("selfie.mp4")
    defer videoFile.Close()

    videoResp, err := client.Video.Process(ctx, biometry.ProcessVideoInput{
        Video:     videoFile,
        VideoName: "selfie.mp4",
        Phrase:    "one two three four",
    }, biometry.WithSessionID(sessionID))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Video score: %.2f\n", videoResp.Score)

    // 6. Match face (document image vs session video)
    matchFace, _ := os.Open("passport.jpg")
    defer matchFace.Close()

    matchResp, err := client.FaceMatch.MatchWithSessionVideo(ctx,
        matchFace, "passport.jpg", sessionID,
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Face match score: %.2f\n", matchResp.Score)

    // 7. End session
    err = client.Sessions.End(ctx, sessionID, &biometry.EndSessionRequest{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Session ended")
}
```

## Services

| Service | Methods | Description |
|---------|---------|-------------|
| `Sessions` | `Start`, `End` | Session lifecycle management |
| `Video` | `Process` | Video processing with liveness and voice verification |
| `Enrollment` | `Face`, `Voice` | Face and voice biometric enrollment |
| `FaceMatch` | `Match`, `MatchWithSessionVideo` | Face comparison between images and videos |
| `DocAuth` | `Check` | Document authenticity verification |
| `Consents` | `Approve`, `MyApprovals` | Custom consent management |
| `Images` | `Upload`, `Get` | Face image storage and retrieval |
| `AVF` | `StartInference`, `CheckStatus` | Anti-video-forgery detection |
| `Transactions` | `List`, `Count`, `GetBySession`, `GetSamples`, `UpdateTags`, `Flag` | Transaction queries and management |
| `APIKeys` | `Mint`, `Revoke`, `Info` | API key lifecycle |

## Client Options

Configure the client at initialization:

```go
client := biometry.NewClient("key",
    biometry.WithBaseURL("https://custom.endpoint.com"),
    biometry.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
    biometry.WithUserFullname("Jane Doe"),
    biometry.WithUserAgent("my-app/1.0"),
)
```

| Option | Description |
|--------|-------------|
| `WithBaseURL` | Override the default API base URL |
| `WithHTTPClient` | Provide a custom `*http.Client` |
| `WithUserFullname` | Set the default `X-User-Fullname` header |
| `WithUserAgent` | Override the default `User-Agent` header |

## Per-Request Options

Any method accepts variadic `CallOption` values to set headers on individual requests:

```go
resp, err := client.Video.Process(ctx, input,
    biometry.WithSessionID(sessionID),
    biometry.WithDeviceInfo(`{"model":"iPhone 15","os":"iOS 18"}`),
    biometry.WithGeoLocation(`{"lat":-35.28,"lng":149.13,"country":"AU"}`),
)
```

| CallOption | Header | Description |
|------------|--------|-------------|
| `WithFullname` | `X-User-Fullname` | Override fullname for one request |
| `WithSessionID` | `X-Session-ID` | Attach request to a session |
| `WithActionTrigger` | `X-Action-Trigger` | Set the action trigger label |
| `WithDeviceInfo` | `X-Device-Info` | JSON device metadata |
| `WithGeoLocation` | `X-Geo-Location` | JSON lat/lng/country/city |
| `WithRequestTags` | `X-Request-Tags` | Comma-separated request labels |
| `WithClientApp` | `X-Client-App` | Calling application identifier |
| `WithClientVersion` | `X-Client-Version` | Calling application version |
| `WithInhouseDocAuth` | `X-Inhouse-Docauth` | Toggle in-house vs IDScan docauth |
| `WithMRZValidation` | `X-Inhouse-MRZ` | Enable MRZ validation for documents |

## Response Headers

Endpoints that return `*Response` expose useful metadata from HTTP response headers:

```go
resp, err := client.Video.Process(ctx, input)

resp.RequestID()       // X-Request-Id  -- unique request trace ID
resp.AutoEnrolled()    // X-Auto-Enroll -- true if auto-enrollment occurred
resp.RemovedServices() // X-Removed-Services -- true if services were dropped (missing consent)
```

## Error Handling

All API errors are returned as `*biometry.APIError` with the HTTP status code and
message from the server:

```go
resp, err := client.Video.Process(ctx, input)
if err != nil {
    var apiErr *biometry.APIError
    if errors.As(err, &apiErr) {
        fmt.Println("Status:", apiErr.StatusCode)
        fmt.Println("Message:", apiErr.Message)
    }
    log.Fatal(err)
}
```

Convenience helpers for common status codes:

```go
if biometry.IsUnauthorized(err) { /* 401 -- bad or expired API key */ }
if biometry.IsForbidden(err)    { /* 403 -- insufficient permissions */ }
if biometry.IsNotFound(err)     { /* 404 -- resource not found */ }
```

## Documentation

Full API reference and integration guides are available at
**[docs.biometrysolutions.com/sdks/go](https://docs.biometrysolutions.com/sdks/go/)**.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
