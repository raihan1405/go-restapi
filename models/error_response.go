package models

// ErrorResponse defines the structure for error responses
// @Description Struktur respons kesalahan yang dikembalikan oleh API
type ErrorResponse struct {
    // Message provides a brief description of the error
    // @Description Pesan singkat yang mendeskripsikan kesalahan
    Message string `json:"message" example:"Validation error"`

    // Error provides detailed information about the error
    // @Description Informasi lebih rinci mengenai kesalahan
    Error string `json:"error" example:"Invalid Operator ID"`
}
