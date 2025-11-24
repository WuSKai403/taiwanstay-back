package service

import (
	"testing"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/stretchr/testify/assert"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
)

func TestDetermineImageStatus(t *testing.T) {
	// Setup default config
	cfg := &config.Config{
		Image: config.ImageConfig{
			RejectAdult:    "LIKELY",
			RejectSpoof:    "LIKELY",
			RejectMedical:  "LIKELY",
			RejectViolence: "LIKELY",
			RejectRacy:     "LIKELY",
		},
	}

	// Create a service instance with just the config (other fields can be nil for this test)
	svc := &imageService{
		cfg: cfg,
	}

	tests := []struct {
		name        string
		annotations *visionpb.SafeSearchAnnotation
		want        domain.ImageStatus
	}{
		{
			name: "All Very Unlikely -> Approved",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_VERY_UNLIKELY,
				Spoof:    visionpb.Likelihood_VERY_UNLIKELY,
				Medical:  visionpb.Likelihood_VERY_UNLIKELY,
				Violence: visionpb.Likelihood_VERY_UNLIKELY,
				Racy:     visionpb.Likelihood_VERY_UNLIKELY,
			},
			want: domain.ImageStatusApproved,
		},
		{
			name: "Adult Likely -> Rejected",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_LIKELY,
				Spoof:    visionpb.Likelihood_VERY_UNLIKELY,
				Medical:  visionpb.Likelihood_VERY_UNLIKELY,
				Violence: visionpb.Likelihood_VERY_UNLIKELY,
				Racy:     visionpb.Likelihood_VERY_UNLIKELY,
			},
			want: domain.ImageStatusRejected,
		},
		{
			name: "Violence Possible (Below Threshold) -> Approved",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_VERY_UNLIKELY,
				Spoof:    visionpb.Likelihood_VERY_UNLIKELY,
				Medical:  visionpb.Likelihood_VERY_UNLIKELY,
				Violence: visionpb.Likelihood_POSSIBLE,
				Racy:     visionpb.Likelihood_VERY_UNLIKELY,
			},
			want: domain.ImageStatusApproved,
		},
		{
			name: "Violence Likely (At Threshold) -> Rejected",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_VERY_UNLIKELY,
				Spoof:    visionpb.Likelihood_VERY_UNLIKELY,
				Medical:  visionpb.Likelihood_VERY_UNLIKELY,
				Violence: visionpb.Likelihood_LIKELY,
				Racy:     visionpb.Likelihood_VERY_UNLIKELY,
			},
			want: domain.ImageStatusRejected,
		},
		{
			name:        "Nil Annotations -> Pending",
			annotations: nil,
			want:        domain.ImageStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.determineImageStatus(tt.annotations)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetermineImageStatus_CustomConfig(t *testing.T) {
	// Setup custom config (Strict on Adult, Lenient on Violence)
	cfg := &config.Config{
		Image: config.ImageConfig{
			RejectAdult:    "POSSIBLE", // Stricter
			RejectSpoof:    "LIKELY",
			RejectMedical:  "LIKELY",
			RejectViolence: "VERY_LIKELY", // More lenient
			RejectRacy:     "LIKELY",
		},
	}

	svc := &imageService{
		cfg: cfg,
	}

	tests := []struct {
		name        string
		annotations *visionpb.SafeSearchAnnotation
		want        domain.ImageStatus
	}{
		{
			name: "Adult Possible (At Stricter Threshold) -> Rejected",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_POSSIBLE,
				Violence: visionpb.Likelihood_LIKELY, // Below custom threshold
			},
			want: domain.ImageStatusRejected,
		},
		{
			name: "Violence Likely (Below Lenient Threshold) -> Approved",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_UNLIKELY,
				Violence: visionpb.Likelihood_LIKELY,
			},
			want: domain.ImageStatusApproved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.determineImageStatus(tt.annotations)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetermineImageStatus_EmptyConfig(t *testing.T) {
	// Setup empty config
	cfg := &config.Config{
		Image: config.ImageConfig{
			RejectAdult:    "",
			RejectSpoof:    "",
			RejectMedical:  "",
			RejectViolence: "",
			RejectRacy:     "",
		},
	}

	svc := &imageService{
		cfg: cfg,
	}

	tests := []struct {
		name        string
		annotations *visionpb.SafeSearchAnnotation
		want        domain.ImageStatus
	}{
		{
			name: "Empty Config -> All Approved",
			annotations: &visionpb.SafeSearchAnnotation{
				Adult:    visionpb.Likelihood_VERY_LIKELY,
				Violence: visionpb.Likelihood_VERY_LIKELY,
			},
			want: domain.ImageStatusApproved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.determineImageStatus(tt.annotations)
			assert.Equal(t, tt.want, got)
		})
	}
}
