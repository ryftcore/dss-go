// Ported from dss-enumerations/.../ImageScaling.java (DSS 6.5.RC1).
package enumerations

// ImageScaling defines an enumeration that can be used to define the image
// scaling behavior for signature fields with a fixed size.
type ImageScaling string

const (
	// ImageScalingStretch stretches the image in both directions in order
	// to fill the signature field box.
	ImageScalingStretch ImageScaling = "STRETCH"
	// ImageScalingZoomAndCenter zooms the image to the closest dimension
	// without stretching and centers the image in other direction.
	ImageScalingZoomAndCenter ImageScaling = "ZOOM_AND_CENTER"
	// ImageScalingCenter keeps the original image size and centers the
	// image in both directions.
	ImageScalingCenter ImageScaling = "CENTER"
)

// ImageScalingValues returns all constants in declaration order.
func ImageScalingValues() []ImageScaling {
	return []ImageScaling{
		ImageScalingStretch,
		ImageScalingZoomAndCenter,
		ImageScalingCenter,
	}
}

// ImageScalingValueOf returns the constant matching the given Java enum
// name.
func ImageScalingValueOf(name string) (ImageScaling, error) {
	for _, v := range ImageScalingValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &imageScalingInvalidValueError{name}
}

type imageScalingInvalidValueError struct {
	name string
}

func (e *imageScalingInvalidValueError) Error() string {
	return "no enum constant ImageScaling." + e.name
}
