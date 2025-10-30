// internal/notes/repo_test.go
package notes

import (
	"context"
	"os"
	"testing"

	"example.com/pz8-mongo/internal/db"
)

func TestCreateAndGet(t *testing.T) {
	// Obtener URI desde .env o usar valor por defecto
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://127.0.0.1:27017/?authSource=admin"
	}

	// Nombre único para la base de datos de prueba
	testDBName := "pz8_test_" + t.Name()

	ctx := context.Background()

	// Conectar a MongoDB
	deps, err := db.ConnectMongo(ctx, uri, testDBName)
	if err != nil {
		t.Fatalf("No se pudo conectar a MongoDB: %v", err)
	}

	// Limpiar: desconectar y eliminar la base de datos de prueba
	t.Cleanup(func() {
		_ = deps.Database.Drop(ctx)
		_ = deps.Client.Disconnect(ctx)
	})

	// Crear repositorio
	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("No se pudo crear el repositorio: %v", err)
	}

	// Crear una nota
	title := "Título de prueba"
	content := "Contenido de prueba"
	created, err := repo.Create(ctx, title, content)
	if err != nil {
		t.Fatalf("No se pudo crear la nota: %v", err)
	}

	// Recuperar la nota por ID
	got, err := repo.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatalf("No se pudo obtener la nota por ID: %v", err)
	}

	// Verificar que los datos coincidan
	if got.Title != title {
		t.Errorf("Título esperado %q, obtenido %q", title, got.Title)
	}
	if got.Content != content {
		t.Errorf("Contenido esperado %q, obtenido %q", content, got.Content)
	}
}
