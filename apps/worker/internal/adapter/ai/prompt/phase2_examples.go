package prompt

import "github.com/kubrickcode/specvital/apps/worker/internal/domain/specview"

// Phase2Example represents a single conversion example for few-shot learning.
type Phase2Example struct {
	Input  string
	Output string
}

// phase2Examples maps languages to their conversion examples.
// Few-shot examples significantly influence LLM output language.
var phase2Examples = map[specview.Language][]Phase2Example{
	"Korean": {
		{"should_login_with_valid_credentials", "유효한 자격 증명으로 로그인 성공"},
		{"returns_404_when_not_found", "존재하지 않으면 404 반환"},
	},
	"Japanese": {
		{"should_login_with_valid_credentials", "有効な資格情報でログイン成功"},
		{"returns_404_when_not_found", "存在しない場合は404を返す"},
	},
	"English": {
		{"should_login_with_valid_credentials", "Successfully logs in with valid credentials"},
		{"returns_404_when_not_found", "Returns 404 when not found"},
	},
	"Chinese": {
		{"should_login_with_valid_credentials", "使用有效凭证成功登录"},
		{"returns_404_when_not_found", "不存在时返回404"},
	},
	"Spanish": {
		{"should_login_with_valid_credentials", "Inicio de sesión exitoso con credenciales válidas"},
		{"returns_404_when_not_found", "Devuelve 404 cuando no se encuentra"},
	},
}

// GetPhase2Examples returns examples for the given language.
// Returns nil if no examples exist for the language.
func GetPhase2Examples(lang specview.Language) []Phase2Example {
	return phase2Examples[lang]
}
