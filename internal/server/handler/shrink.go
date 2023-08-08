package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/config"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/fetch"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"go.uber.org/zap"
)

const (
	// DefaultQuality is the default quality to use when shrinking an image.
	DefaultQuality int = 60

	// DefaultCompression is the default compression level to use when shrinking an image.
	DefaultCompression int = 9

	// DefaultQuantTable is the default quantization table to use when shrinking an image.
	DefaultQuantTable int = 3

	// DefaultOptimizeCoding is the default optimize coding setting to use when shrinking an image.
	DefaultOptimizeCoding bool = true

	// DefaultInterlace is the default interlace setting to use when shrinking an image.
	DefaultInterlace bool = false

	// DefaultStripMetadata is the default strip metadata setting to use when shrinking an image.
	DefaultStripMetadata bool = true

	// DefaultOptimizeICCProfile is the default optimize ICC profile setting to use when shrinking an image.
	DefaultOptimizeICCProfile bool = true

	// DefaultTrellisQuant is the default trellis quant setting to use when shrinking an image.
	DefaultTrellisQuant bool = true

	// DefaultOvershootDeringing is the default overshoot deringing setting to use when shrinking an image.
	DefaultOvershootDeringing bool = true

	// DefaultOptimizeScans is the default optimize scans setting to use when shrinking an image.
	DefaultOptimizeScans bool = true

	// DefaultResize is the default resize setting to use when shrinking an image.
	DefaultResize bool = false

	// DefaultWidth is the default width to use when shrinking an image.
	DefaultWidth int = 0

	// DefaultHeight is the default height to use when shrinking an image.
	DefaultHeight int = 0
)

// ShrinkHandler is an HTTP handler for the /shrink endpoint.
type ShrinkHandler struct {
	cfg         *config.Config
	fetchClient *fetch.Client
	logger      *zap.Logger
}

// NewShrinkHandler creates a new instance of ShrinkHandler.
func NewShrinkHandler(cfg *config.Config, fetchClient *fetch.Client, logger *zap.Logger) *ShrinkHandler {
	return &ShrinkHandler{
		cfg:         cfg,
		fetchClient: fetchClient,
		logger:      logger,
	}
}

