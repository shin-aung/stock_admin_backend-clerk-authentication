package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "strings"

    clerk "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/jwks"
    "github.com/clerk/clerk-sdk-go/v2/jwt"
)

type ctxKey string

const ClerkUserIDKey ctxKey = "clerkUserID"

func extractBearerToken(r *http.Request) (string, error) {
    auth := r.Header.Get("Authorization")
    if auth == "" {
        return "", fmt.Errorf("missing authorization header")
    }
    parts := strings.SplitN(auth, " ", 2)
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", fmt.Errorf("invalid authorization header format")
    }
    return parts[1], nil
}

func verifyToken(ctx context.Context, token string, jwksClient *jwks.Client) (string, error) {
    unsafe, err := jwt.Decode(ctx, &jwt.DecodeParams{Token: token})
    if err != nil {
        return "", fmt.Errorf("failed to decode token: %w", err)
    }

    key, err := jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
        KeyID:      unsafe.KeyID,
        JWKSClient: jwksClient,
    })
    if err != nil {
        return "", fmt.Errorf("failed to get JWK: %w", err)
    }

    claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
        Token: token,
        JWK:   key,
    })
    if err != nil {
        return "", fmt.Errorf("token verification failed: %w", err)
    }

    if claims == nil || claims.Subject == "" {
        return "", fmt.Errorf("invalid token claims")
    }

    return claims.Subject, nil
}

func ClerkAuthMiddleware(next http.Handler) http.Handler {
    secret := getEnv("CLERK_SECRET", "sk_test_dh1gBNQXBQ11CdVVXEd83PGhOxBCFQZEKBHFyvz7EG")
    secret := os.GetEnv("CLERK_SECRET")
    if secret == "" {
        panic("CLERK_SECRET environment variable not set")
    }

    clerk.SetKey("Bearer " + secret)
    jwksClient := jwks.NewClient(&clerk.ClientConfig{})

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, err := extractBearerToken(r)
if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
}

userID, err := verifyToken(r.Context(), token, jwksClient)
if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
}

ctx := context.WithValue(r.Context(), ClerkUserIDKey, userID)
next.ServeHTTP(w, r.WithContext(ctx))
    })
}
