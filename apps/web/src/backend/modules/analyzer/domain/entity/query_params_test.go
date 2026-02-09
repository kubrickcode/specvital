package entity

import "testing"

func TestParseSortBy(t *testing.T) {
	tests := []struct {
		input string
		want  SortBy
	}{
		{"name", SortByName},
		{"recent", SortByRecent},
		{"tests", SortByTests},
		{"", SortByRecent},
		{"invalid", SortByRecent},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseSortBy(tt.input)
			if got != tt.want {
				t.Errorf("ParseSortBy(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseSortOrder(t *testing.T) {
	tests := []struct {
		input string
		want  SortOrder
	}{
		{"asc", SortOrderAsc},
		{"desc", SortOrderDesc},
		{"", SortOrderDesc},
		{"invalid", SortOrderDesc},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseSortOrder(tt.input)
			if got != tt.want {
				t.Errorf("ParseSortOrder(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseViewFilter(t *testing.T) {
	tests := []struct {
		input string
		want  ViewFilter
	}{
		{"all", ViewFilterAll},
		{"my", ViewFilterMy},
		{"community", ViewFilterCommunity},
		{"", ViewFilterAll},
		{"invalid", ViewFilterAll},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseViewFilter(tt.input)
			if got != tt.want {
				t.Errorf("ParseViewFilter(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDefaultSortOrder(t *testing.T) {
	tests := []struct {
		sortBy SortBy
		want   SortOrder
	}{
		{SortByName, SortOrderAsc},
		{SortByRecent, SortOrderDesc},
		{SortByTests, SortOrderDesc},
	}

	for _, tt := range tests {
		t.Run(string(tt.sortBy), func(t *testing.T) {
			got := DefaultSortOrder(tt.sortBy)
			if got != tt.want {
				t.Errorf("DefaultSortOrder(%v) = %v, want %v", tt.sortBy, got, tt.want)
			}
		})
	}
}

func TestSortBy_String(t *testing.T) {
	tests := []struct {
		sortBy SortBy
		want   string
	}{
		{SortByName, "name"},
		{SortByRecent, "recent"},
		{SortByTests, "tests"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.sortBy.String(); got != tt.want {
				t.Errorf("SortBy.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortOrder_String(t *testing.T) {
	tests := []struct {
		sortOrder SortOrder
		want      string
	}{
		{SortOrderAsc, "asc"},
		{SortOrderDesc, "desc"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.sortOrder.String(); got != tt.want {
				t.Errorf("SortOrder.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestViewFilter_String(t *testing.T) {
	tests := []struct {
		viewFilter ViewFilter
		want       string
	}{
		{ViewFilterAll, "all"},
		{ViewFilterMy, "my"},
		{ViewFilterCommunity, "community"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.viewFilter.String(); got != tt.want {
				t.Errorf("ViewFilter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
