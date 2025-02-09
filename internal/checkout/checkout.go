package checkout

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rminder/internal/app"
	"rminder/web"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

func CreatePremiumCheckoutSession(ctx *gin.Context) {
	domain := os.Getenv("STRIPE_RETURN_TO_URL")
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(os.Getenv("STRIPE_PREMIUM_PRICE_ID")),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/checkout/success"),
		CancelURL:  stripe.String(domain + "/"),
	}

	params.PaymentIntentData = &stripe.CheckoutSessionPaymentIntentDataParams{}
	params.PaymentIntentData.AddMetadata("user_id", app.GetUser(ctx).Id)

	s, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}

	ctx.Redirect(http.StatusSeeOther, s.URL)
}

func PremiumCheckoutSuccessHandler(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := web.PremiumPaymentSuccessful().Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in CreatePremiumCheckoutSession: %e :: %v", err, e)
	}
}

func CheckoutWebhookHandler(s *app.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const MaxBodyBytes = int64(65536)
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
		payload, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
			ctx.Writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		event := stripe.Event{}

		if err := json.Unmarshal(payload, &event); err != nil {
			fmt.Fprintf(os.Stderr, "[warning] Webhook error while parsing basic request. %v\n", err.Error())
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Replace this endpoint secret with your endpoint's unique secret
		// If you are testing with the CLI, find the secret by running 'stripe listen'
		// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
		// at https://dashboard.stripe.com/webhooks
		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		signatureHeader := ctx.Request.Header.Get("Stripe-Signature")
		event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[warning] Webhook signature verification failed. %v\n", err)
			ctx.Writer.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
			return
		}
		// Unmarshal the event data into an appropriate struct depending on its Type
		switch event.Type {
		case "payment_intent.succeeded":
			var paymentIntent stripe.PaymentIntent
			err := json.Unmarshal(event.Data.Raw, &paymentIntent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
				ctx.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("Successful payment for %d.", paymentIntent.Amount)
			// Then define and call a func to handle the successful payment intent.
			handlePaymentIntentSucceeded(s, paymentIntent)
		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func handlePaymentIntentSucceeded(s *app.App, paymentIntent stripe.PaymentIntent) {
	user_id, ok := paymentIntent.Metadata["user_id"]
	if !ok {
		log.Printf("No user_id in metadata")
		return
	}

	user, err := s.GetUser(user_id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}

	user.HasPremium = true

	err = s.SaveUser(user)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		return
	}
}
