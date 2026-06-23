package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/inertia"
	"github.com/imohamedsheta/xvalid"
)

// Handler represents the base handler for all handlers
type Handler struct {
	Inertia *inertia.Inertia
	flash   *inertia.InmemFlashProvider
}

func NewBaseHandler(i *inertia.Inertia) *Handler {
	return &Handler{
		Inertia: i,
		flash:   x.Flash(),
	}
}

// BindFilters binds query params to a struct and fills map[string]string fields
// using either `form` or `json` tag
func (h *Handler) BindFilters(c *gin.Context, filters any) error {
	// Bind normal query params first
	if err := c.ShouldBindQuery(filters); err != nil {
		return fmt.Errorf("failed to bind filters: %w", err)
	}

	// Use reflection to fill map[string]string fields from QueryMap
	v := reflect.ValueOf(filters)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Only handle map[string]string
		if field.Kind() != reflect.Map || field.Type().Key().Kind() != reflect.String || field.Type().Elem().Kind() != reflect.String {
			continue
		}

		// Get tag name: prefer `form`, fallback to `json`
		tag := fieldType.Tag.Get("form")
		if tag == "" {
			tag = fieldType.Tag.Get("json")
		}
		if tag == "" {
			continue
		}

		// Populate map from query
		queryMap := c.QueryMap(tag)
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}
		for k, val := range queryMap {
			field.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(val))
		}
	}

	return nil
}

// CustomMultipartBind handles complex multipart form binding with nested structures
func (h *Handler) CustomMultipartBind(c *gin.Context, req xvalid.Validatable) error {
	// Parse the multipart form with a reasonable max memory
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	form := c.Request.MultipartForm
	if form == nil {
		return fmt.Errorf("multipart form is nil")
	}

	// Create a map to store all form data for processing
	formData := make(map[string]any)

	// Process regular form fields
	for key, values := range form.Value {
		if len(values) == 1 {
			formData[key] = values[0]
		} else if len(values) > 1 {
			formData[key] = values
		}
	}

	// Process file fields
	for key, files := range form.File {
		if len(files) == 1 {
			formData[key] = files[0]
		} else if len(files) > 1 {
			formData[key] = files
		}
	}

	// // Set the raw form data
	// req.SetRequestSentFields(formData)

	// Now bind the data to the struct
	return h.bindFormDataToStruct(formData, req)
}

// bindFormDataToStruct uses reflection to bind form data to struct
func (h *Handler) bindFormDataToStruct(formData map[string]any, req xvalid.Validatable) error {
	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() == reflect.Pointer {
		reqValue = reqValue.Elem()
	}
	reqType := reqValue.Type()

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		fieldType := reqType.Field(i)

		if !field.CanSet() {
			continue
		}

		// Get the form tag
		formTag := fieldType.Tag.Get("form")
		if formTag == "" || formTag == "-" {
			continue
		}

		// Handle different field types
		switch field.Kind() {
		case reflect.String:
			if value, ok := formData[formTag].(string); ok {
				field.SetString(value)
			}
		case reflect.Int, reflect.Int64:
			if strValue, ok := formData[formTag].(string); ok {
				if intValue, err := strconv.Atoi(strValue); err == nil {
					field.SetInt(int64(intValue))
				}
			}
		case reflect.Float64:
			if strValue, ok := formData[formTag].(string); ok {
				if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
					field.SetFloat(floatValue)
				}
			}
		case reflect.Slice:
			h.handleSliceField(field, fieldType, formTag, formData)
		case reflect.Pointer:
			// Handle file uploads (multipart.FileHeader)
			if field.Type() == reflect.TypeFor[*multipart.FileHeader]() {
				if file, ok := formData[formTag].(*multipart.FileHeader); ok {
					field.Set(reflect.ValueOf(file))
				}
			}
		}
	}

	return nil
}

