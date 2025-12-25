package auth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/GoAdminGroup/go-admin/modules/db"
)

// HashPassword はパスワードをbcryptでハッシュ化する
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash はパスワードとハッシュを比較する
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// UpdateAdminPassword はadminユーザーのパスワードを更新する
func UpdateAdminPassword(conn db.Connection, username, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = conn.Exec(
		"UPDATE goadmin_users SET password = ? WHERE username = ?",
		hashedPassword, username,
	)
	return err
}

// CreateAdminUser は管理者ユーザーを作成する（存在しない場合）
func CreateAdminUser(conn db.Connection, username, password, name string) error {
	// 既存ユーザーのチェック
	result, err := conn.Query("SELECT COUNT(*) as count FROM goadmin_users WHERE username = ?", username)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		count, ok := result[0]["count"]
		if ok {
			switch v := count.(type) {
			case int64:
				if v > 0 {
					return nil // ユーザーが既に存在
				}
			case int:
				if v > 0 {
					return nil
				}
			case float64:
				if v > 0 {
					return nil
				}
			}
		}
	}

	// 新規ユーザーを作成
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = conn.Exec(
		"INSERT INTO goadmin_users (username, password, name, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))",
		username, hashedPassword, name,
	)
	if err != nil {
		return err
	}

	// 管理者ロールを割り当て
	result, err = conn.Query("SELECT id FROM goadmin_users WHERE username = ?", username)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}

	userID, ok := result[0]["id"]
	if !ok {
		return nil
	}

	_, err = conn.Exec(
		"INSERT OR IGNORE INTO goadmin_role_users (role_id, user_id, created_at, updated_at) VALUES (1, ?, datetime('now'), datetime('now'))",
		userID,
	)

	return err
}
