// seed-icons загружает иконки из директории в таблицу coin.icons.
//
// Использование:
//
//	go run ./cmd/seed-icons -dir ./icons
//
// Ожидается структура: <dir>/<icon-name>/60.png. Для каждой подпапки делается
// UPSERT в coin.icons по полю name (имя = имя подпапки), в img пишется
// содержимое файла 60.png. Параметры подключения берутся из env: PGSQL_HOST,
// PGSQL_DATABASE, PGSQL_USER, PGSQL_PASSWORD.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	var loaded int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		id, ok := iconIDs[name]
		if !ok {
			log.Fatalf("no UUID mapping for icon %q", name)
		}
		path := filepath.Join(*dir, name, "60.png")
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read %s: %v", path, err)
		}

		if err := upsertIcon(ctx, db, id, name, data); err != nil {
			log.Fatalf("upsert %s: %v", name, err)
		}
		fmt.Printf("loaded %s (%d bytes)\n", name, len(data))
		loaded++
	}

	fmt.Printf("done, %d icons loaded\n", loaded)
}

func upsertIcon(ctx context.Context, db *sql.DB, id, name string, data []byte) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO coin.icons (id, name, img) VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, img = EXCLUDED.img`,
		id, name, data)
	return err
}

var iconIDs = map[string]string{
	"apple-pay":        "afacb0b8-9189-4d01-8ef6-b71e8e43d883",
	"baby-bottle":      "4e9df577-7a92-4ee6-b301-a38033a4c638",
	"bank":             "991f5d56-221d-4aec-adc3-9adb7d99e3bc",
	"bank-card":        "17470097-73da-42af-a2b5-a7087e956d59",
	"basket":           "98dafdb6-5a5a-43fe-8346-3afc08a42dc5",
	"beer":             "740241bb-f581-4708-9873-809fda2fb124",
	"bike":             "4c21150a-4f53-4f13-8aff-cc6d7b25a62b",
	"book":             "ac932dcd-752a-4cd8-9ad5-4313b60a3870",
	"bottle":           "ea3f787b-cc7c-4e7e-a28b-f272e4b6f94f",
	"building":         "eeb127d9-9d7e-4b7a-a72f-229a123b2a6d",
	"bus":              "240bf5eb-12f5-4fe7-9679-c538409fa98d",
	"calculator":       "f35b6090-57a3-4e4d-b0bc-b7421be85136",
	"camera":           "9a99b2f4-57e3-4687-9957-6f1d2df77cb8",
	"car":              "6777a17d-e070-47bd-9d6f-0d0c2d2f228e",
	"cards":            "c62f7ef7-bda8-43b7-8c19-ff0fd0fddc94",
	"cash":             "cc0ca597-8573-4809-86d6-1c36a27f13e3",
	"chart":            "11b85386-a66f-4521-a94b-e240b6aba530",
	"credit-card":      "a54007a1-5320-440c-891f-294a6a0fd554",
	"cutlery":          "ba6cf0b1-27d3-4f4e-9f52-5462949e010f",
	"dead":             "ff773744-abe6-4b49-982b-b4ad1f6e3741",
	"dumbbell":         "fa37b6e7-539d-48d6-a9a7-f275ea0130ea",
	"email":            "100cb1ce-58ef-4589-8fb1-18c5d0689472",
	"female-profile":   "ee34321f-00c8-4e3c-a47d-21cc0c75daf4",
	"film":             "7f899034-b2e2-4fd5-8353-55024a81ba21",
	"gamepad":          "caf2da31-181f-4164-b14b-2824142c102c",
	"gas-station":      "d30d26e9-88b0-4b55-b613-8a0e171c001f",
	"gear":             "1906d2cb-ecd4-48bc-8e8b-2043dab63c19",
	"gift":             "2a87d14f-313a-49ad-bb66-730d25a8ddc3",
	"graduation-cap":   "07cc5af1-f600-425e-b5b1-85921546ad15",
	"grape":            "013b2d28-e79a-4fa8-9a91-6beda680c4bc",
	"guitar":           "435dd306-17ce-4052-85cc-49b6ab1196df",
	"hamburger":        "c3f148d0-4c7e-449d-acdc-06670ea80303",
	"hanger":           "ca895268-a8d3-4051-99f9-025f403d2568",
	"house":            "a5693a93-9582-4c74-ba6e-93e4517385a3",
	"iphone":           "acc6c710-49c0-4919-aff8-5c96a2b2ec09",
	"license":          "6d119fbf-d87a-454b-a141-6ba58295ee4f",
	"male-user":        "6d5484eb-f44d-44f1-8029-8a57363d827f",
	"martini-glass":    "260c416a-060b-43cc-a2d0-23f1595dddff",
	"mastercard-logo":  "16a5f0bb-3583-46d5-9758-213dc6fd8104",
	"medicine":         "6e989bbe-acb8-4679-996a-649da733aab6",
	"money":            "a962f67d-b30c-491e-ac42-d3a57916b42d",
	"music":            "e8ba54b9-66e8-4823-8c63-cf5758ca9053",
	"music-record":     "16c18940-f55a-4d89-927d-e76a11290a9a",
	"newspaper":        "7e450a34-91b4-41b7-958b-d309a6d59896",
	"open-end-wrench":  "b0ff2385-5864-4c46-b8f4-9765b9d1c440",
	"painting":         "58bd8d13-6390-4748-acc7-7630575e7631",
	"pet":              "ad2836d1-8e7f-412e-a2ed-f6a73568aeac",
	"pet-bowl":         "7432bf2b-0553-4e59-ada5-a34f146f030c",
	"piggy-bank":       "f55842a7-2576-4f52-9569-b58fc35e95b2",
	"plane":            "efbc4cef-62ec-45f4-ab30-8e221d371896",
	"plug":             "6d330a3c-227f-4f45-a4a2-332ad793118d",
	"prize":            "ea0be2b2-41e2-4be1-84b2-f4ac248768d8",
	"rice-bowl":        "03be84e4-b50f-41b4-beef-8d3b22a58a00",
	"scales":           "4c4e4d3d-e804-4e65-bed1-7cab04d2eb20",
	"shopping-bag":     "5058de78-d439-4da2-b7d7-063d7e46ce61",
	"shopping-cart":    "a10dc6d1-2816-4312-a46e-73718a7b51ea",
	"speaker":          "614c2a37-2cdf-4711-b717-7d2c1a520322",
	"spider-web":       "cedac543-3279-4748-83a3-fa3a73559b92",
	"spray":            "acace456-3be8-49a0-830b-d712423c9ff5",
	"stethoscope":      "238ccf29-f929-475f-a76b-b81b5b2b246f",
	"t-shirt":          "a75a427a-cc24-40d1-bac6-1238c1b74624",
	"telecom":          "3dc47efa-64f1-45c0-9027-43239358cf86",
	"theater":          "89fdf7dc-6bec-4eee-8699-1ca127149757",
	"tv":               "84d5f52d-ae8c-48c6-a909-19b4fba3030c",
	"world":            "2814901c-bf5f-4379-a9a9-0bad7b162795",
}