// ServeHTTP serves the /shrink endpoint.
func (h *ShrinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > int64(h.cfg.Service.MaxUploadSize)<<20 {
		serror.JSON(w, h.logger, serror.ErrorResponse{
			Code:    http.StatusRequestEntityTooLarge,
			Message: fmt.Sprintf("The image you uploaded is too large. The maximum upload size is %d. Please try again.", h.cfg.Service.MaxUploadSize),
		})

		return
	}

	var (
		quality            = DefaultQuality
		compression        = DefaultCompression
		quantTable         = DefaultQuantTable
		optimizeCoding     = DefaultOptimizeCoding
		interlace          = DefaultInterlace
		stripMetadata      = DefaultStripMetadata
		optimizeICCProfile = DefaultOptimizeICCProfile
		trellisQuant       = DefaultTrellisQuant
		overshootDeringing = DefaultOvershootDeringing
		optimizeScans      = DefaultOptimizeScans
		resize             = DefaultResize
		width              = DefaultWidth
		height             = DefaultHeight
		err                error
	)

	if r.URL.Query().Get("quality") != "" {
		quality, err = strconv.Atoi(r.URL.Query().Get("quality"))
		if err != nil {
			h.logger.Error("failed to parse image quality parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image quality parameter. Please provide a valid integer and try again.",
			})

			return
		}

		if quality < 1 || quality > 100 {
			quality = DefaultQuality
		}
	}

	if r.URL.Query().Get("compression") != "" {
		compression, err = strconv.Atoi(r.URL.Query().Get("compression"))
		if err != nil {
			h.logger.Error("failed to parse image compression parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image compression parameter. Please provide a valid integer and try again.",
			})

			return
		}

		if compression < 0 || compression > 9 {
			compression = DefaultCompression
		}
	}

	if r.URL.Query().Get("quant_table") != "" {
		quantTable, err = strconv.Atoi(r.URL.Query().Get("quant_table"))
		if err != nil {
			h.logger.Error("failed to parse image quantization table parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image quantization table parameter. Please provide a valid integer and try again.",
			})

			return
		}

		if quantTable < 0 || quantTable > 8 {
			quantTable = DefaultQuantTable
		}
	}

	if r.URL.Query().Get("optimize_coding") != "" {
		optimizeCoding, err = strconv.ParseBool(r.URL.Query().Get("optimize_coding"))
		if err != nil {
			h.logger.Error("failed to parse image optimize coding parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image optimize coding parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("interlace") != "" {
		interlace, err = strconv.ParseBool(r.URL.Query().Get("interlace"))
		if err != nil {
			h.logger.Error("failed to parse image interlace parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image interlace parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("strip") != "" {
		stripMetadata, err = strconv.ParseBool(r.URL.Query().Get("strip_metadata"))
		if err != nil {
			h.logger.Error("failed to parse image strip metadata parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image strip metadata parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("optimize_icc_profile") != "" {
		optimizeICCProfile, err = strconv.ParseBool(r.URL.Query().Get("optimize_icc_profile"))
		if err != nil {
			h.logger.Error("failed to parse image optimize ICC profile parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image optimize ICC profile parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("trellis_quant") != "" {
		trellisQuant, err = strconv.ParseBool(r.URL.Query().Get("trellis_quant"))
		if err != nil {
			h.logger.Error("failed to parse image trellis quant parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image trellis quant parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("overshoot_deringing") != "" {
		overshootDeringing, err = strconv.ParseBool(r.URL.Query().Get("overshoot_deringing"))
		if err != nil {
			h.logger.Error("failed to parse image overshoot deringing parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image overshoot deringing parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("optimize_scans") != "" {
		optimizeScans, err = strconv.ParseBool(r.URL.Query().Get("optimize_scans"))
		if err != nil {
			h.logger.Error("failed to parse image optimize scans parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image optimize scans parameter. Please provide a valid boolean and try again.",
			})

			return
		}
	}

	if r.URL.Query().Get("width") != "" {
		width, err = strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil {
			h.logger.Error("failed to parse image width parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image width parameter. Please provide a valid integer and try again.",
			})

			return
		}

		if width > int(h.cfg.Service.MaxAllowedWidth) {
			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("The image width parameter cannot be greater than %d. Please provide a valid integer and try again.", h.cfg.Service.MaxAllowedWidth),
			})

			return
		}

		if width < 0 {
			width = DefaultWidth
		}
	}

	if r.URL.Query().Get("height") != "" {
		height, err = strconv.Atoi(r.URL.Query().Get("height"))
		if err != nil {
			h.logger.Error("failed to parse image height parameter", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the image height parameter. Please provide a valid integer and try again.",
			})

			return
		}

		if height > int(h.cfg.Service.MaxAllowedHeight) {
			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("The image height parameter cannot be greater than %d. Please provide a valid integer and try again.", h.cfg.Service.MaxAllowedHeight),
			})

			return
		}

		if height < 0 {
			height = DefaultHeight
		}
	}

	if width > 0 || height > 0 {
		resize = true
	}

	options := &imgdiet.Options{
		Quality:            uint(quality),
		Compression:        uint(compression),
		QuantTable:         uint(quantTable),
		OptimizeCoding:     optimizeCoding,
		Interlaced:         interlace,
		StripMetadata:      stripMetadata,
		OptimizeICCProfile: optimizeICCProfile,
		TrellisQuant:       trellisQuant,
		OvershootDeringing: overshootDeringing,
		OptimizeScans:      optimizeScans,
	}

	var (
		maxUploadSize = h.cfg.Service.MaxUploadSize * 1024 * 1024
		uri           = r.URL.Query().Get("url")
		file          io.ReadCloser
		header        *multipart.FileHeader
		data          []byte
		filename      string
	)

	if uri != "" {
		data, filename, err = h.fetchClient.Remote(r.Context(), uri)
		if err != nil {
			h.logger.Error("failed to fetch remote image", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot fetch the remote image. Please provide a valid image and try again.",
			})

			return
		}

		if len(data) > int(maxUploadSize) {
			h.logger.Error("image size is too large", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("The image size cannot be greater than %d MB. Please provide a valid image and try again.", h.cfg.Service.MaxUploadSize),
			})

			return
		}

		file = io.NopCloser(bytes.NewReader(data))
		defer file.Close()
	} else {
		r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize))
		if err = r.ParseMultipartForm(int64(maxUploadSize)); err != nil {
			h.logger.Error("failed to parse multipart form", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot process the image. Please provide a valid image and try again.",
			})

			return
		}

		file, header, err = r.FormFile("input")
		if err != nil {
			h.logger.Error("failed to get input file", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot process the image. Please provide a valid image and try again.",
			})

			return
		}
		defer file.Close()

		filename = header.Filename
	}

	img, err := imgdiet.Open(file)
	if err != nil {
		if errors.Is(err, imgdiet.ErrUnsupportedImageFormat) {
			h.logger.Error("unsupported image format", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusUnsupportedMediaType,
				Message: "Unsupported image format. Please provide a valid JPEG or PNG image and try again.",
			})

			return
		}

		h.logger.Error("failed to open image", zap.Error(err))

		serror.JSON(w, h.logger, serror.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Cannot process the image. Please try again.",
		})

		return
	}
	defer img.Close()

	var optimizedImage []byte

	if resize {
		optimizedImage, err = img.Resize(uint(width), uint(height), options)
		if err != nil {
			h.logger.Error("failed to resize image", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Cannot process the image. Please try again.",
			})

			return
		}
	} else {
		optimizedImage, err = img.Optimize(options)
		if err != nil {
			h.logger.Error("failed to optimize image", zap.Error(err))

			serror.JSON(w, h.logger, serror.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Cannot process the image. Please try again.",
			})

			return
		}
	}

	w.Header().Set("Content-Type", "octet/stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	_, err = io.Copy(w, bytes.NewReader(optimizedImage))
	if err != nil {
		h.logger.Error("failed to write optimized image", zap.Error(err))

		serror.JSON(w, h.logger, serror.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Cannot process the image. Please try again.",
		})

		return
	}
}