// handleSliceField handles slice fields (arrays) in form data
func (h *Handler) handleSliceField(field reflect.Value, fieldType reflect.StructField, formTag string, formData map[string]any) {
	sliceType := field.Type().Elem()

	if sliceType.Kind() == reflect.String {
		// Handle []string (like categories)
		var stringSlice []string

		// Look for array-style form data: category[0], category[1], etc.
		for key, value := range formData {
			if strings.HasPrefix(key, formTag+"[") {
				if strValue, ok := value.(string); ok {
					stringSlice = append(stringSlice, strValue)
				}
			}
		}

		// Also check for direct array values
		if values, ok := formData[formTag].([]string); ok {
			stringSlice = append(stringSlice, values...)
		}

		if len(stringSlice) > 0 {
			field.Set(reflect.ValueOf(stringSlice))
		}
	} else if sliceType.Kind() == reflect.Struct {
		// Handle []Struct (like tickets)
		h.handleStructSlice(field, sliceType, formTag, formData)
	}
}

// handleStructSlice handles slices of structs (like tickets array)
func (h *Handler) handleStructSlice(field reflect.Value, structType reflect.Type, formTag string, formData map[string]any) {
	// Find all indices for this struct slice
	indices := make(map[int]bool)
	prefix := formTag + "["

	for key := range formData {
		if strings.HasPrefix(key, prefix) {
			// Extract index: tickets[0][name] -> 0
			parts := strings.Split(key[len(prefix):], "]")
			if len(parts) > 0 {
				if index, err := strconv.Atoi(parts[0]); err == nil {
					indices[index] = true
				}
			}
		}
	}

	if len(indices) == 0 {
		return
	}

	// Find the maximum index to create slice
	maxIndex := -1
	for index := range indices {
		if index > maxIndex {
			maxIndex = index
		}
	}

	// Create and populate the slice
	slice := reflect.MakeSlice(field.Type(), maxIndex+1, maxIndex+1)

	for i := 0; i <= maxIndex; i++ {
		structValue := slice.Index(i)
		h.populateStructFromForm(structValue, structType, fmt.Sprintf("%s[%d]", formTag, i), formData)
	}

	field.Set(slice)
}

// populateStructFromForm populates a struct from form data
func (h *Handler) populateStructFromForm(structValue reflect.Value, structType reflect.Type, prefix string, formData map[string]any) {
	for i := 0; i < structType.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		if !field.CanSet() {
			continue
		}

		formTag := fieldType.Tag.Get("form")
		if formTag == "" || formTag == "-" {
			continue
		}

		fieldKey := prefix + "[" + formTag + "]"

		switch field.Kind() {
		case reflect.String:
			if value, ok := formData[fieldKey].(string); ok {
				field.SetString(value)
			}
		case reflect.Int, reflect.Int64:
			if strValue, ok := formData[fieldKey].(string); ok {
				if intValue, err := strconv.Atoi(strValue); err == nil {
					field.SetInt(int64(intValue))
				}
			}
		case reflect.Float64:
			if strValue, ok := formData[fieldKey].(string); ok {
				if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
					field.SetFloat(floatValue)
				}
			}
		case reflect.Bool:
			if strValue, ok := formData[fieldKey].(string); ok {
				if boolValue, err := strconv.ParseBool(strValue); err == nil {
					field.SetBool(boolValue)
				}
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				// Handle []string fields like benefits
				var stringSlice []string
				benefitsPrefix := fieldKey + "["

				for key, value := range formData {
					if strings.HasPrefix(key, benefitsPrefix) {
						if strValue, ok := value.(string); ok {
							stringSlice = append(stringSlice, strValue)
						}
					}
				}

				if len(stringSlice) > 0 {
					field.Set(reflect.ValueOf(stringSlice))
				}
			}
		}
	}
}

// BindAndExtract handles JSON binding + raw field extraction
func (h *Handler) BindBodyAndExtractToRequest(c *gin.Context, req xvalid.Validatable) error {
	if err := h.UniversalBind(c, req); err != nil {
		return err
	}

	// if c.ContentType() == "application/json" {
	// 	var raw map[string]any
	// 	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err == nil {
	// 		req.SetRequestSentFields(raw)
	// 	}
	// }

	return nil
}

