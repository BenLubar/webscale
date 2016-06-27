package model_test

import (
	"testing"

	"github.com/BenLubar/webscale/db"
	"github.com/BenLubar/webscale/model"
)

func TestCategories(t *testing.T) {
	t.Parallel()

	ctx := Context(t)
	defer ctx.Tx.Rollback()

	if _, err := ctx.Tx.Exec(db.Prepare(`delete from categories;`)); err != nil {
		t.Fatalf("delete all categories: %+v", err)
	}

	newCategory := db.Prepare(`insert into categories (name, parent_category_id) values ($1::text, $2::bigint) returning id;`)

	var zero, top1, top2, child11, child12, child121 model.CategoryID
	for _, c := range []struct {
		Name   string
		ID     *model.CategoryID
		Parent *model.CategoryID
	}{
		{
			Name:   "Top Level Category 1",
			ID:     &top1,
			Parent: &zero,
		},
		{
			Name:   "Child Category 1-1",
			ID:     &child11,
			Parent: &top1,
		},
		{
			Name:   "Child Category 1-2",
			ID:     &child12,
			Parent: &top1,
		},
		{
			Name:   "Child Category 1-2-1",
			ID:     &child121,
			Parent: &child12,
		},
		{
			Name:   "Top Level Category 2",
			ID:     &top2,
			Parent: &zero,
		},
	} {
		if err := ctx.Tx.QueryRow(newCategory, c.Name, *c.Parent).Scan(c.ID); err != nil {
			t.Fatalf("insert category %q: %+v", c.Name, err)
		}
	}

	categories, err := model.TopLevelCategories(ctx)
	if err != nil {
		t.Fatalf("model.TopLevelCategories: %+v", err)
	}

	if len(categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(categories))
	}

	compareCategory := func(c *model.Category, id model.CategoryID, name string) {
		if c.ID != id || c.Name != name {
			t.Errorf("expected %d:%q, got %d:%q", id, name, c.ID, c.Name)
		}
	}

	compareCategory(categories[0], top1, "Top Level Category 1")
	compareCategory(categories[1], top2, "Top Level Category 2")

	if t.Failed() {
		t.FailNow()
	}

	children, err := categories[0].Children(ctx)
	if err != nil {
		t.Fatalf("model.Category.Children: %+v", err)
	}

	if len(children) != 2 {
		t.Fatalf("expected 2 child categories, got %d", len(children))
	}

	compareCategory(children[0], child11, "Child Category 1-1")
	compareCategory(children[1], child12, "Child Category 1-2")

	if t.Failed() {
		t.FailNow()
	}

	children1, err := children[0].Children(ctx)
	if err != nil {
		t.Fatalf("model.Category.Children (1): %+v", err)
	}

	if len(children1) != 0 {
		t.Fatalf("expected 0 child categories, got %d", len(children1))
	}

	children2, err := children[1].Children(ctx)
	if err != nil {
		t.Fatalf("model.Category.Children (2): %+v", err)
	}

	if len(children2) != 1 {
		t.Fatalf("expected 1 child category, got %d", len(children2))
	}

	compareCategory(children2[0], child121, "Child Category 1-2-1")

	if t.Failed() {
		t.FailNow()
	}

	path, err := children2[0].Path.Get(ctx)
	if err != nil {
		t.Fatalf("model.CategoryIDs.Path: %+v", err)
	}

	if len(path) != 3 {
		t.Fatalf("expected 3 categories in path, got %d", len(path))
	}

	compareCategory(path[0], top1, "Top Level Category 1")
	compareCategory(path[1], child12, "Child Category 1-2")
	compareCategory(path[2], child121, "Child Category 1-2-1")
}
