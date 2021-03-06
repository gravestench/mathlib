package mathlib

import "math"

// define the different euler orders
const (
	EulerOrderXYZ = iota
	EulerOrderYXZ
	EulerOrderZXY
	EulerOrderZYX
	EulerOrderYZX
	EulerOrderXZY
	numEulerOrders
)

// and define the default
const (
	EulerOrderDefault = EulerOrderXYZ
)

func eulerNoop(_ *Euler) { /* do nothing */ }

// static check that euler is Vector3Like
var _ Vector3Like = &Euler{}

// NewEuler creates a Euler
func NewEuler(x, y, z float64, order int) *Euler {
	return &Euler{
		X:                x,
		Y:                y,
		Z:                z,
		Order:            order,
		OnChangeCallback: eulerNoop,
	}
}

// Euler is an abstraction of a Euler angle
type Euler struct {
	X, Y, Z          float64
	Order            int
	OnChangeCallback func(*Euler)
}

// XY returns the x and y components of the quaternion
func (e *Euler) XY() (x, y float64) {
	return e.X, e.Y
}

// XYZ returns the x, y and z components of the quaternion
func (e *Euler) XYZ() (x, y, z float64) {
	return e.X, e.Y, e.Z
}

// SetX sets the x component
func (e *Euler) SetX(v float64) *Euler {
	return e.Set(v, e.Y, e.Z, e.Order)
}

// SetY sets the y component
func (e *Euler) SetY(v float64) *Euler {
	return e.Set(e.X, v, e.Z, e.Order)
}

// SetZ sets the z component
func (e *Euler) SetZ(v float64) *Euler {
	return e.Set(e.X, e.Y, v, e.Order)
}

// SetOrder sets the order of the components
func (e *Euler) SetOrder(v int) *Euler {
	return e.Set(e.X, e.Y, e.Z, v)
}

// Set sets the x, y, and z components, as well as the order
func (e *Euler) Set(x, y, z float64, order int) *Euler {
	order = int(Clamp(float64(order), 0, numEulerOrders-1))
	e.X, e.Y, e.Z, e.Order = x, y, z, order

	e.OnChangeCallback(e)

	return e
}

// Copy copies the values and order from the given Euler into this one
func (e *Euler) Copy(other *Euler) *Euler {
	return e.Set(other.X, other.Y, other.Z, other.Order)
}

// SetFromQuaternion sets the values from a Quarternion in the specified order
func (e *Euler) SetFromQuaternion(q *Quaternion, order int) *Euler {
	tmpMat4 := NewMatrix4(nil)

	tmpMat4.FromQuaternion(q)

	return e.SetFromRotationMatrix(tmpMat4, order)
}

// SetFromRotationMatrix sets the values from a matrix in the specified order
func (e *Euler) SetFromRotationMatrix(m4 *Matrix4, order int) *Euler {
	m := m4.Values

	m11, m12, m13,
		m21, m22, m23,
		m31, m32, m33 :=
		m[0], m[4], m[8],
		m[1], m[5], m[9],
		m[2], m[6], m[10]

	x, y, z := 0., 0., 0.
	epsilon := 0.99999

	switch e.Order {
	case EulerOrderYXZ:
		x = math.Asin(-Clamp(m23, -1, 1))

		if math.Abs(m23) < epsilon {
			y = math.Atan2(m13, m33)
			z = math.Atan2(m21, m22)
		} else {
			y = math.Atan2(-m31, m11)
		}
	case EulerOrderZXY:
		x = math.Asin(Clamp(m32, -1, 1))

		if math.Abs(m32) < epsilon {
			y = math.Atan2(-m31, m33)
			z = math.Atan2(-m12, m22)
		} else {
			z = math.Atan2(m21, m11)
		}
	case EulerOrderZYX:
		y = math.Asin(-Clamp(m31, -1, 1))

		if math.Abs(m31) < epsilon {
			x = math.Atan2(m32, m33)
			z = math.Atan2(m21, m11)
		} else {
			z = math.Atan2(-m12, m22)
		}
	case EulerOrderYZX:
		z = math.Asin(Clamp(m21, -1, 1))

		if math.Abs(m21) < epsilon {
			x = math.Atan2(-m23, m22)
			y = math.Atan2(-m31, m11)
		} else {
			y = math.Atan2(m13, m33)
		}
	case EulerOrderXZY:
		z = math.Asin(-Clamp(m12, -1, 1))

		if math.Abs(m12) < epsilon {
			x = math.Atan2(m32, m22)
			y = math.Atan2(m13, m11)
		} else {
			x = math.Atan2(-m23, m33)
		}
	case EulerOrderXYZ:
		fallthrough //nolint:gocritic // it's better to be explicit and include the fallthrough to default
	default:
		y = math.Asin(Clamp(m13, -1, 1))

		if math.Abs(m13) < epsilon {
			x = math.Atan2(-m23, m33)
			z = math.Atan2(-m12, m11)
		} else {
			x = math.Atan2(m32, m22)
		}
	}

	return e.Set(x, y, z, order)
}
