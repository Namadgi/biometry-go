# Changelog

## Unreleased (BM-532)

### Added

- **`Consents.Give` / `Consents.Revoke` / `Consents.Status`** record and read the user's authorization and storage consent (`PUT`/`GET /v2/users/{user_id}/consents`), which face and voice enrollment require. `ConsentChange` selects the consents to change and can link the change to a session; selecting neither returns `ErrNoConsentSelected` without a request. All three return `UserConsents`.

### Changed

- `Consents.Approve` / `MyApprovals` are documented as custom consent template calls: they do not unlock enrollment.
- README: the verification flow records consent before enrolling.

## 0.2.0 (BM-338)

Migrated the SDK to the **Biometry API v2** gateway (`/api-gateway/v2`). Endpoint paths below are shortened: `/v2/...` means `/api-gateway/v2/...`. This is a breaking release (the module is pre-1.0, so no `/v2` import-path change).

### Breaking changes

- **Identity model:** biometric calls now take an opaque `UserID` on their input structs, serialized into the API v2 JSON `request` part. The `X-User-Fullname` header and `WithUserFullname`/`WithFullname` options are removed (the gateway strips identity headers on v2).
- **`Client.Video` removed**, replaced by **`Client.Verification`** with `Liveness`, `FaceVerify`, and `VoiceVerify` (the v1 `process-video` split). Each takes a typed input struct.
- **`Client.AVF` → `Client.Deepfake`** (`Submit` / `Status`), targeting `/v2/deepfake/checks`; response type `AVFInferenceResponse`/`AVFStatusResponse` → `DeepfakeCheck`.
- **`Enrollment.Face` / `Enrollment.Voice`** take `EnrollFaceInput` / `EnrollVoiceInput` with `UserID`; voice enrollment is now a single file + phrase (was a list).
- **`FaceMatch.Match`** takes `FaceMatchInput` (`reference_image` + optional `video`, `UseSessionVideo`); `MatchWithSessionVideo` folded into `UseSessionVideo`.
- **`DocAuth.Check`** takes `DocAuthInput`; the `WithInhouseDocAuth`/`WithMRZValidation` options are replaced by `Provider`/`MRZProvider` (`"inhouse" | "idscan"`).
- **`Images`** now manages the user reference image: `Upload(ctx, userID, image, name)` → `PUT /v2/users/{id}/reference-image`; `Get(ctx, userID)` returns raw image bytes.
- **`Sessions.Start`** targets `POST /v2/sessions` and reads `data.session_id`; `Sessions.End` targets `POST /v2/sessions/{id}/end`.
- **`Consents`** approvals target `/v2/consents/*`.
- **Response envelope** is now `{ data, decision, meta }`: `Response.Score`/`ScoringResult`/`DecisionReasons` replaced by `Response.Decision` (`Status`/`Score`/`Reasons`); `AutoEnrolled`/`RemovedServices` removed.
- **Errors:** `APIError` now carries `Code` from the v2 error envelope; the message parser handles v2 (`{error:{code,message}}`), field-error arrays, and the legacy shape.

### Unchanged

- `Transactions` (`/api-transactions`) and `APIKeys` (`/api-projects`) are separate services, untouched by the gateway v2 migration.
- Request-metadata options (`WithDeviceInfo`, `WithGeoLocation`, `WithRequestTags`, `WithClientApp`, `WithClientVersion`) still apply as headers.
