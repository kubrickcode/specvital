package entity

type SortBy string

const (
	SortByName   SortBy = "name"
	SortByRecent SortBy = "recent"
	SortByTests  SortBy = "tests"
)

func ParseSortBy(s string) SortBy {
	switch s {
	case "name":
		return SortByName
	case "tests":
		return SortByTests
	default:
		return SortByRecent
	}
}

func (s SortBy) String() string {
	return string(s)
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func ParseSortOrder(s string) SortOrder {
	switch s {
	case "asc":
		return SortOrderAsc
	default:
		return SortOrderDesc
	}
}

func (s SortOrder) String() string {
	return string(s)
}

func DefaultSortOrder(sortBy SortBy) SortOrder {
	if sortBy == SortByName {
		return SortOrderAsc
	}
	return SortOrderDesc
}

type ViewFilter string

const (
	ViewFilterAll       ViewFilter = "all"
	ViewFilterCommunity ViewFilter = "community"
	ViewFilterMy        ViewFilter = "my"
)

func ParseViewFilter(s string) ViewFilter {
	switch s {
	case "my":
		return ViewFilterMy
	case "community":
		return ViewFilterCommunity
	default:
		return ViewFilterAll
	}
}

func (v ViewFilter) String() string {
	return string(v)
}

type OwnershipFilter string

const (
	OwnershipFilterAll          OwnershipFilter = "all"
	OwnershipFilterMine         OwnershipFilter = "mine"
	OwnershipFilterOrganization OwnershipFilter = "organization"
	OwnershipFilterOthers       OwnershipFilter = "others"
)

func ParseOwnershipFilter(s string) OwnershipFilter {
	switch s {
	case "mine":
		return OwnershipFilterMine
	case "organization":
		return OwnershipFilterOrganization
	case "others":
		return OwnershipFilterOthers
	default:
		return OwnershipFilterAll
	}
}

func (o OwnershipFilter) String() string {
	return string(o)
}
