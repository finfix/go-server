// seed-icons загружает иконки из директории в таблицу coin.icons.
//
// Использование:
//
//	go run ./cmd/seed-icons -dir ./icons
//
// Для каждого файла *.png|*.svg|*.jpg|*.jpeg|*.webp в директории делается
// UPSERT в coin.icons по полю name (имя берётся из имени файла без расширения).
// Параметры подключения берутся из env: PGSQL_HOST, PGSQL_DATABASE,
// PGSQL_USER, PGSQL_PASSWORD.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dir := flag.String("dir", "./icons", "директория с файлами иконок")
	flag.Parse()

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		os.Getenv("PGSQL_USER"),
		os.Getenv("PGSQL_PASSWORD"),
		os.Getenv("PGSQL_HOST"),
		os.Getenv("PGSQL_DATABASE"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}

	entries, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatalf("read dir %s: %v", *dir, err)
	}

	allowed := map[string]bool{".png": true, ".svg": true, ".jpg": true, ".jpeg": true, ".webp": true}

	var loaded int
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if !allowed[ext] {
			continue
		}
		path := filepath.Join(*dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read %s: %v", path, err)
		}
		name := strings.TrimSuffix(e.Name(), ext)

		// UPSERT по name: если иконка с таким именем уже есть — обновляем img,
		// иначе вставляем новую с автогенерированным uuid.
		_, err = db.ExecContext(ctx, `
			INSERT INTO coin.icons (id, name, img)
			VALUES (gen_random_uuid(), $1, $2)
			ON CONFLICT (name) DO UPDATE SET img = EXCLUDED.img
		`, name, data)
		if err != nil {
			log.Fatalf("upsert %s: %v", name, err)
		}
		fmt.Printf("loaded %s (%d bytes)\n", name, len(data))
		loaded++
	}

	fmt.Printf("done, %d icons loaded\n", loaded)
}