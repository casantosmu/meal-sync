package models_test

import (
	"reflect"
	"testing"

	"github.com/casantosmu/meal-sync/models"
)

func TestIngredientToList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic",
			input:    "tomato\r\nonion\r\nlettuce",
			expected: []string{"tomato", "onion", "lettuce"},
		},
		{
			name:     "empty",
			input:    "",
			expected: []string{},
		},
		{
			name:     "extra spaces",
			input:    "    tomato\r\nonion   \r\n   lettuce   ",
			expected: []string{"tomato", "onion", "lettuce"},
		},
		{
			name:     "multi lines",
			input:    "tomato\r\n\r\nonion\r\n\r\nlettuce",
			expected: []string{"tomato", "onion", "lettuce"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := models.Recipe{Ingredients: tc.input}.IngredientsToList()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Test case: %s - expected %v, but got %v", tc.name, tc.expected, got)
			}
		})
	}
}

func TestDirectionsToGroups(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []models.DirectionGroup
	}{
		{
			name:  "basic",
			input: "Preheat oven to 350F.\r\nMix ingredients in a bowl.\r\nBake for 20 minutes.",
			expected: []models.DirectionGroup{
				{
					Heading:    "",
					Directions: []string{"Preheat oven to 350F.", "Mix ingredients in a bowl.", "Bake for 20 minutes."},
				},
			},
		},
		{
			name:  "basic with heading",
			input: "Title 1:\r\n\r\nPreheat oven to 350F.\r\nMix ingredients in a bowl.\r\nBake for 20 minutes.",
			expected: []models.DirectionGroup{
				{
					Heading:    "Title 1:",
					Directions: []string{"Preheat oven to 350F.", "Mix ingredients in a bowl.", "Bake for 20 minutes."},
				},
			},
		},
		{
			name:  "multiple direction groups",
			input: "Title 1:\r\n\r\nPreheat oven to 350F.\r\nMix ingredients in a bowl.\r\n\r\nTitle 2:\r\n\r\nBake for 20 minutes.",
			expected: []models.DirectionGroup{
				{
					Heading:    "Title 1:",
					Directions: []string{"Preheat oven to 350F.", "Mix ingredients in a bowl."},
				},
				{
					Heading:    "Title 2:",
					Directions: []string{"Bake for 20 minutes."},
				},
			},
		},
		{
			name:     "empty",
			input:    "",
			expected: []models.DirectionGroup{},
		},
		{
			name:  "directions with extra spaces",
			input: "    Preheat oven to 350F.\r\nMix ingredients in a bowl.   \r\n   Bake for 20 minutes.   ",
			expected: []models.DirectionGroup{
				{
					Heading:    "",
					Directions: []string{"Preheat oven to 350F.", "Mix ingredients in a bowl.", "Bake for 20 minutes."},
				},
			},
		},
		{
			name:  "directions with multiple newlines",
			input: "Preheat oven.\r\n\r\nMix ingredients.\r\n\r\nBake.",
			expected: []models.DirectionGroup{
				{
					Heading:    "",
					Directions: []string{"Preheat oven.", "Mix ingredients.", "Bake."},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := models.Recipe{Directions: tc.input}.DirectionsToGroups()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Test case: %s - expected %v, but got %v", tc.name, tc.expected, got)
			}
		})
	}
}
