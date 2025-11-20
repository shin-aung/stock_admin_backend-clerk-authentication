package main

import (
    "context"
    "time"
)

type Admin struct {
    AdminID     int64     `json:"admin_id"`
    ClerkUserID string    `json:"clerk_user_id"`
    Email       string    `json:"email"`
    FullName    string    `json:"full_name"`
    Role        string    `json:"role"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func GetAdminByClerkID(ctx context.Context, clerkUserID string) (*Admin, error) {
    row := pool.QueryRow(ctx, `
        SELECT admin_id, clerk_user_id, email, full_name, role, is_active, created_at, updated_at
        FROM admins WHERE clerk_user_id=$1
    `, clerkUserID)

    var a Admin
    err := row.Scan(&a.AdminID, &a.ClerkUserID, &a.Email, &a.FullName, &a.Role, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return &a, nil
}

func CreateOrUpdateAdmin(ctx context.Context, a *Admin) error {
    ct := time.Now().UTC()
    _, err := pool.Exec(ctx, `
        INSERT INTO admins (clerk_user_id, email, full_name, role, is_active, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$6)
        ON CONFLICT (clerk_user_id) DO UPDATE
        SET email = EXCLUDED.email,
            full_name = EXCLUDED.full_name,
            role = EXCLUDED.role,
            is_active = EXCLUDED.is_active,
            updated_at = EXCLUDED.updated_at
    `, a.ClerkUserID, a.Email, a.FullName, a.Role, a.IsActive, ct)
    return err
}

func ListAdmins(ctx context.Context) ([]Admin, error) {
    rows, err := pool.Query(ctx, `
        SELECT admin_id, clerk_user_id, email, full_name, role, is_active, created_at, updated_at 
        FROM admins ORDER BY admin_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []Admin
    for rows.Next() {
        var a Admin
        if err := rows.Scan(&a.AdminID, &a.ClerkUserID, &a.Email, &a.FullName, &a.Role, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
            return nil, err
        }
        res = append(res, a)
    }
    return res, nil
}
