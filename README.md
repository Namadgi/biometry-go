# biometry-go

[![Go Reference](https://pkg.go.dev/badge/github.com/Namadgi/biometry-go.svg)](https://pkg.go.dev/github.com/Namadgi/biometry-go)
[![CI](https://github.com/Namadgi/biometry-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Namadgi/biometry-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Namadgi/biometry-go)](go.mod)

Official Go SDK for the [Biometry](https://biometrysolutions.com) identity verification platform. It calls the Biometry API v2 (`/api-gateway/v2`).

## Installation

```bash
go get github.com/Namadgi/biometry-go
```

Requires the Go version in [go.mod](go.mod).

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
    client := biometry.NewClient("your-api-key")

    sessionID, err := client.Sessions.Start(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Session started:", sessionID)
}
```

Every biometric call takes a `UserID`: your own opaque, stable key for the end user (not their name or email).

## Full Verification Flow

Enroll the user's face, then verify them in a session: liveness on a video, face verification on the same video, and a face match against their document.

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
    client := biometry.NewClient("your-api-key")
    userID := "user-123"

    // 1. Enroll the face. Enrollment needs the user's authorization and
    //    storage consent (see Consent below).
    face, _ := os.Open("face.jpg")
    defer face.Close()
    if _, err := client.Enrollment.Face(ctx, biometry.EnrollFaceInput{
        UserID:   userID,
        Face:     face,
        FileName: "face.jpg",
    }); err != nil {
        log.Fatal(err)
    }

    // 2. Start a session
    sessionID, err := client.Sessions.Start(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // 3. Check the identity document
    doc, _ := os.Open("passport.jpg")
    defer doc.Close()
    if _, err := client.DocAuth.Check(ctx, biometry.DocAuthInput{
        Document:  doc,
        FileName:  "passport.jpg",
        SessionID: sessionID,
    }); err != nil {
        log.Fatal(err)
    }

    // 4. Liveness on a video of the user reading the challenge phrase
    video, _ := os.Open("selfie.mp4")
    defer video.Close()
    if _, err := client.Verification.Liveness(ctx, biometry.LivenessInput{
        UserID:    userID,
        Video:     video,
        FileName:  "selfie.mp4",
        Phrase:    "one two three four five six seven",
        SessionID: sessionID,
    }); err != nil {
        log.Fatal(err)
    }

    // 5. Face verification on the same video
    if _, err := client.Verification.FaceVerify(ctx, biometry.VerifyInput{
        UserID:          userID,
        SessionID:       sessionID,
        UseSessionVideo: true,
    }); err != nil {
        log.Fatal(err)
    }

    // 6. Face match: document photo against the session video
    passport, _ := os.Open("passport.jpg")
    defer passport.Close()
    match, err := client.FaceMatch.Match(ctx, biometry.FaceMatchInput{
        UserID:          userID,
        ReferenceImage:  passport,
        ReferenceName:   "passport.jpg",
        SessionID:       sessionID,
        UseSessionVideo: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    if match.Decision != nil {
        fmt.Println("Face match decision:", match.Decision.Status)
    }

    // 7. End the session
    if err := client.Sessions.End(ctx, sessionID, &biometry.EndSessionRequest{}); err != nil {
        log.Fatal(err)
    }
}
```

In a session, face verification and face match must use the liveness media: the same file, or `UseSessionVideo`.

## Consent

Face and voice enrollment and `Images.Upload` need the user's **authorization** and **storage** consent, recorded with `PUT /api-gateway/v2/users/{user_id}/consents`.

`Consents.Approve` and `Consents.MyApprovals` manage your project's custom consent templates. They do **not** unlock enrollment.

## Services

| Service | Methods | API v2 |
|---------|---------|--------|
| `Sessions` | `Start`, `End` | `/v2/sessions` |
| `Verification` | `Liveness`, `FaceVerify`, `VoiceVerify` | `/v2/liveness`, `/v2/face-verify`, `/v2/voice-verify` |
| `Enrollment` | `Face`, `Voice` | `/v2/enrollments/face`, `/v2/enrollments/voice` |
| `FaceMatch` | `Match` | `/v2/face-match` |
| `DocAuth` | `Check` | `/v2/documents/check` |
| `Deepfake` | `Submit`, `Status` | `/v2/deepfake/checks` |
| `Images` | `Upload`, `Get` | `/v2/users/{user_id}/reference-image` |
| `Consents` | `Approve`, `MyApprovals` | `/v2/consents` (custom templates) |
| `Transactions` | `List`, `Count`, `GetBySession`, `GetSamples`, `UpdateTags`, `Flag` | `/api-transactions` (not part of v2) |
| `APIKeys` | `Mint`, `Revoke`, `Info` | `/api-projects` (not part of v2) |

`LivenessInput.ExcludeServices` skips checks by name; use the `Service*` constants (`ServiceFaceLiveness`, `ServiceActiveSpeaker`, `ServiceVisualSpeech`, ...).

## Client Options

```go
client := biometry.NewClient("key",
    biometry.WithBaseURL("https://custom.endpoint.com"),
    biometry.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
    biometry.WithUserAgent("my-app/1.0"),
)
```

| Option | Description |
|--------|-------------|
| `WithBaseURL` | Override the default API base URL |
| `WithHTTPClient` | Provide a custom `*http.Client` |
| `WithUserAgent` | Override the default `User-Agent` header |

## Per-Request Options

Any method accepts `CallOption` values that add request metadata headers:

```go
resp, err := client.Verification.Liveness(ctx, input,
    biometry.WithDeviceInfo(`{"model":"iPhone 15","os":"iOS 18"}`),
    biometry.WithGeoLocation(`{"lat":-35.28,"lng":149.13,"country":"AU"}`),
)
```

| CallOption | Header | Description |
|------------|--------|-------------|
| `WithDeviceInfo` | `X-Device-Info` | JSON device metadata |
| `WithGeoLocation` | `X-Geo-Location` | JSON lat/lng/country/city |
| `WithRequestTags` | `X-Request-Tags` | Comma-separated request labels |
| `WithClientApp` | `X-Client-App` | Calling application identifier |
| `WithClientVersion` | `X-Client-Version` | Calling application version |

User, session and trigger are fields on the input structs, not headers.

## Responses

Endpoints that return `*Response` carry the API v2 envelope `{ data, decision, meta }`:

```go
resp, err := client.Verification.Liveness(ctx, input)

resp.Decision      // *Decision{Status, Score, Reasons}; nil when no scoring system applies
resp.Decode(&out)  // decode the endpoint-specific data
resp.RequestID()   // request trace ID
```

## Error Handling

API errors are returned as `*biometry.APIError`, with the HTTP status, the v2 error code and the message:

```go
resp, err := client.Verification.Liveness(ctx, input)
if err != nil {
    var apiErr *biometry.APIError
    if errors.As(err, &apiErr) {
        fmt.Println("Status:", apiErr.StatusCode)
        fmt.Println("Code:", apiErr.Code) // e.g. "face_not_detected"
        fmt.Println("Message:", apiErr.Message)
    }
    log.Fatal(err)
}
```

Helpers for common status codes:

```go
if biometry.IsUnauthorized(err) { /* 401 -- bad or expired API key */ }
if biometry.IsForbidden(err)    { /* 403 -- insufficient permissions */ }
if biometry.IsNotFound(err)     { /* 404 -- resource not found */ }
```

## Documentation

Full API reference and integration guides: **[docs.biometrysolutions.com/sdks/go](https://docs.biometrysolutions.com/sdks/go/)**.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
