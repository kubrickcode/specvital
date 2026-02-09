package entity

type GenerationMode string

const (
	GenerationModeInitial          GenerationMode = "initial"
	GenerationModeRegenerateCached GenerationMode = "regenerate_cached"
	GenerationModeRegenerateFresh  GenerationMode = "regenerate_fresh"
)

func (m GenerationMode) IsRegeneration() bool {
	return m == GenerationModeRegenerateCached || m == GenerationModeRegenerateFresh
}
