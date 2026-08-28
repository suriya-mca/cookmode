package services

import (
	"context"
	"strings"

	sdk "github.com/meilisearch/meilisearch-go"

	"cookmode/internal/models"
)

type SearchIndexer struct {
	client sdk.ServiceManager
}

func NewSearchIndexer(host, apiKey string) *SearchIndexer {
	return &SearchIndexer{client: sdk.New(host, sdk.WithAPIKey(apiKey))}
}

func (s *SearchIndexer) Client() sdk.ServiceManager {
	return s.client
}

type recipeDoc struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Cuisine       string   `json:"cuisine"`
	Difficulty    string   `json:"difficulty"`
	DietaryTags   []string `json:"dietary_tags"`
	Ingredients   []string `json:"ingredients"`
	IngredientRaw []string `json:"ingredient_raw"`
	TotalTime     int      `json:"total_time"`
	PrepTime      int      `json:"prep_time"`
	CookTime      int      `json:"cook_time"`
	Servings      int      `json:"servings"`
}

func recipeToDoc(r *models.Recipe) recipeDoc {
	ings := make([]string, 0, len(r.Ingredients))
	for _, ing := range r.Ingredients {
		ings = append(ings, strings.ToLower(ing.Name))
	}
	tags := r.DietaryTags
	if tags == nil {
		tags = []string{}
	}
	return recipeDoc{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Cuisine:     r.Cuisine,
		Difficulty:  string(r.Difficulty),
		DietaryTags: tags,
		Ingredients: ings,
		TotalTime:   r.PrepTimeMin + r.CookTimeMin,
		PrepTime:    r.PrepTimeMin,
		CookTime:    r.CookTimeMin,
		Servings:    r.Servings,
	}
}

func (s *SearchIndexer) EnsureIndex(ctx context.Context) error {
	_, err := s.client.GetIndexWithContext(ctx, "recipes")
	if err == nil {
		return s.configure(ctx)
	}
	_, err = s.client.CreateIndexWithContext(ctx, &sdk.IndexConfig{
		Uid:        "recipes",
		PrimaryKey: "id",
	})
	if err != nil {
		return err
	}
	return s.configure(ctx)
}

func (s *SearchIndexer) configure(ctx context.Context) error {
	idx := s.client.Index("recipes")
	settings := &sdk.Settings{
		SearchableAttributes: []string{"title", "description", "ingredients", "cuisine"},
		FilterableAttributes: []string{"ingredients", "dietary_tags", "difficulty", "cuisine", "total_time", "prep_time", "cook_time"},
		SortableAttributes:   []string{"total_time"},
		TypoTolerance: &sdk.TypoTolerance{
			Enabled: true,
		},
	}
	_, err := idx.UpdateSettingsWithContext(ctx, settings)
	return err
}

func (s *SearchIndexer) IndexRecipe(ctx context.Context, r *models.Recipe) error {
	if r.Status != models.StatusPublished {
		return s.DeleteRecipe(ctx, r.ID)
	}
	if err := s.EnsureIndex(ctx); err != nil {
		return err
	}
	doc := recipeToDoc(r)
	task, err := s.client.Index("recipes").AddDocumentsWithContext(ctx, []recipeDoc{doc}, nil)
	if err != nil {
		return err
	}
	_, err = s.client.WaitForTaskWithContext(ctx, task.TaskUID, 0)
	return err
}

func (s *SearchIndexer) DeleteRecipe(ctx context.Context, id string) error {
	task, err := s.client.Index("recipes").DeleteDocumentWithContext(ctx, id, nil)
	if err != nil {
		return err
	}
	_, err = s.client.WaitForTaskWithContext(ctx, task.TaskUID, 0)
	return err
}

type SearchParams struct {
	Query       string
	Ingredients []string
	Dietary     []string
	Difficulty  string
	Cuisine     string
	MaxTime     *int
	Limit       int
	Offset      int
}

func (s *SearchIndexer) Search(ctx context.Context, p SearchParams) (*sdk.SearchResponse, error) {
	if err := s.EnsureIndex(ctx); err != nil {
		return nil, err
	}
	var filters []string
	if len(p.Ingredients) > 0 {
		for _, ing := range p.Ingredients {
			filters = append(filters, `ingredients = "`+strings.ToLower(strings.TrimSpace(ing))+`"`)
		}
	}
	if len(p.Dietary) > 0 {
		for _, d := range p.Dietary {
			filters = append(filters, `dietary_tags = "`+strings.TrimSpace(d)+`"`)
		}
	}
	if p.Difficulty != "" {
		filters = append(filters, `difficulty = "`+p.Difficulty+`"`)
	}
	if p.Cuisine != "" {
		filters = append(filters, `cuisine = "`+p.Cuisine+`"`)
	}
	if p.MaxTime != nil {
		filters = append(filters, `total_time <= `+itoa(*p.MaxTime))
	}
	filterStr := ""
	if len(filters) > 0 {
		filterStr = strings.Join(filters, " AND ")
	}
	limit := p.Limit
	if limit == 0 {
		limit = 20
	}
	req := &sdk.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(p.Offset),
		Filter: filterStr,
	}
	return s.client.Index("recipes").SearchWithContext(ctx, p.Query, req)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
