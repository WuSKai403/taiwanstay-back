package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/google/uuid"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
	"github.com/taiwanstay/taiwanstay-back/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageService interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID string) (*domain.Image, error)
	GetImage(ctx context.Context, id string) (*domain.Image, error)
	UpdateImageStatus(ctx context.Context, id string, status domain.ImageStatus) error
	GetImageContent(ctx context.Context, id string) (io.ReadCloser, error)
}

type imageService struct {
	repo          repository.ImageRepository
	storageClient *storage.Client
	visionClient  *vision.ImageAnnotatorClient
	publicBucket  string
	privateBucket string
	cfg           *config.Config
}

func NewImageService(repo repository.ImageRepository, storageClient *storage.Client, visionClient *vision.ImageAnnotatorClient, cfg *config.Config) ImageService {
	return &imageService{
		repo:          repo,
		storageClient: storageClient,
		visionClient:  visionClient,
		publicBucket:  cfg.GCP.PublicBucket,
		privateBucket: cfg.GCP.PrivateBucket,
		cfg:           cfg,
	}
}

func (s *imageService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID string) (*domain.Image, error) {
	// 1. Generate unique filename
	ext := "jpg" // Simplified: assume jpg or detect from header
	filename := fmt.Sprintf("%s/%s.%s", userID, uuid.New().String(), ext)

	// 2. Upload to Private Bucket initially
	wc := s.storageClient.Bucket(s.privateBucket).Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, fmt.Errorf("failed to upload to GCS: %w", err)
	}
	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("failed to close GCS writer: %w", err)
	}

	// 3. Analyze with Vision API (from Private Bucket)
	gcsURI := fmt.Sprintf("gs://%s/%s", s.privateBucket, filename)

	req := &visionpb.AnnotateImageRequest{
		Image: &visionpb.Image{
			Source: &visionpb.ImageSource{
				GcsImageUri: gcsURI,
			},
		},
		Features: []*visionpb.Feature{
			{Type: visionpb.Feature_SAFE_SEARCH_DETECTION},
		},
	}

	batchReq := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	}

	resp, err := s.visionClient.BatchAnnotateImages(ctx, batchReq)
	if err != nil {
		logger.ErrorContext(ctx, "Vision API failed", "error", err)
		// Don't fail upload, just mark as pending manual review
	}

	// 4. Determine Status
	status := domain.ImageStatusPending
	visionData := domain.VisionAIRawData{}

	var annotations *visionpb.SafeSearchAnnotation
	if resp != nil && len(resp.Responses) > 0 {
		annotations = resp.Responses[0].SafeSearchAnnotation
	}

	if annotations != nil {
		visionData.Adult = annotations.Adult.String()
		visionData.Racy = annotations.Racy.String()
		visionData.Violence = annotations.Violence.String()
		visionData.Medical = annotations.Medical.String()
		visionData.Spoof = annotations.Spoof.String()
	}

	// 4. Determine Status
	status = s.determineImageStatus(annotations)

	// 5. If Approved, move to Public Bucket and generate Public URL
	publicURL := ""
	if status == domain.ImageStatusApproved {
		if err := s.moveFile(ctx, s.privateBucket, filename, s.publicBucket, filename); err != nil {
			logger.ErrorContext(ctx, "Failed to move approved image to public bucket", "error", err)
			// Fallback to pending if move fails
			status = domain.ImageStatusPending
		} else {
			// Generate ImageKit URL
			// Format: https://ik.imagekit.io/your_id/path/to/file
			if s.cfg.Image.ImageKitEndpoint != "" {
				// Ensure endpoint doesn't end with slash and filename doesn't start with slash (it doesn't)
				publicURL = fmt.Sprintf("%s/%s", s.cfg.Image.ImageKitEndpoint, filename)
			}
		}
	}

	// 6. Save to DB
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	image := &domain.Image{
		UserID:     userObjID,
		GCSPath:    filename,
		PublicURL:  publicURL,
		Status:     status,
		VisionData: visionData,
	}

	if err := s.repo.Create(ctx, image); err != nil {
		return nil, err
	}

	return image, nil
}

func (s *imageService) determineImageStatus(annotations *visionpb.SafeSearchAnnotation) domain.ImageStatus {
	if annotations == nil {
		return domain.ImageStatusPending
	}

	// Helper to parse likelihood string to enum
	parseLikelihood := func(s string) visionpb.Likelihood {
		v, ok := visionpb.Likelihood_value[s]
		if !ok {
			return visionpb.Likelihood_UNKNOWN
		}
		return visionpb.Likelihood(v)
	}

	// Helper to check if any category exceeds reject threshold
	isRejected := func() bool {
		if s.cfg.Image.RejectAdult != "" && annotations.Adult >= parseLikelihood(s.cfg.Image.RejectAdult) {
			return true
		}
		if s.cfg.Image.RejectSpoof != "" && annotations.Spoof >= parseLikelihood(s.cfg.Image.RejectSpoof) {
			return true
		}
		if s.cfg.Image.RejectMedical != "" && annotations.Medical >= parseLikelihood(s.cfg.Image.RejectMedical) {
			return true
		}
		if s.cfg.Image.RejectViolence != "" && annotations.Violence >= parseLikelihood(s.cfg.Image.RejectViolence) {
			return true
		}
		if s.cfg.Image.RejectRacy != "" && annotations.Racy >= parseLikelihood(s.cfg.Image.RejectRacy) {
			return true
		}
		return false
	}

	if isRejected() {
		return domain.ImageStatusRejected
	}
	return domain.ImageStatusApproved
}

func (s *imageService) GetImage(ctx context.Context, id string) (*domain.Image, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *imageService) UpdateImageStatus(ctx context.Context, id string, newStatus domain.ImageStatus) error {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if image.Status == newStatus {
		return nil
	}

	// Move file logic
	// If moving TO Approved -> Private to Public
	// If moving FROM Approved -> Public to Private

	if newStatus == domain.ImageStatusApproved && image.Status != domain.ImageStatusApproved {
		// Move Private -> Public
		if err := s.moveFile(ctx, s.privateBucket, image.GCSPath, s.publicBucket, image.GCSPath); err != nil {
			return err
		}
	} else if image.Status == domain.ImageStatusApproved && newStatus != domain.ImageStatusApproved {
		// Move Public -> Private
		if err := s.moveFile(ctx, s.publicBucket, image.GCSPath, s.privateBucket, image.GCSPath); err != nil {
			return err
		}
	}

	return s.repo.UpdateStatus(ctx, id, newStatus)
}

func (s *imageService) GetImageContent(ctx context.Context, id string) (io.ReadCloser, error) {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bucket := s.privateBucket
	if image.Status == domain.ImageStatusApproved {
		bucket = s.publicBucket
	}

	rc, err := s.storageClient.Bucket(bucket).Object(image.GCSPath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// Helper to move file between buckets
func (s *imageService) moveFile(ctx context.Context, srcBucket, srcObject, dstBucket, dstObject string) error {
	src := s.storageClient.Bucket(srcBucket).Object(srcObject)
	dst := s.storageClient.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete source object: %w", err)
	}
	return nil
}
