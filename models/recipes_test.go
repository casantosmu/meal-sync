package models

import (
	"reflect"
	"testing"
)

func TestIngredientToList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic",
			input:    "tomato\nonion\nlettuce",
			expected: []string{"tomato", "onion", "lettuce"},
		},
		{
			name:     "empty",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "single",
			input:    "tomato",
			expected: []string{"tomato"},
		},
		{
			name:     "extra spaces",
			input:    "  tomato \n onion  \n lettuce ",
			expected: []string{"tomato", "onion", "lettuce"},
		},
		{
			name:     "multi lines",
			input:    "tomato\n\nonion\n\nlettuce",
			expected: []string{"tomato", "onion", "lettuce"},
		},
		{
			name:     "duplicated",
			input:    "tomato\nonion\ntomato",
			expected: []string{"tomato", "onion", "tomato"},
		},
		{
			name:     "newline at the end",
			input:    "tomato\nonion\nlettuce\n",
			expected: []string{"tomato", "onion", "lettuce"},
		},
		{
			name:     "special characters",
			input:    "tomato\nonion\nlettuce & spinach",
			expected: []string{"tomato", "onion", "lettuce & spinach"},
		},
		{
			name:     "long names",
			input:    "tomato\nonion\nlong grain brown rice",
			expected: []string{"tomato", "onion", "long grain brown rice"},
		},
		{
			name:     "with numbers",
			input:    "tomato\n2 onions\nlettuce",
			expected: []string{"tomato", "2 onions", "lettuce"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Recipe{Ingredients: tc.input}.IngredientsToList()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Test case: %s - expected %v, but got %v", tc.name, tc.expected, got)
			}
		})
	}
}

func TestDirectionsToList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic",
			input:    "Preheat oven to 350F.\n\nMix ingredients in a bowl.\n\nBake for 20 minutes.",
			expected: []string{"Preheat oven to 350F.", "Mix ingredients in a bowl.", "Bake for 20 minutes."},
		},
		{
			name:     "empty",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "single instruction",
			input:    "Preheat oven to 350F.",
			expected: []string{"Preheat oven to 350F."},
		},
		{
			name:     "extra spaces",
			input:    "  Preheat oven to 350F.  \n\n  Mix ingredients.  \n\n  Bake for 20 minutes.  ",
			expected: []string{"Preheat oven to 350F.", "Mix ingredients.", "Bake for 20 minutes."},
		},
		{
			name:     "multiple newlines",
			input:    "Preheat oven.\n\n\n\nMix ingredients.\n\n\nBake.",
			expected: []string{"Preheat oven.", "Mix ingredients.", "Bake."},
		},
		{
			name:     "newline at the end",
			input:    "Preheat oven to 350F.\n\nMix ingredients.\n\nBake for 20 minutes.\n",
			expected: []string{"Preheat oven to 350F.", "Mix ingredients.", "Bake for 20 minutes."},
		},
		{
			name:     "long directions",
			input:    "Preheat oven to 350F.\n\nIn a large bowl, mix sugar, flour, and eggs until smooth.\n\nBake for 45 minutes or until golden brown.",
			expected: []string{"Preheat oven to 350F.", "In a large bowl, mix sugar, flour, and eggs until smooth.", "Bake for 45 minutes or until golden brown."},
		},
		{
			name:     "special characters",
			input:    "Preheat oven to 350F.\n\nMix sugar, flour, & eggs.\n\nBake at 350°F for 30 minutes.",
			expected: []string{"Preheat oven to 350F.", "Mix sugar, flour, & eggs.", "Bake at 350°F for 30 minutes."},
		},
		{
			name:     "windows newlines",
			input:    "Preheat oven to 350F.\r\n\r\nMix ingredients in a bowl.\r\n\r\nBake for 20 minutes.",
			expected: []string{"Preheat oven to 350F.", "Mix ingredients in a bowl.", "Bake for 20 minutes."},
		},
		{
			name:     "mixed newlines",
			input:    "Preheat oven to 350F.\r\n\nMix ingredients.\n\nBake.",
			expected: []string{"Preheat oven to 350F.", "Mix ingredients.", "Bake."},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Recipe{Directions: tc.input}.DirectionsToList()
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Test case: %s - expected %v, but got %v", tc.name, tc.expected, got)
			}
		})
	}
}
