package main

import (
    "context"
    "encoding/json"
    "net/http"
)

func requireClerkUserID(r *http.Request) (string, bool) {
    v := r.Context().Value(ClerkUserIDKey)
    if v == nil {
        return "", false
    }
    id, ok := v.(string)
    return id, ok
}

func listAdminsHandler(w http.ResponseWriter, r *http.Request) {
    _, ok := requireClerkUserID(r)
    if !ok {
        http.Error(w, "unauthenticated", http.StatusUnauthorized)
        return
    }

    admins, err := ListAdmins(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(admins)
}

func upsertAdminHandler(w http.ResponseWriter, r *http.Request) {
    clerkUserID, ok := requireClerkUserID(r)
    if !ok {
        http.Error(w, "unauthenticated", http.StatusUnauthorized)
        return
    }

    var a Admin
    if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }

    a.ClerkUserID = clerkUserID

    if err := CreateOrUpdateAdmin(context.Background(), &a); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
