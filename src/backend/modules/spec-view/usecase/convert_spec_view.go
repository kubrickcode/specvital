package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type ConvertSpecViewInput struct {
	CommitSHA      string
	IsForceRefresh bool
	Language       entity.Language
	Owner          string
	Repo           string
}

type ConvertSpecViewOutput struct {
	Files   []FileConversionResult
	Summary ConversionSummary
}

type FileConversionResult struct {
	FilePath  string
	Framework string
	Suites    []SuiteConversionResult
}

type SuiteConversionResult struct {
	Hierarchy string
	Name      string
	Tests     []TestConversionResult
}

type TestConversionResult struct {
	ConvertedName string
	IsFromCache   bool
	Line          int
	Modifier      string
	OriginalName  string
	Status        string
}

type ConversionSummary struct {
	CachedCount    int
	ConvertedAt    time.Time
	ConvertedCount int
	TotalTests     int
}

type ConvertSpecViewUseCase struct {
	aiProvider       port.AIProvider
	analysisProvider port.AnalysisProvider
	cacheRepo        port.CacheRepository
	rateLimiter      port.RateLimiter
}

func NewConvertSpecViewUseCase(
	aiProvider port.AIProvider,
	analysisProvider port.AnalysisProvider,
	cacheRepo port.CacheRepository,
	rateLimiter port.RateLimiter,
) *ConvertSpecViewUseCase {
	return &ConvertSpecViewUseCase{
		aiProvider:       aiProvider,
		analysisProvider: analysisProvider,
		cacheRepo:        cacheRepo,
		rateLimiter:      rateLimiter,
	}
}

func (uc *ConvertSpecViewUseCase) Execute(ctx context.Context, input ConvertSpecViewInput) (*ConvertSpecViewOutput, error) {
	if !input.Language.IsValid() {
		return nil, domain.ErrInvalidLanguage
	}

	analysisInfo, err := uc.analysisProvider.GetAnalysisForConversion(ctx, input.Owner, input.Repo, input.CommitSHA)
	if err != nil {
		return nil, err
	}

	rateLimitKey := fmt.Sprintf("spec-view:convert:%s", analysisInfo.CodebaseID.String())
	if !uc.rateLimiter.Allow(ctx, rateLimitKey) {
		return nil, domain.ErrRateLimited
	}

	testMetas := collectTestMetas(analysisInfo.CodebaseID, analysisInfo.Files, input.Language)

	var cachedResults map[string]*entity.CacheEntry
	if !input.IsForceRefresh {
		cacheKeys := extractCacheKeys(testMetas)
		var err error
		cachedResults, err = uc.cacheRepo.GetCachedConversions(ctx, cacheKeys, uc.aiProvider.ModelID())
		if err != nil {
			slog.WarnContext(ctx, "failed to get cached conversions", "error", err)
			cachedResults = make(map[string]*entity.CacheEntry)
		} else if cachedResults == nil {
			cachedResults = make(map[string]*entity.CacheEntry)
		}
	} else {
		cachedResults = make(map[string]*entity.CacheEntry)
	}

	uncachedByFile := groupUncachedByFile(testMetas, cachedResults)

	newConversions := make(map[string]string)
	for filePath, fileData := range uncachedByFile {
		if len(fileData.suites) == 0 {
			continue
		}

		convertInput := port.ConvertInput{
			FilePath: filePath,
			Language: input.Language,
			Suites:   fileData.suites,
		}

		result, err := uc.aiProvider.ConvertTestNames(ctx, convertInput)
		if err != nil {
			if errors.Is(err, domain.ErrRateLimited) {
				slog.WarnContext(ctx, "rate limited during conversion, returning partial results",
					"file", filePath, "converted_so_far", len(newConversions))
				break
			}
			return nil, fmt.Errorf("convert file %s: %w", filePath, err)
		}

		for globalIdx, convertedName := range result {
			if meta, ok := fileData.indexToMeta[globalIdx]; ok {
				keyHex := hex.EncodeToString(meta.cacheKeyHash)
				newConversions[keyHex] = convertedName
			}
		}
	}

	if len(newConversions) > 0 {
		entries := buildCacheEntries(testMetas, newConversions, uc.aiProvider.ModelID())
		if err := uc.cacheRepo.UpsertCachedConversions(ctx, entries); err != nil {
			slog.WarnContext(ctx, "failed to save cached conversions", "error", err, "count", len(entries))
		}
	}

	return buildOutput(analysisInfo.Files, testMetas, cachedResults, newConversions), nil
}

type testMeta struct {
	cacheKeyHash   []byte
	codebaseID     uuid.UUID
	filePath       string
	framework      string
	language       entity.Language
	line           int
	modifier       string
	originalName   string
	status         string
	suiteHierarchy string
}

func collectTestMetas(codebaseID uuid.UUID, files []port.FileInfo, language entity.Language) []testMeta {
	var metas []testMeta

	for _, file := range files {
		for _, suite := range file.Suites {
			for _, test := range suite.Tests {
				cacheKey := entity.GenerateCacheKey(codebaseID, file.FilePath, suite.Hierarchy, test.Name, language)
				metas = append(metas, testMeta{
					cacheKeyHash:   cacheKey,
					codebaseID:     codebaseID,
					filePath:       file.FilePath,
					framework:      file.Framework,
					language:       language,
					line:           test.Line,
					modifier:       test.Modifier,
					originalName:   test.Name,
					status:         test.Status,
					suiteHierarchy: suite.Hierarchy,
				})
			}
		}
	}

	return metas
}