func (h *Handler) BindAndValidate(c *gin.Context, req xvalid.Validatable) error {
	ct := c.ContentType()

	// Use custom multipart binding for form data with files
	if strings.Contains(ct, "multipart/form-data") {
		if err := h.CustomMultipartBind(c, req); err != nil {
			return fmt.Errorf("multipart binding failed: %w", err)
		}
	} else {
		if err := h.BindBodyAndExtractToRequest(c, req); err != nil {
			return err
		}
	}

	if xe := x.Validator().ValidateRequest(c, req); xe != nil {
		var ve *xvalid.ValidationError
		if errors.As(xe, &ve) {
			if err := h.flash.FlashErrors(c, ve.Fields); err != nil {
				return err
			}

			h.Inertia.Back(c, 422)
			return xe
		}
		return xe
	}

	return nil
}

func (h *Handler) BindAndValidateWithFlashError(c *gin.Context, req xvalid.Validatable, flashMessage string) error {
	ct := c.ContentType()

	// Use custom multipart binding for form data with files
	if strings.Contains(ct, "multipart/form-data") {
		if err := h.CustomMultipartBind(c, req); err != nil {
			return fmt.Errorf("multipart binding failed: %w", err)
		}
	} else {
		if err := h.BindBodyAndExtractToRequest(c, req); err != nil {
			return err
		}
	}

	if xe := x.Validator().ValidateRequest(c, req); xe != nil {
		var ve *xvalid.ValidationError
		if errors.As(xe, &ve) {
			props := inertia.Props{
				"toast": []map[string]any{
					utils.ToastError(flashMessage, flashMessage),
				},
			}
			if err := h.Flash(c, props); err != nil {
				return err
			}

			h.Inertia.Back(c, 422)
			return nil
		}
		return xe
	}

	return nil
}

func (h *Handler) BackWithFlash(c *gin.Context, message string, status int) error {
	props := inertia.Props{
		"drawer": "close",
		"toast": []map[string]any{
			utils.ToastSuccess("نجحت العملية", message),
		},
	}

	if err := h.Flash(c, props); err != nil {
		return err
	}

	h.Inertia.Back(c, status)
	return nil
}

// UniversalBind attempts to bind request based on Content-Type
func (h *Handler) UniversalBind(c *gin.Context, req xvalid.Validatable) error {
	ct := c.ContentType()

	switch {
	case ct == "application/json":
		if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
			return err
		}

		// var raw map[string]any
		// if err := c.ShouldBindBodyWith(&raw, binding.JSON); err == nil {
		// 	req.SetRequestSentFields(raw)
		// }

	case ct == "application/x-www-form-urlencoded":
		if err := c.ShouldBindWith(req, binding.Form); err != nil {
			return err
		}

	case strings.Contains(ct, "multipart/form-data"):
		// Use our custom multipart binding instead of default
		return h.CustomMultipartBind(c, req)

	default:
		if err := c.ShouldBind(req); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) FlashAlertBack(c *gin.Context, status int, toast ...map[string]any) error {
	return h.FlashBack(c, inertia.Props{
		"toast": toast,
	}, status)
}

func (h *Handler) FlashBack(c *gin.Context, props inertia.Props, status int) error {
	if err := h.Flash(c, props); err != nil {
		return err
	}

	h.Inertia.Back(c, status)
	return nil
}

func (h *Handler) Flash(c *gin.Context, props inertia.Props) error {
	err := h.flash.Flash(c, props)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) FlashProvider(c *gin.Context) *inertia.InmemFlashProvider {
	return h.flash
}

func (h *Handler) HandleValidationErrors(c *gin.Context, err error) bool {
	var xe *xerr.XErr
	if errors.As(err, &xe) && xe.IsType(enums.XErrValidationError) {
		utils.Dump(xe.Details)

		_ = h.flash.FlashErrors(c, xe.Details)
		h.Inertia.Back(c, 422)
		return true
	}
	return false
}
