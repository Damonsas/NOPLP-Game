package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type IndexEntry struct {
	Titre   string `json:"titre"`
	Artiste string `json:"artiste"`
	Ligne   string `json:"ligne"`
}

func UpdateIndex() error {
	dataPath := "data/paroledata"
	indexPath := "data/serverdata/paroledata/index.json"

	files, err := os.ReadDir(dataPath)
	if err != nil {
		return fmt.Errorf("lecture du dossier impossible: %v", err)
	}

	var index []IndexEntry

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dataPath, file.Name()))
		if err != nil {
			fmt.Printf("Impossible de lire %s: %v\n", file.Name(), err)
			continue
		}

		var p IndexEntry
		if err := json.Unmarshal(content, &p); err != nil {
			fmt.Printf("Impossible de parser %s: %v\n", file.Name(), err)
			continue
		}

		p.Ligne = fmt.Sprintf("%s - %s", p.Artiste, p.Titre)
		index = append(index, p)
	}

	os.MkdirAll(filepath.Dir(indexPath), os.ModePerm)

	indexContent, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors du marshalling: %v", err)
	}

	if err := os.WriteFile(indexPath, indexContent, 0644); err != nil {
		return fmt.Errorf("erreur lors de l'écriture de index.json: %v", err)
	}

	fmt.Println("index.json mis à jour !")
	return nil
}

func WatchIndex() {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("Erreur watcher:", err)
			return
		}
		defer watcher.Close()

		if err := watcher.Add("data/paroledata"); err != nil {
			fmt.Println("Erreur watcher add:", err)
			return
		}

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("Changement détecté, mise à jour index.json...")
					UpdateIndex()
				}
			case err := <-watcher.Errors:
				fmt.Println("Erreur watcher:", err)
			}
		}
	}()
}
