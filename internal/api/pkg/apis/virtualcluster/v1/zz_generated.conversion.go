//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by conversion-gen. DO NOT EDIT.

package v1

import (
	unsafe "unsafe"

	virtualcluster "dev.khulnasoft.com/api/pkg/apis/virtualcluster"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*HelmRelease)(nil), (*virtualcluster.HelmRelease)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HelmRelease_To_virtualcluster_HelmRelease(a.(*HelmRelease), b.(*virtualcluster.HelmRelease), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*virtualcluster.HelmRelease)(nil), (*HelmRelease)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_virtualcluster_HelmRelease_To_v1_HelmRelease(a.(*virtualcluster.HelmRelease), b.(*HelmRelease), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmReleaseList)(nil), (*virtualcluster.HelmReleaseList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HelmReleaseList_To_virtualcluster_HelmReleaseList(a.(*HelmReleaseList), b.(*virtualcluster.HelmReleaseList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*virtualcluster.HelmReleaseList)(nil), (*HelmReleaseList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_virtualcluster_HelmReleaseList_To_v1_HelmReleaseList(a.(*virtualcluster.HelmReleaseList), b.(*HelmReleaseList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmReleaseSpec)(nil), (*virtualcluster.HelmReleaseSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HelmReleaseSpec_To_virtualcluster_HelmReleaseSpec(a.(*HelmReleaseSpec), b.(*virtualcluster.HelmReleaseSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*virtualcluster.HelmReleaseSpec)(nil), (*HelmReleaseSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_virtualcluster_HelmReleaseSpec_To_v1_HelmReleaseSpec(a.(*virtualcluster.HelmReleaseSpec), b.(*HelmReleaseSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmReleaseStatus)(nil), (*virtualcluster.HelmReleaseStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HelmReleaseStatus_To_virtualcluster_HelmReleaseStatus(a.(*HelmReleaseStatus), b.(*virtualcluster.HelmReleaseStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*virtualcluster.HelmReleaseStatus)(nil), (*HelmReleaseStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_virtualcluster_HelmReleaseStatus_To_v1_HelmReleaseStatus(a.(*virtualcluster.HelmReleaseStatus), b.(*HelmReleaseStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_HelmRelease_To_virtualcluster_HelmRelease(in *HelmRelease, out *virtualcluster.HelmRelease, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_HelmReleaseSpec_To_virtualcluster_HelmReleaseSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_HelmReleaseStatus_To_virtualcluster_HelmReleaseStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_HelmRelease_To_virtualcluster_HelmRelease is an autogenerated conversion function.
func Convert_v1_HelmRelease_To_virtualcluster_HelmRelease(in *HelmRelease, out *virtualcluster.HelmRelease, s conversion.Scope) error {
	return autoConvert_v1_HelmRelease_To_virtualcluster_HelmRelease(in, out, s)
}

func autoConvert_virtualcluster_HelmRelease_To_v1_HelmRelease(in *virtualcluster.HelmRelease, out *HelmRelease, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_virtualcluster_HelmReleaseSpec_To_v1_HelmReleaseSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_virtualcluster_HelmReleaseStatus_To_v1_HelmReleaseStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_virtualcluster_HelmRelease_To_v1_HelmRelease is an autogenerated conversion function.
func Convert_virtualcluster_HelmRelease_To_v1_HelmRelease(in *virtualcluster.HelmRelease, out *HelmRelease, s conversion.Scope) error {
	return autoConvert_virtualcluster_HelmRelease_To_v1_HelmRelease(in, out, s)
}

func autoConvert_v1_HelmReleaseList_To_virtualcluster_HelmReleaseList(in *HelmReleaseList, out *virtualcluster.HelmReleaseList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]virtualcluster.HelmRelease)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1_HelmReleaseList_To_virtualcluster_HelmReleaseList is an autogenerated conversion function.
func Convert_v1_HelmReleaseList_To_virtualcluster_HelmReleaseList(in *HelmReleaseList, out *virtualcluster.HelmReleaseList, s conversion.Scope) error {
	return autoConvert_v1_HelmReleaseList_To_virtualcluster_HelmReleaseList(in, out, s)
}

func autoConvert_virtualcluster_HelmReleaseList_To_v1_HelmReleaseList(in *virtualcluster.HelmReleaseList, out *HelmReleaseList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]HelmRelease)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_virtualcluster_HelmReleaseList_To_v1_HelmReleaseList is an autogenerated conversion function.
func Convert_virtualcluster_HelmReleaseList_To_v1_HelmReleaseList(in *virtualcluster.HelmReleaseList, out *HelmReleaseList, s conversion.Scope) error {
	return autoConvert_virtualcluster_HelmReleaseList_To_v1_HelmReleaseList(in, out, s)
}

func autoConvert_v1_HelmReleaseSpec_To_virtualcluster_HelmReleaseSpec(in *HelmReleaseSpec, out *virtualcluster.HelmReleaseSpec, s conversion.Scope) error {
	out.HelmReleaseSpec = in.HelmReleaseSpec
	return nil
}

// Convert_v1_HelmReleaseSpec_To_virtualcluster_HelmReleaseSpec is an autogenerated conversion function.
func Convert_v1_HelmReleaseSpec_To_virtualcluster_HelmReleaseSpec(in *HelmReleaseSpec, out *virtualcluster.HelmReleaseSpec, s conversion.Scope) error {
	return autoConvert_v1_HelmReleaseSpec_To_virtualcluster_HelmReleaseSpec(in, out, s)
}

func autoConvert_virtualcluster_HelmReleaseSpec_To_v1_HelmReleaseSpec(in *virtualcluster.HelmReleaseSpec, out *HelmReleaseSpec, s conversion.Scope) error {
	out.HelmReleaseSpec = in.HelmReleaseSpec
	return nil
}

// Convert_virtualcluster_HelmReleaseSpec_To_v1_HelmReleaseSpec is an autogenerated conversion function.
func Convert_virtualcluster_HelmReleaseSpec_To_v1_HelmReleaseSpec(in *virtualcluster.HelmReleaseSpec, out *HelmReleaseSpec, s conversion.Scope) error {
	return autoConvert_virtualcluster_HelmReleaseSpec_To_v1_HelmReleaseSpec(in, out, s)
}

func autoConvert_v1_HelmReleaseStatus_To_virtualcluster_HelmReleaseStatus(in *HelmReleaseStatus, out *virtualcluster.HelmReleaseStatus, s conversion.Scope) error {
	out.HelmReleaseStatus = in.HelmReleaseStatus
	return nil
}

// Convert_v1_HelmReleaseStatus_To_virtualcluster_HelmReleaseStatus is an autogenerated conversion function.
func Convert_v1_HelmReleaseStatus_To_virtualcluster_HelmReleaseStatus(in *HelmReleaseStatus, out *virtualcluster.HelmReleaseStatus, s conversion.Scope) error {
	return autoConvert_v1_HelmReleaseStatus_To_virtualcluster_HelmReleaseStatus(in, out, s)
}

func autoConvert_virtualcluster_HelmReleaseStatus_To_v1_HelmReleaseStatus(in *virtualcluster.HelmReleaseStatus, out *HelmReleaseStatus, s conversion.Scope) error {
	out.HelmReleaseStatus = in.HelmReleaseStatus
	return nil
}

// Convert_virtualcluster_HelmReleaseStatus_To_v1_HelmReleaseStatus is an autogenerated conversion function.
func Convert_virtualcluster_HelmReleaseStatus_To_v1_HelmReleaseStatus(in *virtualcluster.HelmReleaseStatus, out *HelmReleaseStatus, s conversion.Scope) error {
	return autoConvert_virtualcluster_HelmReleaseStatus_To_v1_HelmReleaseStatus(in, out, s)
}
