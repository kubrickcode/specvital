package api

import "net/http"

func NewBadRequest(detail string) BadRequestApplicationProblemPlusJSONResponse {
	return BadRequestApplicationProblemPlusJSONResponse{
		Detail: detail,
		Status: http.StatusBadRequest,
		Title:  "Bad Request",
	}
}

func NewNotFound(detail string) NotFoundApplicationProblemPlusJSONResponse {
	return NotFoundApplicationProblemPlusJSONResponse{
		Detail: detail,
		Status: http.StatusNotFound,
		Title:  "Not Found",
	}
}

func NewInternalError(detail string) InternalErrorApplicationProblemPlusJSONResponse {
	return InternalErrorApplicationProblemPlusJSONResponse{
		Detail: detail,
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
	}
}

func NewUnauthorized(detail string) UnauthorizedApplicationProblemPlusJSONResponse {
	return UnauthorizedApplicationProblemPlusJSONResponse{
		Detail: detail,
		Status: http.StatusUnauthorized,
		Title:  "Unauthorized",
	}
}

func NewTooManyRequests(detail string) TooManyRequestsApplicationProblemPlusJSONResponse {
	return TooManyRequestsApplicationProblemPlusJSONResponse{
		Detail: detail,
		Status: http.StatusTooManyRequests,
		Title:  "Too Many Requests",
	}
}

func NewForbidden(detail string) ForbiddenApplicationProblemPlusJSONResponse {
	return ForbiddenApplicationProblemPlusJSONResponse{
		Detail: detail,
		Status: http.StatusForbidden,
		Title:  "Forbidden",
	}
}
