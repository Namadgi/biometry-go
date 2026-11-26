package biometry_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Namadgi/biometry-go"
)

func Example() {
	client := biometry.NewClient(os.Getenv("BIOMETRY_API_KEY"))

	ctx := context.Background()

	// Start a biometric session
	sessionID, err := client.Sessions.Start(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Session started:", sessionID)

	// Run a liveness check
	video := strings.NewReader("video-bytes") // use a real file in practice
	resp, err := client.Verification.Liveness(ctx, biometry.LivenessInput{
		UserID:    "user-123",
		Video:     video,
		FileName:  "sample.mp4",
		Phrase:    "one two three four five",
		SessionID: sessionID,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if resp.Decision != nil {
		fmt.Println("Decision:", resp.Decision.Status)
	}

	// End the session
	_ = client.Sessions.End(ctx, sessionID, nil)
}

func Example_docauth() {
	client := biometry.NewClient(os.Getenv("BIOMETRY_API_KEY"))
	ctx := context.Background()

	doc := strings.NewReader("document-bytes") // use a real file
	resp, err := client.DocAuth.Check(ctx, biometry.DocAuthInput{
		Document: doc,
		FileName: "passport.jpg",
		Provider: "inhouse",
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var info biometry.DocAuthInfo
	if err := resp.Decode(&info); err != nil {
		fmt.Println("Decode error:", err)
		return
	}
	fmt.Printf("Document: %s, Country: %s\n", info.DocumentCategory, info.CountryCode)
}

func Example_transactions() {
	client := biometry.NewClient(os.Getenv("BIOMETRY_API_KEY"))
	ctx := context.Background()

	txns, err := client.Transactions.List(ctx, biometry.TransactionListRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Found %d transactions\n", len(txns))
}