func extractCacheKeys(metas []testMeta) [][]byte {
	keys := make([][]byte, len(metas))
	for i, m := range metas {
		keys[i] = m.cacheKeyHash
	}
	return keys
}

type uncachedFileData struct {
	indexToMeta map[string]testMeta
	suites      []port.SuiteInput
}

func groupUncachedByFile(metas []testMeta, cached map[string]*entity.CacheEntry) map[string]*uncachedFileData {
	result := make(map[string]*uncachedFileData)

	fileSuiteOrder := make(map[string][]string)
	fileSuiteTestMetas := make(map[string]map[string][]testMeta)

	for i, meta := range metas {
		keyHex := hex.EncodeToString(meta.cacheKeyHash)
		if _, ok := cached[keyHex]; ok {
			continue
		}

		if result[meta.filePath] == nil {
			result[meta.filePath] = &uncachedFileData{
				indexToMeta: make(map[string]testMeta),
			}
			fileSuiteTestMetas[meta.filePath] = make(map[string][]testMeta)
		}

		if _, exists := fileSuiteTestMetas[meta.filePath][meta.suiteHierarchy]; !exists {
			fileSuiteOrder[meta.filePath] = append(fileSuiteOrder[meta.filePath], meta.suiteHierarchy)
		}

		fileSuiteTestMetas[meta.filePath][meta.suiteHierarchy] = append(
			fileSuiteTestMetas[meta.filePath][meta.suiteHierarchy], metas[i])
	}

	for filePath, hierarchies := range fileSuiteOrder {
		localIdx := 1
		for _, hierarchy := range hierarchies {
			testMetas := fileSuiteTestMetas[filePath][hierarchy]
			testNames := make([]string, 0, len(testMetas))

			for _, tm := range testMetas {
				result[filePath].indexToMeta[fmt.Sprintf("%d", localIdx)] = tm
				testNames = append(testNames, tm.originalName)
				localIdx++
			}

			result[filePath].suites = append(result[filePath].suites, port.SuiteInput{
				Hierarchy: hierarchy,
				Tests:     testNames,
			})
		}
	}

	return result
}

func buildCacheEntries(metas []testMeta, newConversions map[string]string, modelID string) []*entity.CacheEntry {
	var entries []*entity.CacheEntry

	for _, meta := range metas {
		keyHex := hex.EncodeToString(meta.cacheKeyHash)
		convertedName, ok := newConversions[keyHex]
		if !ok {
			continue
		}

		entries = append(entries, &entity.CacheEntry{
			CacheKeyHash:   meta.cacheKeyHash,
			CodebaseID:     meta.codebaseID,
			ConvertedName:  convertedName,
			FilePath:       meta.filePath,
			Framework:      meta.framework,
			Language:       meta.language,
			ModelID:        modelID,
			OriginalName:   meta.originalName,
			SuiteHierarchy: meta.suiteHierarchy,
		})
	}

	return entries
}

func buildOutput(files []port.FileInfo, metas []testMeta, cached map[string]*entity.CacheEntry, newConversions map[string]string) *ConvertSpecViewOutput {
	var fileResults []FileConversionResult
	cachedCount := 0
	convertedCount := 0
	totalTests := 0
	metaIdx := 0

	for _, file := range files {
		var suiteResults []SuiteConversionResult

		for _, suite := range file.Suites {
			var testResults []TestConversionResult

			for _, test := range suite.Tests {
				if metaIdx >= len(metas) {
					break
				}
				meta := metas[metaIdx]
				metaIdx++
				totalTests++

				keyHex := hex.EncodeToString(meta.cacheKeyHash)

				var convertedName string
				var isFromCache bool

				if entry, ok := cached[keyHex]; ok {
					convertedName = entry.ConvertedName
					isFromCache = true
					cachedCount++
				} else if name, ok := newConversions[keyHex]; ok {
					convertedName = name
					isFromCache = false
					convertedCount++
				} else {
					convertedName = test.Name
					isFromCache = false
				}

				testResults = append(testResults, TestConversionResult{
					ConvertedName: convertedName,
					IsFromCache:   isFromCache,
					Line:          test.Line,
					Modifier:      test.Modifier,
					OriginalName:  test.Name,
					Status:        test.Status,
				})
			}

			suiteResults = append(suiteResults, SuiteConversionResult{
				Hierarchy: suite.Hierarchy,
				Name:      suite.Name,
				Tests:     testResults,
			})
		}

		fileResults = append(fileResults, FileConversionResult{
			FilePath:  file.FilePath,
			Framework: file.Framework,
			Suites:    suiteResults,
		})
	}

	return &ConvertSpecViewOutput{
		Files: fileResults,
		Summary: ConversionSummary{
			CachedCount:    cachedCount,
			ConvertedAt:    time.Now().UTC(),
			ConvertedCount: convertedCount,
			TotalTests:     totalTests,
		},
	}
}
