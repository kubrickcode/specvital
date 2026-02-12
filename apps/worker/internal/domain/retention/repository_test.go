package retention

import "testing"

func TestDeleteResult_HasMore(t *testing.T) {
	tests := []struct {
		name         string
		deletedCount int64
		batchSize    int
		want         bool
	}{
		{"deleted less than batch", 50, 100, false},
		{"deleted exactly batch", 100, 100, true},
		{"deleted more than batch", 150, 100, true},
		{"zero deleted", 0, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DeleteResult{DeletedCount: tt.deletedCount}
			if got := r.HasMore(tt.batchSize); got != tt.want {
				t.Errorf("HasMore(%d) = %v, want %v", tt.batchSize, got, tt.want)
			}
		})
	}
}
